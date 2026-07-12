package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	dbrepo "app/db/repo"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// AllowedEntityTypes defines valid entity types for file association
var AllowedEntityTypes = map[string]bool{
	"machine":              true,
	"spare_part":           true,
	"vendor":               true,
	"maintenance_schedule": true,
}

// AllowedUsageTypes defines valid usage types for files
var AllowedUsageTypes = map[string]bool{
	"thumbnail":  true,
	"attachment": true,
	"report":     true,
	"image":      true,
	"diagram":    true,
	"document":   true,
	"manual":     true,
	"warranty":   true,
}

// AllowedMimeTypes defines valid mime types for upload
var AllowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"application/pdf": true,
}

// FilePatchRequest allows updating limited fields on a file link
type FilePatchRequest struct {
	IsPrimary *bool `json:"is_primary"`
}

// getUploadDir returns the upload directory path
func getUploadDir() string {
	dir := os.Getenv("UPLOAD_DIR")
	if dir == "" {
		dir = "./uploads"
	}
	return dir
}

// getImageLimitPerEntity returns the max number of image files allowed per entity
func getImageLimitPerEntity() int {
	limitStr := os.Getenv("FILE_LIMIT_IMAGE_PER_ENTITY")
	if limitStr == "" {
		return 4 // default
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return 4
	}
	return limit
}

// getDocumentLimitPerEntity returns the max number of document files allowed per entity
func getDocumentLimitPerEntity() int {
	limitStr := os.Getenv("FILE_LIMIT_DOCUMENT_PER_ENTITY")
	if limitStr == "" {
		return 10 // default
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return 10
	}
	return limit
}

// UploadFile godoc
// @Summary Upload File
// @Description Upload a file and associate it with an entity (machine, spare_part, vendor)
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload (image or PDF)"
// @Param entity_type formData string true "Entity type (machine, spare_part, vendor)"
// @Param entity_uid formData string true "Entity UID"
// @Param usage_type formData string true "Usage type (thumbnail, attachment, report, image, diagram, document , warranty , manual)"
// @Param display_order formData int false "Display order"
// @Param is_primary formData bool false "Is primary file"
// @Success 201 {object} map[string]string "{\"uid\":\"550e8400-e29b-41d4-a716-446655440000\"}"
// @Failure 400 {object} map[string]string "{\"error\":\"Invalid request\"}"
// @Failure 500 {object} map[string]string "{\"error\":\"Internal server error\"}"
// @Router /files [post]
func UploadFile(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UploadFile: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	entityType := strings.TrimSpace(c.FormValue("entity_type"))
	entityUID := strings.TrimSpace(c.FormValue("entity_uid"))
	usageType := strings.TrimSpace(c.FormValue("usage_type"))
	displayOrderStr := c.FormValue("display_order")
	isPrimaryStr := c.FormValue("is_primary")

	if !AllowedEntityTypes[entityType] {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": fmt.Sprintf("invalid entity_type: %s. Allowed: machine, spare_part, vendor", entityType),
		})
	}

	if entityUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "entity_uid is required"})
	}

	if !AllowedUsageTypes[usageType] {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": fmt.Sprintf("invalid usage_type: %s. Allowed: thumbnail, attachment, report, image, diagram, document , warranty , manual", usageType),
		})
	}

	displayOrder := 0
	if displayOrderStr != "" {
		var err error
		displayOrder, err = strconv.Atoi(displayOrderStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid display_order"})
		}
	}

	isPrimary := false
	if isPrimaryStr != "" {
		isPrimary = isPrimaryStr == "true" || isPrimaryStr == "1"
	}

	entityID, err := dbrepo.GetEntityIDByTypeAndUID(pool, entityType, entityUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{
				"error": fmt.Sprintf("%s with UID %s not found", entityType, entityUID),
			})
		}
		c.Logger().Errorf("UploadFile: error resolving entity: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to resolve entity"})
	}

	// Check file limits per entity and usage_type
	if usageType == "image" || usageType == "document" {
		fileFilters := dbrepo.FileListFilters{
			EntityType: &entityType,
			EntityID:   &entityID,
			UsageType:  &usageType,
		}
		existingFiles, _, _, err := dbrepo.ListFilesWithFilters(pool, fileFilters)
		if err != nil {
			c.Logger().Errorf("UploadFile: error checking existing files: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to check existing files"})
		}

		currentCount := len(existingFiles)
		var limit int
		if usageType == "image" {
			limit = getImageLimitPerEntity()
		} else {
			limit = getDocumentLimitPerEntity()
		}

		if currentCount >= limit {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error":         fmt.Sprintf("maximum %s files limit reached for this %s", usageType, entityType),
				"current_count": currentCount,
				"limit":         limit,
			})
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "file is required"})
	}

	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "file size exceeds 10MB limit"})
	}

	src, err := file.Open()
	if err != nil {
		c.Logger().Errorf("UploadFile: error opening file: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to open uploaded file"})
	}
	defer src.Close()

	fileContent, err := io.ReadAll(src)
	if err != nil {
		c.Logger().Errorf("UploadFile: error reading file: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to read uploaded file"})
	}

	mimeType := http.DetectContentType(fileContent)
	if !AllowedMimeTypes[mimeType] {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": fmt.Sprintf("invalid file type: %s. Allowed: images (jpeg, png, gif, webp) and PDF", mimeType),
		})
	}

	hash := sha256.Sum256(fileContent)
	checksum := hex.EncodeToString(hash[:])

	originalName := file.Filename
	extension := strings.ToLower(strings.TrimPrefix(filepath.Ext(originalName), "."))
	if extension == "" {
		switch mimeType {
		case "image/jpeg":
			extension = "jpg"
		case "image/png":
			extension = "png"
		case "image/gif":
			extension = "gif"
		case "image/webp":
			extension = "webp"
		case "application/pdf":
			extension = "pdf"
		}
	}

	storedName := fmt.Sprintf("%s_%d.%s", uuid.New().String(), time.Now().UnixNano(), extension)

	uploadDir := getUploadDir()
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(uploadDir, entityType, dateDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		c.Logger().Errorf("UploadFile: error creating directory: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create storage directory"})
	}

	storagePath := filepath.Join(entityType, dateDir, storedName)
	fullPath := filepath.Join(uploadDir, storagePath)

	dst, err := os.Create(fullPath)
	if err != nil {
		c.Logger().Errorf("UploadFile: error creating file: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to save file"})
	}
	defer dst.Close()

	if _, err := dst.Write(fileContent); err != nil {
		c.Logger().Errorf("UploadFile: error writing file: %v", err)
		os.Remove(fullPath)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to write file"})
	}

	fileParams := dbrepo.CreateFileParams{
		OriginalName:    originalName,
		StoredName:      storedName,
		Extension:       extension,
		MimeType:        mimeType,
		SizeBytes:       file.Size,
		StorageProvider: "local",
		StoragePath:     storagePath,
		ChecksumSHA256:  &checksum,
		UploadedBy:      1,
	}

	fileUID, fileID, err := dbrepo.CreateFile(pool, fileParams)
	if err != nil {
		c.Logger().Errorf("UploadFile: error creating file record: %v", err)
		os.Remove(fullPath)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create file record"})
	}

	entityFileParams := dbrepo.CreateEntityFileParams{
		FileID:       fileID,
		EntityType:   entityType,
		EntityID:     entityID,
		UsageType:    usageType,
		DisplayOrder: displayOrder,
		IsPrimary:    isPrimary,
		LinkedBy:     1,
	}

	_, err = dbrepo.CreateEntityFile(pool, entityFileParams)
	if err != nil {
		c.Logger().Errorf("UploadFile: error creating entity_file link: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to link file to entity"})
	}

	// If isPrimary is true, reset is_primary for other files in the same group
	if isPrimary {
		if err := dbrepo.ResetIsPrimaryForOthers(pool, entityType, entityID, usageType, fileID); err != nil {
			c.Logger().Warnf("UploadFile: failed to reset is_primary for other files: %v", err)
		}
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"uid":      fileUID,
		"message":  "File uploaded successfully",
		"filename": originalName,
	})
}

// ViewFile godoc
// @Summary Download File
// @Description Download/view a file by UID
// @Tags Files
// @Produce octet-stream
// @Param uid path string true "File UID"
// @Success 200 {file} file
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /files/{uid} [get]
func ViewFile(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DownloadFile: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	fileRow, err := dbrepo.GetFileByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "file not found"})
		}
		c.Logger().Errorf("DownloadFile: db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to retrieve file"})
	}

	// Build full file path
	uploadDir := getUploadDir()
	fullPath := filepath.Join(uploadDir, fileRow.StoragePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.Logger().Errorf("DownloadFile: file not found on disk: %s", fullPath)
		return c.JSON(http.StatusNotFound, echo.Map{"error": "file not found on storage"})
	}

	// Set content disposition for download
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, fileRow.OriginalName))
	c.Response().Header().Set("Content-Type", fileRow.MimeType)

	return c.File(fullPath)
}

// DeleteFile godoc
// @Summary Delete File
// @Description Delete a file by UID
// @Tags Files
// @Produce json
// @Param uid path string true "File UID"
// @Success 200 {object} map[string]string "{\"message\":\"File deleted successfully\"}"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /files/{uid} [delete]
func DeleteFile(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteFile: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	// Get file info first to delete from storage
	fileRow, err := dbrepo.GetFileByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "file not found"})
		}
		c.Logger().Errorf("DeleteFile: db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to retrieve file"})
	}

	// Delete from database (cascade deletes entity_files)
	if err := dbrepo.DeleteFileByUID(pool, uid); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "file not found"})
		}
		c.Logger().Errorf("DeleteFile: db delete error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to delete file record"})
	}

	// Delete from storage (best effort - don't fail if file doesn't exist)
	uploadDir := getUploadDir()
	fullPath := filepath.Join(uploadDir, fileRow.StoragePath)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		c.Logger().Warnf("DeleteFile: failed to delete file from storage: %v", err)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "File deleted successfully"})
}

// PatchFile godoc
// @Summary Update file primary flag
// @Description Update is_primary for a file link; when set true, other files in the same entity/type/usage are cleared
// @Tags Files
// @Accept json
// @Produce json
// @Param uid path string true "File UID"
// @Param file body FilePatchRequest true "Patch payload (is_primary only)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /files/{uid} [patch]
func PatchFile(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}

	var req FilePatchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if req.IsPrimary == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "no fields to update"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("PatchFile: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	filters := dbrepo.FileListFilters{UID: &uid}
	files, _, _, err := dbrepo.ListFilesWithFilters(pool, filters)
	if err != nil {
		c.Logger().Errorf("PatchFile: error fetching file linkage: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch file linkage"})
	}
	if len(files) == 0 || files[0].EntityType == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "file not linked to any entity"})
	}
	file := files[0]

	updateReq := dbrepo.EntityFileUpdateRequest{IsPrimary: req.IsPrimary}
	if err := dbrepo.UpdateEntityFile(pool, file.ID, updateReq); err != nil {
		c.Logger().Errorf("PatchFile: update error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update file"})
	}

	// If marked primary, unset others in same entity/type/usage group
	if req.IsPrimary != nil && *req.IsPrimary {
		if err := dbrepo.ResetIsPrimaryForOthers(pool, file.EntityType, file.EntityID, file.UsageType, file.ID); err != nil {
			c.Logger().Warnf("PatchFile: failed to reset other primaries: %v", err)
		}
	}

	return c.JSON(http.StatusOK, echo.Map{"uid": uid, "is_primary": req.IsPrimary})
}
