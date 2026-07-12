package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FeedAlertRow struct {
	ID             int
	UID            string
	AlertUID       string
	LowAlertValue  float64
	HighAlertValue float64
	FeedUID        string
}

type FeedAlertCreateRequest struct {
	AlertID        int     `json:"alert_id" validate:"required"`
	LowAlertValue  float64 `json:"low_alert_value" validate:"required"`
	HighAlertValue float64 `json:"high_alert_value" validate:"required"`
	FeedID         int     `json:"feed_id" validate:"required"`
}

type FeedAlertUpdateRequest struct {
	AlertID        *int     `json:"alert_id,omitempty"`
	LowAlertValue  *float64 `json:"low_alert_value,omitempty"`
	HighAlertValue *float64 `json:"high_alert_value,omitempty"`
	FeedID         *int     `json:"feed_id" validate:"required"`
}

type FeedAlertsParams struct {
	UID     *string `json:"uid" example:"025dd763-4068-49e8-b62e-b82006281736"`
	AlertID *int    `json:"alert_uid"`
	FeedID  *int    `json:"feed_uid"`
}

func ListFeedAlerts(
	pool *pgxpool.Pool,
	qp FeedAlertsParams,
) ([]FeedAlertRow, error) {
	ctx := context.Background()
	var filters []string
	var filterValues []any
	idx := 1

	query := `SELECT fa.id, 
					fa.uid,
					a.uid as alert_uid, 
					fa.low_alert_value, 
					fa.high_alert_value,
					f.uid as feed_uid 
			FROM feed_alerts fa
			JOIN alerts a ON fa.alert_id = a.id 
			JOIN feeds f ON fa.feed_id = f.id
			`

	if qp.FeedID != nil {
		filters = append(filters, fmt.Sprintf(" f.id = $%d ", idx))
		filterValues = append(filterValues, &qp.FeedID)
		idx++
	}

	if qp.AlertID != nil {
		filters = append(filters, fmt.Sprintf(" a.id = $%d ", idx))
		filterValues = append(filterValues, &qp.AlertID)
		idx++
	}

	if qp.UID != nil {
		filters = append(filters, fmt.Sprintf(" fa.uid = $%d ", idx))
		filterValues = append(filterValues, &qp.UID)
		idx++
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += " ORDER BY fa.id ASC"

	rows, err := pool.Query(ctx, query, filterValues...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []FeedAlertRow

	for rows.Next() {
		var r FeedAlertRow
		if err := rows.Scan(&r.ID,
			&r.UID,
			&r.AlertUID,
			&r.LowAlertValue,
			&r.HighAlertValue,
			&r.FeedUID,
		); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func GetFeedAlertByUID(pool *pgxpool.Pool, uid string) (FeedAlertRow, error) {
	ctx := context.Background()
	query := `SELECT fa.id, 
					fa.uid,
					a.uid as alert_uid, 
					fa.low_alert_value, 
					fa.high_alert_value,
					f.uid as feed_uid
			FROM feed_alerts fa
			JOIN alerts a ON fa.alert_id = a.id 
			JOIN feeds f ON fa.feed_id = f.id
			WHERE fa.uid = $1`
	var r FeedAlertRow
	err := pool.QueryRow(ctx, query, uid).Scan(&r.ID,
		&r.UID,
		&r.AlertUID,
		&r.LowAlertValue,
		&r.HighAlertValue,
		&r.FeedUID,
	)
	return r, err
}

func CreateFeedAlert(pool *pgxpool.Pool, req FeedAlertCreateRequest) (string, error) {
	ctx := context.Background()
	var newUID string
	err := pool.QueryRow(ctx,
		`INSERT INTO feed_alerts 
		(	alert_id, 
			low_alert_value, 
			high_alert_value,
			feed_id
		) 
		VALUES ($1, $2, $3, $4) 
		RETURNING uid`,
		req.AlertID,
		req.LowAlertValue,
		req.HighAlertValue,
		req.FeedID,
	).Scan(&newUID)
	if err != nil {
		return "", err
	}
	return newUID, nil
}

func UpdateFeedAlert(pool *pgxpool.Pool, uid string, req FeedAlertUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1
	if req.AlertID != nil {
		setClauses = append(setClauses, fmt.Sprintf("alert_id = $%d", idx))
		args = append(args, *req.AlertID)
		idx++
	}
	if req.LowAlertValue != nil {
		setClauses = append(setClauses, fmt.Sprintf("low_alert_value = $%d", idx))
		args = append(args, *req.LowAlertValue)
		idx++
	}
	if req.HighAlertValue != nil {
		setClauses = append(setClauses, fmt.Sprintf("high_alert_value = $%d", idx))
		args = append(args, *req.HighAlertValue)
		idx++
	}
	if req.FeedID != nil {
		setClauses = append(setClauses, fmt.Sprintf("feed_id = $%d", idx))
		args = append(args, *req.FeedID)
		idx++
	}
	if len(setClauses) == 0 {
		return nil
	}
	args = append(args, uid)
	query := fmt.Sprintf("UPDATE feed_alerts SET %s WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteFeedAlert(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM feed_alerts WHERE uid = $1`, uid)
	return err
}
