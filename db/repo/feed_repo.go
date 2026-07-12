package repo

import (
	"app/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FeedRow struct {
	ID                int
	UID               string
	Name              string
	Code              string
	Description       string
	LowValue          float64
	HighValue         float64
	DisplayUnit       string
	TypeUID           string
	MarginOfError     float64
	ConnectionType    string
	ConnectionTypeUID string
	IngestionType     string
	PollingRate       int
	IsRunning         bool
	JSONData          []map[string]interface{}
}

func ListFeeds(pool *pgxpool.Pool) ([]FeedRow, error) {
	ctx := context.Background()
	query := `SELECT f.id, 
					f.uid, 
					f.name, 
					f.code, 
					f.description, 
					f.low_value, 
					f.high_value, 
					t.uid as display_unit, 
					t1.uid as type_id,
					f.margin_of_error, 
					ct.uid as connection_type_uid, 
					f.ingestion_type, 
					f.polling_rate,
					f.is_running, 
					f.json 
				FROM feeds f
				JOIN feed_connection_types ct ON f.connection_type = ct.id
				JOIN types t ON f.display_unit = t.id
				JOIN types t1 ON f.type_id = t1.id
			ORDER BY id ASC`
	rows, err := pool.Query(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []FeedRow
	for rows.Next() {
		var r FeedRow
		var jsonBytes []byte
		if err := rows.Scan(&r.ID,
			&r.UID,
			&r.Name,
			&r.Code,
			&r.Description,
			&r.LowValue,
			&r.HighValue,
			&r.DisplayUnit,
			&r.TypeUID,
			&r.MarginOfError,
			&r.ConnectionTypeUID,
			&r.IngestionType,
			&r.PollingRate,
			&r.IsRunning,
			&jsonBytes); err != nil {
			return nil, err
		}
		// Parse JSON data

		if len(jsonBytes) > 0 {
			if err := json.Unmarshal(jsonBytes, &r.JSONData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
			}
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func GetFeedByUID(pool *pgxpool.Pool, uid string) (FeedRow, error) {
	ctx := context.Background()
	query := `SELECT f.id, 
					f.uid, 
					f.name, 
					f.code, 
					f.description, 
					f.low_value, 
					f.high_value, 
					t.uid as display_unit, 
					t1.uid as type_id,
					f.margin_of_error, 
					ct.name as connection_type,
					ct.uid as connection_type_uid, 
					f.ingestion_type, 
					f.polling_rate,
					f.is_running, 
					f.json 
				FROM feeds f
				JOIN feed_connection_types ct ON f.connection_type = ct.id
				JOIN types t ON f.display_unit = t.id
				JOIN types t1 ON f.type_id = t1.id
				WHERE f.uid = $1`
	var r FeedRow
	var jsonBytes []byte
	err := pool.QueryRow(ctx, query, uid).Scan(&r.ID,
		&r.UID,
		&r.Name,
		&r.Code,
		&r.Description,
		&r.LowValue,
		&r.HighValue,
		&r.DisplayUnit,
		&r.TypeUID,
		&r.MarginOfError,
		&r.ConnectionType,
		&r.ConnectionTypeUID,
		&r.IngestionType,
		&r.PollingRate,
		&r.IsRunning,
		&jsonBytes)

	if err != nil {
		return r, err
	}

	// Parse JSON data
	if len(jsonBytes) > 0 {
		if err := json.Unmarshal(jsonBytes, &r.JSONData); err != nil {
			return r, fmt.Errorf("failed to unmarshal JSON data: %w", err)
		}
	}

	return r, nil
}

func CreateFeed(
	pool *pgxpool.Pool,
	name string,
	code string,
	description string,
	lowValue *float64,
	highValue *float64,
	displayUnitID int,
	typeID int,
	marginOfError *float64,
	connectionTypeID int,
	ingestionType *string,
	pollingRate *int,
	jsonData []map[string]interface{},
) (string, error) {

	ctx := context.Background()

	// Marshal JSON input
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSONData: %w", err)
	}

	var uid string

	query := `
        INSERT INTO feeds 
        (
            name, 
            code, 
            description, 
            low_value, 
            high_value,
            display_unit,
            margin_of_error,
            connection_type,
            ingestion_type,
            polling_rate,
            json,
			type_id
        )
        VALUES
        ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
        RETURNING uid
    `

	err = pool.QueryRow(ctx, query,
		name,
		code,
		description,
		lowValue,
		highValue,
		displayUnitID,
		marginOfError,
		connectionTypeID,
		ingestionType,
		pollingRate,
		jsonBytes,
		typeID,
	).Scan(&uid)

	if err != nil {
		return "", fmt.Errorf("insert feed failed: %w", err)
	}

	return uid, nil
}

func CreateMachineFeedMapping(
	pool *pgxpool.Pool,
	machineID,
	feedID int) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx,
		`INSERT INTO feed_machine_mapping (machine_id, feed_id)
         VALUES ($1, $2)`,
		machineID, feedID,
	)

	return err
}

func UpdateFeed(pool *pgxpool.Pool, uid string, req models.FeedUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1
	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, *req.Name)
		idx++
	}
	if req.Code != nil {
		setClauses = append(setClauses, fmt.Sprintf("code = $%d", idx))
		args = append(args, *req.Code)
		idx++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", idx))
		args = append(args, *req.Description)
		idx++
	}
	if req.LowValue != nil {
		setClauses = append(setClauses, fmt.Sprintf("low_value = $%d", idx))
		args = append(args, *req.LowValue)
		idx++
	}
	if req.HighValue != nil {
		setClauses = append(setClauses, fmt.Sprintf("high_value = $%d", idx))
		args = append(args, *req.HighValue)
		idx++
	}
	if req.DisplayUnitID != nil {
		setClauses = append(setClauses, fmt.Sprintf("display_unit = $%d", idx))
		args = append(args, *req.DisplayUnitID)
		idx++
	}
	if req.FeedTypeID != nil {
		setClauses = append(setClauses, fmt.Sprintf("type_id = $%d", idx))
		args = append(args, *req.FeedTypeID)
		idx++
	}
	if req.MarginOfError != nil {
		setClauses = append(setClauses, fmt.Sprintf("margin_of_error = $%d", idx))
		args = append(args, *req.MarginOfError)
		idx++
	}
	if req.ConnectionTypeID != nil {
		setClauses = append(setClauses, fmt.Sprintf("connection_type = $%d", idx))
		args = append(args, *req.ConnectionTypeID)
		idx++
	}

	if req.IngestionType != nil {
		setClauses = append(setClauses, fmt.Sprintf("ingestion_type = $%d", idx))
		args = append(args, *req.IngestionType)
		idx++
	}
	if req.PollingRate != nil {
		setClauses = append(setClauses, fmt.Sprintf("polling_rate = $%d", idx))
		args = append(args, *req.PollingRate)
		idx++
	}
	if req.JSONData != nil {
		jsonBytes, err := json.Marshal(*req.JSONData)
		if err != nil {
			return fmt.Errorf("failed to marshal JSONData: %w", err)
		}
		setClauses = append(setClauses, fmt.Sprintf("json = $%d", idx))
		args = append(args, jsonBytes)
		idx++
	}
	if req.IsRunning != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_running = $%d", idx))
		args = append(args, *req.IsRunning)
		idx++
	}
	if len(setClauses) == 0 {
		return nil
	}
	args = append(args, uid)
	query := fmt.Sprintf("UPDATE feeds SET %s WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteFeed(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM feeds WHERE uid = $1`, uid)
	return err
}

func CreateFeedHistoryEntry(
	pool *pgxpool.Pool,
	feedID int,
	status *bool,
	remark *string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx,
		`INSERT INTO feed_history (feed_id, status, remark)
		 VALUES ($1, $2, $3)`,
		feedID, status, remark,
	)
	return err
}
