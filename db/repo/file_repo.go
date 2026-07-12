package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FileRow represents a file row from the database
type FileRow struct {
	ID              int64
	UID             string
	OriginalName    string
	StoredName      string
	Extension       string
	MimeType        string
	SizeBytes       int64
	StorageProvider string
	StorageBucket   *string
	StoragePath     string
	ChecksumSHA256  *string
	UploadedBy      int64
	UploadedAt      time.Time
}

// EntityFileRow represents an entity_file row from the database
type EntityFileRow struct {
	ID           int64
	FileID       int64
	EntityType   string
	EntityID     int64
	UsageType    string
	DisplayOrder int
	IsPrimary    bool
	LinkedBy     int64
	LinkedAt     time.Time
}

// FileWithEntityRow combines file and entity information
type FileWithEntityRow struct {
	FileRow
	EntityType   string
	EntityID     int64
	UsageType    string
	DisplayOrder int
	IsPrimary    bool
}

// CreateFileParams holds parameters for creating a file record
type CreateFileParams struct {
	OriginalName    string
	StoredName      string
	Extension       string
	MimeType        string
	SizeBytes       int64
	StorageProvider string
	StorageBucket   *string
	StoragePath     string
	ChecksumSHA256  *string
	UploadedBy      int64
}

// CreateEntityFileParams holds parameters for creating an entity_file link
type CreateEntityFileParams struct {
	FileID       int64
	EntityType   string
	EntityID     int64
	UsageType    string
	DisplayOrder int
	IsPrimary    bool
	LinkedBy     int64
}

// EntityFileUpdateRequest holds parameters for updating an entity_file link
type EntityFileUpdateRequest struct {
	IsPrimary    *bool
	DisplayOrder *int
}

// FileListFilters holds optional filters for listing files
type FileListFilters struct {
	// Pagination (optional, defaults: page=1, pageSize=50, max=100)
	Page     *int
	PageSize *int

	// File filters
	UID             *string
	OriginalName    *string
	Extension       *string
	MimeType        *string
	StorageProvider *string
	UploadedBy      *int64

	// Entity filters
	EntityType *string
	EntityID   *int64
	UsageType  *string
	IsPrimary  *bool
}

// CreateFile inserts a new file record and returns its UID and ID
func CreateFile(pool *pgxpool.Pool, params CreateFileParams) (string, int64, error) {
	ctx := context.Background()
	var newUID string
	var newID int64

	err := pool.QueryRow(ctx, `
		INSERT INTO files (
			original_name, stored_name, extension, mime_type, 
			size_bytes, storage_provider, storage_bucket, storage_path, 
			checksum_sha256, uploaded_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING uid, id`,
		params.OriginalName,
		params.StoredName,
		params.Extension,
		params.MimeType,
		params.SizeBytes,
		params.StorageProvider,
		params.StorageBucket,
		params.StoragePath,
		params.ChecksumSHA256,
		params.UploadedBy,
	).Scan(&newUID, &newID)

	if err != nil {
		return "", 0, err
	}
	return newUID, newID, nil
}

// CreateEntityFile creates a link between a file and an entity
func CreateEntityFile(pool *pgxpool.Pool, params CreateEntityFileParams) (int64, error) {
	ctx := context.Background()
	var newID int64

	err := pool.QueryRow(ctx, `
		INSERT INTO entity_files (
			file_id, entity_type, entity_id, usage_type, 
			display_order, is_primary, linked_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		params.FileID,
		params.EntityType,
		params.EntityID,
		params.UsageType,
		params.DisplayOrder,
		params.IsPrimary,
		params.LinkedBy,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil
}

// UpdateEntityFile updates an entity_file record by file_id
// Only updates is_primary and display_order fields
func UpdateEntityFile(pool *pgxpool.Pool, fileID int64, req EntityFileUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []any{}
	idx := 1

	if req.IsPrimary != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_primary = $%d", idx))
		args = append(args, *req.IsPrimary)
		idx++
	}

	if req.DisplayOrder != nil {
		setClauses = append(setClauses, fmt.Sprintf("display_order = $%d", idx))
		args = append(args, *req.DisplayOrder)
		idx++
	}

	// Nothing to update
	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, fileID)
	query := "UPDATE entity_files SET " + strings.Join(setClauses, ", ") +
		" WHERE file_id = " + fmt.Sprintf("$%d", idx)

	_, err := pool.Exec(ctx, query, args...)
	return err
}

// ResetIsPrimaryForOthers sets is_primary=false for all entity_files
// matching entity_type, entity_id, usage_type except the specified file_id
func ResetIsPrimaryForOthers(pool *pgxpool.Pool, entityType string, entityID int64, usageType string, excludeFileID int64) error {
	ctx := context.Background()
	query := `
		UPDATE entity_files 
		SET is_primary = false 
		WHERE entity_type = $1 
		  AND entity_id = $2 
		  AND usage_type = $3 
		  AND file_id != $4
		  AND is_primary = true
	`
	_, err := pool.Exec(ctx, query, entityType, entityID, usageType, excludeFileID)
	return err
}

// GetFileByUID retrieves a file by its UID
func GetFileByUID(pool *pgxpool.Pool, uid string) (*FileRow, error) {
	ctx := context.Background()
	var r FileRow

	err := pool.QueryRow(ctx, `
		SELECT 
			id, uid, original_name, stored_name, extension, 
			mime_type, size_bytes, storage_provider, storage_bucket, 
			storage_path, checksum_sha256, uploaded_by, uploaded_at
		FROM files
		WHERE uid = $1`, uid).Scan(
		&r.ID, &r.UID, &r.OriginalName, &r.StoredName, &r.Extension,
		&r.MimeType, &r.SizeBytes, &r.StorageProvider, &r.StorageBucket,
		&r.StoragePath, &r.ChecksumSHA256, &r.UploadedBy, &r.UploadedAt,
	)

	if err != nil {
		return nil, err
	}
	return &r, nil
}

// DeleteFileByUID deletes a file by its UID (cascade deletes entity_files)
func DeleteFileByUID(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	result, err := pool.Exec(ctx, `DELETE FROM files WHERE uid = $1`, uid)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// ListFilesWithFilters retrieves files with optional filters and pagination
// If no filters are provided, returns all files
// Default page is 1, default page size is 50, max is 100
// Returns: files, currentPage, totalPages, error
func ListFilesWithFilters(
	pool *pgxpool.Pool,
	filters FileListFilters) ([]FileWithEntityRow, int, int, error) {
	ctx := context.Background()

	// Set pagination defaults
	page := 1
	pageSize := 50
	if filters.Page != nil && *filters.Page > 0 {
		page = *filters.Page
	}
	if filters.PageSize != nil && *filters.PageSize > 0 {
		pageSize = *filters.PageSize
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var whereClauses []string
	var args []any
	idx := 1

	// File filters
	if filters.UID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("f.uid = $%d", idx))
		args = append(args, *filters.UID)
		idx++
	}
	if filters.OriginalName != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("f.original_name ILIKE $%d", idx))
		args = append(args, "%"+*filters.OriginalName+"%")
		idx++
	}
	if filters.Extension != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("f.extension = $%d", idx))
		args = append(args, *filters.Extension)
		idx++
	}
	if filters.MimeType != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("f.mime_type = $%d", idx))
		args = append(args, *filters.MimeType)
		idx++
	}
	if filters.StorageProvider != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("f.storage_provider = $%d", idx))
		args = append(args, *filters.StorageProvider)
		idx++
	}
	if filters.UploadedBy != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("f.uploaded_by = $%d", idx))
		args = append(args, *filters.UploadedBy)
		idx++
	}

	// Entity filters
	if filters.EntityType != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("ef.entity_type = $%d", idx))
		args = append(args, *filters.EntityType)
		idx++
	}
	if filters.EntityID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("ef.entity_id = $%d", idx))
		args = append(args, *filters.EntityID)
		idx++
	}
	if filters.UsageType != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("ef.usage_type = $%d", idx))
		args = append(args, *filters.UsageType)
		idx++
	}
	if filters.IsPrimary != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("ef.is_primary = $%d", idx))
		args = append(args, *filters.IsPrimary)
		idx++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*) 
		FROM files f
		LEFT JOIN entity_files ef ON f.id = ef.file_id` + whereClause
	var total int
	err := pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, 0, err
	}

	// Calculate total pages
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	// Get files
	query := `
		SELECT 
			f.id, f.uid, f.original_name, f.stored_name, f.extension, 
			f.mime_type, f.size_bytes, f.storage_provider, f.storage_bucket, 
			f.storage_path, f.checksum_sha256, f.uploaded_by, f.uploaded_at,
			COALESCE(ef.entity_type, ''), COALESCE(ef.entity_id, 0), 
			COALESCE(ef.usage_type, ''), COALESCE(ef.display_order, 0), COALESCE(ef.is_primary, false)
		FROM files f
		LEFT JOIN entity_files ef ON f.id = ef.file_id
	` + whereClause + fmt.Sprintf(" ORDER BY f.uploaded_at DESC LIMIT $%d OFFSET $%d", idx, idx+1)

	args = append(args, pageSize, offset)

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var out []FileWithEntityRow
	for rows.Next() {
		var r FileWithEntityRow
		if err := rows.Scan(
			&r.ID, &r.UID, &r.OriginalName, &r.StoredName, &r.Extension,
			&r.MimeType, &r.SizeBytes, &r.StorageProvider, &r.StorageBucket,
			&r.StoragePath, &r.ChecksumSHA256, &r.UploadedBy, &r.UploadedAt,
			&r.EntityType, &r.EntityID, &r.UsageType, &r.DisplayOrder, &r.IsPrimary,
		); err != nil {
			return nil, 0, 0, err
		}
		out = append(out, r)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	return out, page, totalPages, nil
}

// GetEntityIDByTypeAndUID resolves entity UID to entity ID based on entity type
func GetEntityIDByTypeAndUID(pool *pgxpool.Pool, entityType, entityUID string) (int64, error) {
	ctx := context.Background()
	var id int64
	var query string

	switch entityType {
	case "machine":
		query = `SELECT id FROM machines WHERE uid = $1`
	case "spare_part":
		query = `SELECT id FROM spare_parts WHERE uid = $1`
	case "vendor":
		query = `SELECT id FROM vendors WHERE uid = $1`
	default:
		return 0, fmt.Errorf("unsupported entity type: %s", entityType)
	}

	err := pool.QueryRow(ctx, query, entityUID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
