package repo

import (
	"app/models"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QualityInspectionRow struct {
	ID            int
	UID           pgtype.UUID
	JobUID        pgtype.UUID
	WIPItemUID    pgtype.UUID
	InspectorName pgtype.Text
	Status        pgtype.Text
	InspectedAt   pgtype.Timestamp
	Findings      pgtype.Text
	Result        pgtype.Text
	Notes         pgtype.Text
	CreatedAt     pgtype.Timestamp
	UpdatedAt     pgtype.Timestamp
}

type NonConformanceReportRow struct {
	ID               int
	UID              pgtype.UUID
	InspectionUID    pgtype.UUID
	JobUID           pgtype.UUID
	WIPItemUID       pgtype.UUID
	Type             pgtype.Text
	Severity         pgtype.Text
	Description      pgtype.Text
	CorrectiveAction pgtype.Text
	Status           pgtype.Text
	CreatedAt        pgtype.Timestamp
	UpdatedAt        pgtype.Timestamp
}

func ListQualityInspections(pool *pgxpool.Pool) ([]QualityInspectionRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT id, uid, job_uid, wip_item_uid, inspector_name, status, inspected_at, findings, result, notes, created_at, updated_at FROM quality_inspections ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []QualityInspectionRow
	for rows.Next() {
		var r QualityInspectionRow
		if err := rows.Scan(&r.ID, &r.UID, &r.JobUID, &r.WIPItemUID, &r.InspectorName, &r.Status, &r.InspectedAt, &r.Findings, &r.Result, &r.Notes, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func GetQualityInspectionByUID(pool *pgxpool.Pool, uid string) (QualityInspectionRow, error) {
	ctx := context.Background()
	var r QualityInspectionRow
	err := pool.QueryRow(ctx, `SELECT id, uid, job_uid, wip_item_uid, inspector_name, status, inspected_at, findings, result, notes, created_at, updated_at FROM quality_inspections WHERE uid = $1`, uid).Scan(&r.ID, &r.UID, &r.JobUID, &r.WIPItemUID, &r.InspectorName, &r.Status, &r.InspectedAt, &r.Findings, &r.Result, &r.Notes, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func CreateQualityInspection(pool *pgxpool.Pool, req models.QualityInspectionCreateRequest) (string, error) {
	ctx := context.Background()
	var uid string
	if req.Status == "" {
		req.Status = "pending"
	}
	inspectedAt, err := nullableTimestamp(req.InspectedAt)
	if err != nil {
		return "", err
	}
	err = pool.QueryRow(ctx, `INSERT INTO quality_inspections (job_uid, wip_item_uid, inspector_name, status, inspected_at, findings, result, notes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING uid`, nullableUUID(req.JobUID), nullableUUID(req.WIPItemUID), strings.TrimSpace(req.InspectorName), strings.TrimSpace(req.Status), inspectedAt, strings.TrimSpace(req.Findings), strings.TrimSpace(req.Result), strings.TrimSpace(req.Notes)).Scan(&uid)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func UpdateQualityInspection(pool *pgxpool.Pool, uid string, req models.QualityInspectionUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	if req.JobUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("job_uid = $%d", idx))
		args = append(args, nullableUUID(strings.TrimSpace(*req.JobUID)))
		idx++
	}
	if req.WIPItemUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("wip_item_uid = $%d", idx))
		args = append(args, nullableUUID(strings.TrimSpace(*req.WIPItemUID)))
		idx++
	}
	if req.InspectorName != nil {
		setClauses = append(setClauses, fmt.Sprintf("inspector_name = $%d", idx))
		args = append(args, strings.TrimSpace(*req.InspectorName))
		idx++
	}
	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Status))
		idx++
	}
	if req.InspectedAt != nil {
		inspectedAt, err := nullableTimestamp(strings.TrimSpace(*req.InspectedAt))
		if err != nil {
			return err
		}
		setClauses = append(setClauses, fmt.Sprintf("inspected_at = $%d", idx))
		args = append(args, inspectedAt)
		idx++
	}
	if req.Findings != nil {
		setClauses = append(setClauses, fmt.Sprintf("findings = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Findings))
		idx++
	}
	if req.Result != nil {
		setClauses = append(setClauses, fmt.Sprintf("result = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Result))
		idx++
	}
	if req.Notes != nil {
		setClauses = append(setClauses, fmt.Sprintf("notes = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Notes))
		idx++
	}
	if len(setClauses) == 0 {
		return nil
	}
	args = append(args, uid)
	query := fmt.Sprintf("UPDATE quality_inspections SET %s, updated_at = NOW() WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteQualityInspection(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM quality_inspections WHERE uid = $1`, uid)
	return err
}

func ListNonConformanceReportsByInspectionUID(pool *pgxpool.Pool, inspectionUID string) ([]NonConformanceReportRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT id, uid, inspection_uid, job_uid, wip_item_uid, type, severity, description, corrective_action, status, created_at, updated_at FROM non_conformance_reports WHERE inspection_uid = $1 ORDER BY id ASC`, inspectionUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []NonConformanceReportRow
	for rows.Next() {
		var r NonConformanceReportRow
		if err := rows.Scan(&r.ID, &r.UID, &r.InspectionUID, &r.JobUID, &r.WIPItemUID, &r.Type, &r.Severity, &r.Description, &r.CorrectiveAction, &r.Status, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func CreateNonConformanceReport(pool *pgxpool.Pool, req models.NonConformanceReportCreateRequest) (string, error) {
	ctx := context.Background()
	var uid string
	if req.Status == "" {
		req.Status = "open"
	}
	err := pool.QueryRow(ctx, `INSERT INTO non_conformance_reports (inspection_uid, job_uid, wip_item_uid, type, severity, description, corrective_action, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING uid`, strings.TrimSpace(req.InspectionUID), nullableUUID(req.JobUID), nullableUUID(req.WIPItemUID), strings.TrimSpace(req.Type), strings.TrimSpace(req.Severity), strings.TrimSpace(req.Description), strings.TrimSpace(req.CorrectiveAction), strings.TrimSpace(req.Status)).Scan(&uid)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func UpdateNonConformanceReport(pool *pgxpool.Pool, uid string, req models.NonConformanceReportUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	if req.JobUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("job_uid = $%d", idx))
		args = append(args, nullableUUID(strings.TrimSpace(*req.JobUID)))
		idx++
	}
	if req.WIPItemUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("wip_item_uid = $%d", idx))
		args = append(args, nullableUUID(strings.TrimSpace(*req.WIPItemUID)))
		idx++
	}
	if req.Type != nil {
		setClauses = append(setClauses, fmt.Sprintf("type = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Type))
		idx++
	}
	if req.Severity != nil {
		setClauses = append(setClauses, fmt.Sprintf("severity = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Severity))
		idx++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Description))
		idx++
	}
	if req.CorrectiveAction != nil {
		setClauses = append(setClauses, fmt.Sprintf("corrective_action = $%d", idx))
		args = append(args, strings.TrimSpace(*req.CorrectiveAction))
		idx++
	}
	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Status))
		idx++
	}
	if len(setClauses) == 0 {
		return nil
	}
	args = append(args, uid)
	query := fmt.Sprintf("UPDATE non_conformance_reports SET %s, updated_at = NOW() WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteNonConformanceReport(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM non_conformance_reports WHERE uid = $1`, uid)
	return err
}

func nullableUUID(value string) interface{} {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func nullableTimestamp(value string) (interface{}, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp %q: %w", value, err)
	}
	return parsed, nil
}
