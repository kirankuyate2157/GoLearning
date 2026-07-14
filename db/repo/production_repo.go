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

type ProductionOrderRow struct {
	ID           int
	UID          pgtype.UUID
	OrderNumber  string
	CustomerName pgtype.Text
	Reference    pgtype.Text
	OrderDate    pgtype.Date
	RequiredDate pgtype.Date
	Priority     pgtype.Text
	Status       pgtype.Text
	Notes        pgtype.Text
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
}

type JobRow struct {
	ID             int
	UID            pgtype.UUID
	OrderUID       pgtype.UUID
	OrderNumber    pgtype.Text
	MachineUID     pgtype.UUID
	MachineName    pgtype.Text
	JobNumber      string
	PartName       pgtype.Text
	Quantity       float64
	Priority       pgtype.Text
	Status         pgtype.Text
	ScheduledStart pgtype.Date
	ScheduledEnd   pgtype.Date
	Notes          pgtype.Text
	CreatedAt      pgtype.Timestamp
	UpdatedAt      pgtype.Timestamp
}

func ListProductionOrders(pool *pgxpool.Pool) ([]ProductionOrderRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT id, uid, order_number, customer_name, reference, order_date, required_date, priority, status, notes, created_at, updated_at FROM production_orders ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ProductionOrderRow
	for rows.Next() {
		var r ProductionOrderRow
		if err := rows.Scan(&r.ID, &r.UID, &r.OrderNumber, &r.CustomerName, &r.Reference, &r.OrderDate, &r.RequiredDate, &r.Priority, &r.Status, &r.Notes, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func GetProductionOrderByUID(pool *pgxpool.Pool, uid string) (ProductionOrderRow, error) {
	ctx := context.Background()
	var r ProductionOrderRow
	err := pool.QueryRow(ctx, `SELECT id, uid, order_number, customer_name, reference, order_date, required_date, priority, status, notes, created_at, updated_at FROM production_orders WHERE uid = $1`, uid).Scan(&r.ID, &r.UID, &r.OrderNumber, &r.CustomerName, &r.Reference, &r.OrderDate, &r.RequiredDate, &r.Priority, &r.Status, &r.Notes, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func CreateProductionOrder(pool *pgxpool.Pool, req models.ProductionOrderCreateRequest) (string, error) {
	ctx := context.Background()
	priority := strings.TrimSpace(req.Priority)
	if priority == "" {
		priority = "medium"
	}
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "draft"
	}

	orderDate, err := nullableDate(req.OrderDate)
	if err != nil {
		return "", err
	}
	requiredDate, err := nullableDate(req.RequiredDate)
	if err != nil {
		return "", err
	}

	var uid string
	err = pool.QueryRow(ctx, `INSERT INTO production_orders (order_number, customer_name, reference, order_date, required_date, priority, status, notes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING uid`, strings.TrimSpace(req.OrderNumber), strings.TrimSpace(req.CustomerName), strings.TrimSpace(req.Reference), orderDate, requiredDate, priority, status, strings.TrimSpace(req.Notes)).Scan(&uid)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func UpdateProductionOrder(pool *pgxpool.Pool, uid string, req models.ProductionOrderUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	if req.OrderNumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("order_number = $%d", idx))
		args = append(args, strings.TrimSpace(*req.OrderNumber))
		idx++
	}
	if req.CustomerName != nil {
		setClauses = append(setClauses, fmt.Sprintf("customer_name = $%d", idx))
		args = append(args, strings.TrimSpace(*req.CustomerName))
		idx++
	}
	if req.Reference != nil {
		setClauses = append(setClauses, fmt.Sprintf("reference = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Reference))
		idx++
	}
	if req.OrderDate != nil {
		dateArg, err := nullableDate(*req.OrderDate)
		if err != nil {
			return err
		}
		setClauses = append(setClauses, fmt.Sprintf("order_date = $%d", idx))
		args = append(args, dateArg)
		idx++
	}
	if req.RequiredDate != nil {
		dateArg, err := nullableDate(*req.RequiredDate)
		if err != nil {
			return err
		}
		setClauses = append(setClauses, fmt.Sprintf("required_date = $%d", idx))
		args = append(args, dateArg)
		idx++
	}
	if req.Priority != nil {
		setClauses = append(setClauses, fmt.Sprintf("priority = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Priority))
		idx++
	}
	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Status))
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
	query := fmt.Sprintf("UPDATE production_orders SET %s, updated_at = NOW() WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteProductionOrder(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM production_orders WHERE uid = $1`, uid)
	return err
}

func ListJobs(pool *pgxpool.Pool) ([]JobRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT j.id, j.uid, o.uid AS order_uid, o.order_number, m.uid AS machine_uid, m.name AS machine_name, j.job_number, j.part_name, j.quantity::float8 AS quantity, j.priority, j.status, j.scheduled_start, j.scheduled_end, j.notes, j.created_at, j.updated_at FROM jobs j LEFT JOIN production_orders o ON j.order_id = o.id LEFT JOIN machines m ON j.machine_id = m.id ORDER BY j.id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []JobRow
	for rows.Next() {
		var r JobRow
		if err := rows.Scan(&r.ID, &r.UID, &r.OrderUID, &r.OrderNumber, &r.MachineUID, &r.MachineName, &r.JobNumber, &r.PartName, &r.Quantity, &r.Priority, &r.Status, &r.ScheduledStart, &r.ScheduledEnd, &r.Notes, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func ListJobsByMachineUID(pool *pgxpool.Pool, machineUID string) ([]JobRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT j.id, j.uid, o.uid AS order_uid, o.order_number, m.uid AS machine_uid, m.name AS machine_name, j.job_number, j.part_name, j.quantity::float8 AS quantity, j.priority, j.status, j.scheduled_start, j.scheduled_end, j.notes, j.created_at, j.updated_at FROM jobs j LEFT JOIN production_orders o ON j.order_id = o.id JOIN machines m ON j.machine_id = m.id WHERE m.uid = $1 ORDER BY j.scheduled_start NULLS LAST, j.id ASC`, machineUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []JobRow
	for rows.Next() {
		var r JobRow
		if err := rows.Scan(&r.ID, &r.UID, &r.OrderUID, &r.OrderNumber, &r.MachineUID, &r.MachineName, &r.JobNumber, &r.PartName, &r.Quantity, &r.Priority, &r.Status, &r.ScheduledStart, &r.ScheduledEnd, &r.Notes, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func ListJobsByProductionOrderUID(pool *pgxpool.Pool, orderUID string) ([]JobRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT j.id, j.uid, o.uid AS order_uid, o.order_number, m.uid AS machine_uid, m.name AS machine_name, j.job_number, j.part_name, j.quantity::float8 AS quantity, j.priority, j.status, j.scheduled_start, j.scheduled_end, j.notes, j.created_at, j.updated_at FROM jobs j JOIN production_orders o ON j.order_id = o.id LEFT JOIN machines m ON j.machine_id = m.id WHERE o.uid = $1 ORDER BY j.scheduled_start NULLS LAST, j.id ASC`, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []JobRow
	for rows.Next() {
		var r JobRow
		if err := rows.Scan(&r.ID, &r.UID, &r.OrderUID, &r.OrderNumber, &r.MachineUID, &r.MachineName, &r.JobNumber, &r.PartName, &r.Quantity, &r.Priority, &r.Status, &r.ScheduledStart, &r.ScheduledEnd, &r.Notes, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func GetJobByUID(pool *pgxpool.Pool, uid string) (JobRow, error) {
	ctx := context.Background()
	var r JobRow
	err := pool.QueryRow(ctx, `SELECT j.id, j.uid, o.uid AS order_uid, o.order_number, m.uid AS machine_uid, m.name AS machine_name, j.job_number, j.part_name, j.quantity::float8 AS quantity, j.priority, j.status, j.scheduled_start, j.scheduled_end, j.notes, j.created_at, j.updated_at FROM jobs j LEFT JOIN production_orders o ON j.order_id = o.id LEFT JOIN machines m ON j.machine_id = m.id WHERE j.uid = $1`, uid).Scan(&r.ID, &r.UID, &r.OrderUID, &r.OrderNumber, &r.MachineUID, &r.MachineName, &r.JobNumber, &r.PartName, &r.Quantity, &r.Priority, &r.Status, &r.ScheduledStart, &r.ScheduledEnd, &r.Notes, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func CreateJob(pool *pgxpool.Pool, req models.JobCreateRequest, orderID *int, machineID *int) (string, error) {
	ctx := context.Background()
	priority := strings.TrimSpace(req.Priority)
	if priority == "" {
		priority = "medium"
	}
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "planned"
	}

	startDate, err := nullableDate(req.ScheduledStart)
	if err != nil {
		return "", err
	}
	endDate, err := nullableDate(req.ScheduledEnd)
	if err != nil {
		return "", err
	}

	var uid string
	err = pool.QueryRow(ctx, `INSERT INTO jobs (order_id, machine_id, job_number, part_name, quantity, priority, status, scheduled_start, scheduled_end, notes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING uid`, orderID, machineID, strings.TrimSpace(req.JobNumber), strings.TrimSpace(req.PartName), req.Quantity, priority, status, startDate, endDate, strings.TrimSpace(req.Notes)).Scan(&uid)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func UpdateJob(pool *pgxpool.Pool, uid string, req models.JobUpdateRequest, orderID *int, machineID *int) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	if req.OrderUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("order_id = $%d", idx))
		args = append(args, orderID)
		idx++
	}
	if req.MachineUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("machine_id = $%d", idx))
		args = append(args, machineID)
		idx++
	}
	if req.JobNumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("job_number = $%d", idx))
		args = append(args, strings.TrimSpace(*req.JobNumber))
		idx++
	}
	if req.PartName != nil {
		setClauses = append(setClauses, fmt.Sprintf("part_name = $%d", idx))
		args = append(args, strings.TrimSpace(*req.PartName))
		idx++
	}
	if req.Quantity != nil {
		setClauses = append(setClauses, fmt.Sprintf("quantity = $%d", idx))
		args = append(args, *req.Quantity)
		idx++
	}
	if req.Priority != nil {
		setClauses = append(setClauses, fmt.Sprintf("priority = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Priority))
		idx++
	}
	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Status))
		idx++
	}
	if req.ScheduledStart != nil {
		dateArg, err := nullableDate(*req.ScheduledStart)
		if err != nil {
			return err
		}
		setClauses = append(setClauses, fmt.Sprintf("scheduled_start = $%d", idx))
		args = append(args, dateArg)
		idx++
	}
	if req.ScheduledEnd != nil {
		dateArg, err := nullableDate(*req.ScheduledEnd)
		if err != nil {
			return err
		}
		setClauses = append(setClauses, fmt.Sprintf("scheduled_end = $%d", idx))
		args = append(args, dateArg)
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
	query := fmt.Sprintf("UPDATE jobs SET %s, updated_at = NOW() WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteJob(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM jobs WHERE uid = $1`, uid)
	return err
}

func GetProductionOrderByUIDFromID(pool *pgxpool.Pool, id int32) (string, error) {
	ctx := context.Background()
	var uid string
	err := pool.QueryRow(ctx, `SELECT uid FROM production_orders WHERE id = $1`, id).Scan(&uid)
	return uid, err
}

func GetJobByUIDFromID(pool *pgxpool.Pool, id int32) (string, error) {
	ctx := context.Background()
	var uid string
	err := pool.QueryRow(ctx, `SELECT uid FROM jobs WHERE id = $1`, id).Scan(&uid)
	return uid, err
}

func GetMachineUIDFromID(pool *pgxpool.Pool, id int) (string, error) {
	ctx := context.Background()
	var uid string
	err := pool.QueryRow(ctx, `SELECT uid FROM machines WHERE id = $1`, id).Scan(&uid)
	return uid, err
}

func nullableDate(value string) (interface{}, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", trimmed)
	if err != nil {
		return nil, fmt.Errorf("invalid date %q: %w", value, err)
	}
	return parsed, nil
}
