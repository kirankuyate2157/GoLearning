package controllers

import (
	dbrepo "app/db/repo"
	"app/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func GetProductionOrders(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.ProductionOrder{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	rows, err := dbrepo.ListProductionOrders(pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query production orders: %v", err)})
	}

	out := make([]models.ProductionOrder, 0, len(rows))
	for _, r := range rows {
		order := models.ProductionOrder{UID: uuidString(r.UID), OrderNumber: r.OrderNumber}
		if r.CustomerName.Valid {
			order.CustomerName = r.CustomerName.String
		}
		if r.Reference.Valid {
			order.Reference = r.Reference.String
		}
		if r.OrderDate.Valid {
			t := r.OrderDate.Time
			order.OrderDate = &t
		}
		if r.RequiredDate.Valid {
			t := r.RequiredDate.Time
			order.RequiredDate = &t
		}
		if r.Priority.Valid {
			order.Priority = r.Priority.String
		}
		if r.Status.Valid {
			order.Status = r.Status.String
		}
		if r.Notes.Valid {
			order.Notes = r.Notes.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			order.CreatedAt = &t
		}
		if r.UpdatedAt.Valid {
			t := r.UpdatedAt.Time
			order.UpdatedAt = &t
		}
		out = append(out, order)
	}
	return c.JSON(http.StatusOK, out)
}

func GetProductionOrder(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	r, err := dbrepo.GetProductionOrderByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query production order: %v", err)})
	}

	order := models.ProductionOrder{UID: uuidString(r.UID), OrderNumber: r.OrderNumber}
	if r.CustomerName.Valid {
		order.CustomerName = r.CustomerName.String
	}
	if r.Reference.Valid {
		order.Reference = r.Reference.String
	}
	if r.OrderDate.Valid {
		t := r.OrderDate.Time
		order.OrderDate = &t
	}
	if r.RequiredDate.Valid {
		t := r.RequiredDate.Time
		order.RequiredDate = &t
	}
	if r.Priority.Valid {
		order.Priority = r.Priority.String
	}
	if r.Status.Valid {
		order.Status = r.Status.String
	}
	if r.Notes.Valid {
		order.Notes = r.Notes.String
	}
	if r.CreatedAt.Valid {
		t := r.CreatedAt.Time
		order.CreatedAt = &t
	}
	if r.UpdatedAt.Valid {
		t := r.UpdatedAt.Time
		order.UpdatedAt = &t
	}
	return c.JSON(http.StatusOK, order)
}

func GetProductionMetrics(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, models.ProductionMetrics{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	metrics := models.ProductionMetrics{}
	var err error

	metrics.ProductionOrderStatusCounts, err = dbrepo.GetProductionOrderStatusCounts(pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query production order status counts: %v", err)})
	}
	metrics.JobStatusCounts, err = dbrepo.GetJobStatusCounts(pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query job status counts: %v", err)})
	}
	metrics.WIPItemStatusCounts, err = dbrepo.GetWIPItemStatusCounts(pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query wip item status counts: %v", err)})
	}

	for _, v := range metrics.ProductionOrderStatusCounts {
		metrics.ProductionOrderCount += v
	}
	for _, v := range metrics.JobStatusCounts {
		metrics.JobCount += v
	}
	for _, v := range metrics.WIPItemStatusCounts {
		metrics.WIPItemCount += v
	}

	metrics.AverageJobCompletionHours, err = dbrepo.GetAverageCompletedJobDurationHours(pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query average job duration: %v", err)})
	}

	return c.JSON(http.StatusOK, metrics)
}

func CreateProductionOrder(c echo.Context) error {
	var req models.ProductionOrderCreateRequest
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

	uid, err := dbrepo.CreateProductionOrder(pool, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": uid})
}

func UpdateProductionOrder(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.ProductionOrderUpdateRequest
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

	if err := dbrepo.UpdateProductionOrder(pool, uid, req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

func DeleteProductionOrder(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if err := dbrepo.DeleteProductionOrder(pool, uid); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete production order: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}

func GetJobs(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.Job{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	rows, err := dbrepo.ListJobs(pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query jobs: %v", err)})
	}

	out := make([]models.Job, 0, len(rows))
	for _, r := range rows {
		job := models.Job{UID: uuidString(r.UID), JobNumber: r.JobNumber, Quantity: r.Quantity}
		if r.OrderUID.Valid {
			job.OrderUID = r.OrderUID.String()
		}
		if r.OrderNumber.Valid {
			job.OrderNumber = r.OrderNumber.String
		}
		if r.MachineUID.Valid {
			job.MachineUID = r.MachineUID.String()
		}
		if r.MachineName.Valid {
			job.MachineName = r.MachineName.String
		}
		if r.PartName.Valid {
			job.PartName = r.PartName.String
		}
		if r.Priority.Valid {
			job.Priority = r.Priority.String
		}
		if r.Status.Valid {
			job.Status = r.Status.String
		}
		if r.ScheduledStart.Valid {
			t := r.ScheduledStart.Time
			job.ScheduledStart = &t
		}
		if r.ScheduledEnd.Valid {
			t := r.ScheduledEnd.Time
			job.ScheduledEnd = &t
		}
		if r.Notes.Valid {
			job.Notes = r.Notes.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			job.CreatedAt = &t
		}
		if r.UpdatedAt.Valid {
			t := r.UpdatedAt.Time
			job.UpdatedAt = &t
		}
		out = append(out, job)
	}
	return c.JSON(http.StatusOK, out)
}

func GetMachineJobs(c echo.Context) error {
	machineUID := strings.TrimSpace(c.Param("uid"))
	if machineUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid machine uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.Job{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	rows, err := dbrepo.ListJobsByMachineUID(pool, machineUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine jobs: %v", err)})
	}

	out := make([]models.Job, 0, len(rows))
	for _, r := range rows {
		job := models.Job{UID: uuidString(r.UID), JobNumber: r.JobNumber, Quantity: r.Quantity}
		if r.OrderUID.Valid {
			job.OrderUID = r.OrderUID.String()
		}
		if r.OrderNumber.Valid {
			job.OrderNumber = r.OrderNumber.String
		}
		if r.MachineUID.Valid {
			job.MachineUID = r.MachineUID.String()
		}
		if r.MachineName.Valid {
			job.MachineName = r.MachineName.String
		}
		if r.PartName.Valid {
			job.PartName = r.PartName.String
		}
		if r.Priority.Valid {
			job.Priority = r.Priority.String
		}
		if r.Status.Valid {
			job.Status = r.Status.String
		}
		if r.ScheduledStart.Valid {
			t := r.ScheduledStart.Time
			job.ScheduledStart = &t
		}
		if r.ScheduledEnd.Valid {
			t := r.ScheduledEnd.Time
			job.ScheduledEnd = &t
		}
		if r.Notes.Valid {
			job.Notes = r.Notes.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			job.CreatedAt = &t
		}
		if r.UpdatedAt.Valid {
			t := r.UpdatedAt.Time
			job.UpdatedAt = &t
		}
		out = append(out, job)
	}
	return c.JSON(http.StatusOK, out)
}

func GetProductionOrderJobs(c echo.Context) error {
	orderUID := strings.TrimSpace(c.Param("uid"))
	if orderUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid production order uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.Job{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	rows, err := dbrepo.ListJobsByProductionOrderUID(pool, orderUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query production order jobs: %v", err)})
	}

	out := make([]models.Job, 0, len(rows))
	for _, r := range rows {
		job := models.Job{UID: uuidString(r.UID), JobNumber: r.JobNumber, Quantity: r.Quantity}
		if r.OrderUID.Valid {
			job.OrderUID = r.OrderUID.String()
		}
		if r.OrderNumber.Valid {
			job.OrderNumber = r.OrderNumber.String
		}
		if r.MachineUID.Valid {
			job.MachineUID = r.MachineUID.String()
		}
		if r.MachineName.Valid {
			job.MachineName = r.MachineName.String
		}
		if r.PartName.Valid {
			job.PartName = r.PartName.String
		}
		if r.Priority.Valid {
			job.Priority = r.Priority.String
		}
		if r.Status.Valid {
			job.Status = r.Status.String
		}
		if r.ScheduledStart.Valid {
			t := r.ScheduledStart.Time
			job.ScheduledStart = &t
		}
		if r.ScheduledEnd.Valid {
			t := r.ScheduledEnd.Time
			job.ScheduledEnd = &t
		}
		if r.Notes.Valid {
			job.Notes = r.Notes.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			job.CreatedAt = &t
		}
		if r.UpdatedAt.Valid {
			t := r.UpdatedAt.Time
			job.UpdatedAt = &t
		}
		out = append(out, job)
	}
	return c.JSON(http.StatusOK, out)
}

func GetJob(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	r, err := dbrepo.GetJobByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query job: %v", err)})
	}

	job := models.Job{UID: uuidString(r.UID), JobNumber: r.JobNumber, Quantity: r.Quantity}
	if r.OrderUID.Valid {
		job.OrderUID = r.OrderUID.String()
	}
	if r.OrderNumber.Valid {
		job.OrderNumber = r.OrderNumber.String
	}
	if r.MachineUID.Valid {
		job.MachineUID = r.MachineUID.String()
	}
	if r.MachineName.Valid {
		job.MachineName = r.MachineName.String
	}
	if r.PartName.Valid {
		job.PartName = r.PartName.String
	}
	if r.Priority.Valid {
		job.Priority = r.Priority.String
	}
	if r.Status.Valid {
		job.Status = r.Status.String
	}
	if r.ScheduledStart.Valid {
		t := r.ScheduledStart.Time
		job.ScheduledStart = &t
	}
	if r.ScheduledEnd.Valid {
		t := r.ScheduledEnd.Time
		job.ScheduledEnd = &t
	}
	if r.Notes.Valid {
		job.Notes = r.Notes.String
	}
	if r.CreatedAt.Valid {
		t := r.CreatedAt.Time
		job.CreatedAt = &t
	}
	if r.UpdatedAt.Valid {
		t := r.UpdatedAt.Time
		job.UpdatedAt = &t
	}
	return c.JSON(http.StatusOK, job)
}

func GetJobExecutionLogs(c echo.Context) error {
	jobUID := strings.TrimSpace(c.Param("uid"))
	if jobUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid job uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.JobExecutionLog{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	rows, err := dbrepo.ListJobExecutionLogsByJobUID(pool, jobUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "job not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query job execution logs: %v", err)})
	}

	out := make([]models.JobExecutionLog, 0, len(rows))
	for _, r := range rows {
		logEntry := models.JobExecutionLog{UID: uuidString(r.UID), Action: r.Action, Status: r.Status.String, OperatorName: r.OperatorName.String, Note: r.Note.String}
		if r.JobUID.Valid {
			logEntry.JobUID = r.JobUID.String()
		}
		if r.MachineUID.Valid {
			logEntry.MachineUID = r.MachineUID.String()
		}
		if r.LoggedAt.Valid {
			t := r.LoggedAt.Time
			logEntry.LoggedAt = &t
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			logEntry.CreatedAt = &t
		}
		if r.UpdatedAt.Valid {
			t := r.UpdatedAt.Time
			logEntry.UpdatedAt = &t
		}
		out = append(out, logEntry)
	}
	return c.JSON(http.StatusOK, out)
}

func CreateJobExecutionLog(c echo.Context) error {
	jobUID := strings.TrimSpace(c.Param("uid"))
	if jobUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid job uid"})
	}
	var req models.JobExecutionLogCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if strings.TrimSpace(req.Action) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "action is required"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	uid, err := dbrepo.CreateJobExecutionLog(pool, jobUID, req)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "job not found"})
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": uid})
}

func CreateJob(c echo.Context) error {
	var req models.JobCreateRequest
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

	orderRow, err := dbrepo.GetProductionOrderByUID(pool, req.OrderUID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "order uid not found"})
	}
	var machineID *int
	if strings.TrimSpace(req.MachineUID) != "" {
		machineRow, err := dbrepo.GetMachineByUID(pool, req.MachineUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine uid not found"})
		}
		machineID = &machineRow.ID
	}
	orderID := &orderRow.ID

	uid, err := dbrepo.CreateJob(pool, req, orderID, machineID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": uid})
}

func UpdateJob(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.JobUpdateRequest
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

	var orderID *int
	var machineID *int
	if req.OrderUID != nil && strings.TrimSpace(*req.OrderUID) != "" {
		orderRow, err := dbrepo.GetProductionOrderByUID(pool, *req.OrderUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "order uid not found"})
		}
		orderID = &orderRow.ID
	}
	if req.MachineUID != nil && strings.TrimSpace(*req.MachineUID) != "" {
		machineRow, err := dbrepo.GetMachineByUID(pool, *req.MachineUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine uid not found"})
		}
		machineID = &machineRow.ID
	}

	if err := dbrepo.UpdateJob(pool, uid, req, orderID, machineID); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

func DeleteJob(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if err := dbrepo.DeleteJob(pool, uid); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete job: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}

func uuidString(v interface{ String() string }) string {
	return v.String()
}
