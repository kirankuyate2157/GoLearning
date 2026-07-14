package controllers

import (
	dbrepo "app/db/repo"
	"app/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func GetQualityInspections(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.QualityInspection{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	rows, err := dbrepo.ListQualityInspections(pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query quality inspections: %v", err)})
	}

	out := make([]models.QualityInspection, 0, len(rows))
	for _, r := range rows {
		inspection := models.QualityInspection{UID: r.UID.String(), InspectorName: r.InspectorName.String, Status: r.Status.String, Findings: r.Findings.String, Result: r.Result.String, Notes: r.Notes.String}
		if r.JobUID.Valid {
			inspection.JobUID = r.JobUID.String()
		}
		if r.WIPItemUID.Valid {
			inspection.WIPItemUID = r.WIPItemUID.String()
		}
		if r.InspectedAt.Valid {
			inspection.InspectedAt = r.InspectedAt.Time.Format(time.RFC3339)
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			inspection.CreatedAt = &t
		}
		if r.UpdatedAt.Valid {
			t := r.UpdatedAt.Time
			inspection.UpdatedAt = &t
		}
		out = append(out, inspection)
	}
	return c.JSON(http.StatusOK, out)
}

func GetQualityInspection(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	r, err := dbrepo.GetQualityInspectionByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query quality inspection: %v", err)})
	}

	inspection := models.QualityInspection{UID: r.UID.String(), InspectorName: r.InspectorName.String, Status: r.Status.String, Findings: r.Findings.String, Result: r.Result.String, Notes: r.Notes.String}
	if r.JobUID.Valid {
		inspection.JobUID = r.JobUID.String()
	}
	if r.WIPItemUID.Valid {
		inspection.WIPItemUID = r.WIPItemUID.String()
	}
	if r.InspectedAt.Valid {
		inspection.InspectedAt = r.InspectedAt.Time.Format(time.RFC3339)
	}
	if r.CreatedAt.Valid {
		t := r.CreatedAt.Time
		inspection.CreatedAt = &t
	}
	if r.UpdatedAt.Valid {
		t := r.UpdatedAt.Time
		inspection.UpdatedAt = &t
	}
	return c.JSON(http.StatusOK, inspection)
}

func CreateQualityInspection(c echo.Context) error {
	var req models.QualityInspectionCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if strings.TrimSpace(req.JobUID) == "" && strings.TrimSpace(req.WIPItemUID) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "job_uid or wip_item_uid is required"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	if strings.TrimSpace(req.JobUID) != "" {
		if _, err := dbrepo.GetJobByUID(pool, strings.TrimSpace(req.JobUID)); err != nil {
			if err == pgx.ErrNoRows {
				return c.JSON(http.StatusNotFound, echo.Map{"error": "job uid not found"})
			}
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("job lookup error: %v", err)})
		}
	}
	if strings.TrimSpace(req.WIPItemUID) != "" {
		if _, err := dbrepo.GetWIPItemByUID(pool, strings.TrimSpace(req.WIPItemUID)); err != nil {
			if err == pgx.ErrNoRows {
				return c.JSON(http.StatusNotFound, echo.Map{"error": "wip item uid not found"})
			}
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("wip lookup error: %v", err)})
		}
	}

	uid, err := dbrepo.CreateQualityInspection(pool, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": uid})
}

func UpdateQualityInspection(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}
	var req models.QualityInspectionUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	if req.JobUID != nil && strings.TrimSpace(*req.JobUID) != "" {
		if _, err := dbrepo.GetJobByUID(pool, strings.TrimSpace(*req.JobUID)); err != nil {
			if err == pgx.ErrNoRows {
				return c.JSON(http.StatusNotFound, echo.Map{"error": "job uid not found"})
			}
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("job lookup error: %v", err)})
		}
	}
	if req.WIPItemUID != nil && strings.TrimSpace(*req.WIPItemUID) != "" {
		if _, err := dbrepo.GetWIPItemByUID(pool, strings.TrimSpace(*req.WIPItemUID)); err != nil {
			if err == pgx.ErrNoRows {
				return c.JSON(http.StatusNotFound, echo.Map{"error": "wip item uid not found"})
			}
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("wip lookup error: %v", err)})
		}
	}

	if err := dbrepo.UpdateQualityInspection(pool, uid, req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

func DeleteQualityInspection(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if err := dbrepo.DeleteQualityInspection(pool, uid); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete quality inspection: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}

func GetQualityInspectionNonConformances(c echo.Context) error {
	inspectionUID := strings.Trim(c.Param("uid"), "/")
	if inspectionUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.NonConformanceReport{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	if _, err := dbrepo.GetQualityInspectionByUID(pool, inspectionUID); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "inspection not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("inspection lookup error: %v", err)})
	}

	rows, err := dbrepo.ListNonConformanceReportsByInspectionUID(pool, inspectionUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query non-conformances: %v", err)})
	}

	out := make([]models.NonConformanceReport, 0, len(rows))
	for _, r := range rows {
		ncr := models.NonConformanceReport{UID: r.UID.String(), Type: r.Type.String, Severity: r.Severity.String, Description: r.Description.String, CorrectiveAction: r.CorrectiveAction.String, Status: r.Status.String}
		if r.InspectionUID.Valid {
			ncr.InspectionUID = r.InspectionUID.String()
		}
		if r.JobUID.Valid {
			ncr.JobUID = r.JobUID.String()
		}
		if r.WIPItemUID.Valid {
			ncr.WIPItemUID = r.WIPItemUID.String()
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			ncr.CreatedAt = &t
		}
		if r.UpdatedAt.Valid {
			t := r.UpdatedAt.Time
			ncr.UpdatedAt = &t
		}
		out = append(out, ncr)
	}
	return c.JSON(http.StatusOK, out)
}

func CreateQualityInspectionNonConformance(c echo.Context) error {
	inspectionUID := strings.Trim(c.Param("uid"), "/")
	if inspectionUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}
	var req models.NonConformanceReportCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if strings.TrimSpace(req.InspectionUID) == "" {
		req.InspectionUID = inspectionUID
	} else if strings.TrimSpace(req.InspectionUID) != inspectionUID {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "inspection uid mismatch"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	if _, err := dbrepo.GetQualityInspectionByUID(pool, inspectionUID); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "inspection not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("inspection lookup error: %v", err)})
	}

	if strings.TrimSpace(req.JobUID) != "" {
		if _, err := dbrepo.GetJobByUID(pool, strings.TrimSpace(req.JobUID)); err != nil {
			if err == pgx.ErrNoRows {
				return c.JSON(http.StatusNotFound, echo.Map{"error": "job uid not found"})
			}
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("job lookup error: %v", err)})
		}
	}
	if strings.TrimSpace(req.WIPItemUID) != "" {
		if _, err := dbrepo.GetWIPItemByUID(pool, strings.TrimSpace(req.WIPItemUID)); err != nil {
			if err == pgx.ErrNoRows {
				return c.JSON(http.StatusNotFound, echo.Map{"error": "wip item uid not found"})
			}
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("wip lookup error: %v", err)})
		}
	}

	uid, err := dbrepo.CreateNonConformanceReport(pool, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": uid})
}

func UpdateQualityInspectionNonConformance(c echo.Context) error {
	ncUID := strings.Trim(c.Param("nc_uid"), "/")
	if ncUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}
	var req models.NonConformanceReportUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	if req.JobUID != nil && strings.TrimSpace(*req.JobUID) != "" {
		if _, err := dbrepo.GetJobByUID(pool, strings.TrimSpace(*req.JobUID)); err != nil {
			if err == pgx.ErrNoRows {
				return c.JSON(http.StatusNotFound, echo.Map{"error": "job uid not found"})
			}
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("job lookup error: %v", err)})
		}
	}
	if req.WIPItemUID != nil && strings.TrimSpace(*req.WIPItemUID) != "" {
		if _, err := dbrepo.GetWIPItemByUID(pool, strings.TrimSpace(*req.WIPItemUID)); err != nil {
			if err == pgx.ErrNoRows {
				return c.JSON(http.StatusNotFound, echo.Map{"error": "wip item uid not found"})
			}
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("wip lookup error: %v", err)})
		}
	}

	if err := dbrepo.UpdateNonConformanceReport(pool, ncUID, req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": ncUID})
}

func DeleteQualityInspectionNonConformance(c echo.Context) error {
	ncUID := strings.Trim(c.Param("nc_uid"), "/")
	if ncUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"uid": ncUID})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if err := dbrepo.DeleteNonConformanceReport(pool, ncUID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete non-conformance report: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}
