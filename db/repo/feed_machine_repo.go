package repo

import (
	"app/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ListFeedMachineMappings(pool *pgxpool.Pool, feedUID string) ([]models.FeedMachineMappingResponse, error) {
	ctx := context.Background()

	query := `
        SELECT 
            fmm.uid AS mapping_uid,
            f.uid  AS feed_uid,
            m.uid  AS machine_uid
        FROM feed_machine_mapping fmm
        JOIN feeds f ON fmm.feed_id = f.id
        JOIN machines m ON fmm.machine_id = m.id
        WHERE f.uid = $1
        ORDER BY fmm.id ASC
    `

	rows, err := pool.Query(ctx, query, feedUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.FeedMachineMappingResponse

	for rows.Next() {
		var r models.FeedMachineMappingResponse
		if err := rows.Scan(&r.UID, &r.FeedUID, &r.MachineUID); err != nil {
			return nil, err
		}
		out = append(out, r)
	}

	return out, nil
}

func ListMachineFeedMappings(
	pool *pgxpool.Pool,
	machineUID string,
) ([]models.FeedMachineMappingResponse, error) {

	ctx := context.Background()

	query := `
        SELECT 
            fmm.uid AS mapping_uid,
            f.uid  AS feed_uid,
            m.uid  AS machine_uid
        FROM feed_machine_mapping fmm
        JOIN feeds f ON fmm.feed_id = f.id
        JOIN machines m ON fmm.machine_id = m.id
        WHERE m.uid = $1
        ORDER BY fmm.id ASC
    `

	rows, err := pool.Query(ctx, query, machineUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.FeedMachineMappingResponse

	for rows.Next() {
		var r models.FeedMachineMappingResponse
		if err := rows.Scan(&r.UID, &r.FeedUID, &r.MachineUID); err != nil {
			return nil, err
		}
		out = append(out, r)
	}

	return out, nil
}

func CreateFeedMachineMapping(
	pool *pgxpool.Pool,
	feedID int,
	machineID int) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx,
		`INSERT INTO feed_machine_mapping (machine_id, feed_id)
         VALUES ($1, $2)`,
		machineID, feedID,
	)

	return err
}

func FeedMachineMappingExists(pool *pgxpool.Pool, feedID, machineID int) (bool, error) {
	ctx := context.Background()
	var exists bool

	query := `
        SELECT EXISTS (
            SELECT 1 
            FROM feed_machine_mapping 
            WHERE feed_id = $1 AND machine_id = $2
        )
    `
	err := pool.QueryRow(ctx, query, feedID, machineID).Scan(&exists)
	return exists, err
}

func DeleteFeedMachineMapping(pool *pgxpool.Pool, feedID, machineID int) error {
	ctx := context.Background()

	query := `
        DELETE FROM feed_machine_mapping
        WHERE feed_id = $1 AND machine_id = $2
    `

	cmd, err := pool.Exec(ctx, query, feedID, machineID)
	if err != nil {
		return fmt.Errorf("failed to delete mapping: %w", err)
	}

	// Optional: ensure deletion actually happened
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("mapping not found")
	}

	return nil
}
