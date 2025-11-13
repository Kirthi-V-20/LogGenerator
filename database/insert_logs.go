package database

import (
	"LogGenerator/model"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func InsertLogs(ctx context.Context, conn *pgx.Conn, store model.LogStore) error {
	for _, segment := range store.Segments {
		for _, entry := range segment.LogEntries {
			_, err := conn.Exec(ctx, `
				INSERT INTO logs (time, level, host, component, req_id, msg)
				VALUES (
					$1,
					(SELECT id FROM log_Level WHERE level = $2),
					(SELECT id FROM log_Host WHERE host = $3),
					(SELECT id FROM log_Component WHERE component = $4),
					$5,
					$6
				)`,
				entry.Time,
				entry.Level,
				entry.Host,
				entry.Component,
				entry.RequestID,
				entry.Message,
			)
			if err != nil {
				slog.Warn("Failed to insert log entry", "error", err, "entry", entry.Raw)
				continue
			}
		}
	}
	return nil
}
