package postgres

import (
	"context"
	"tp-db-project/internal/models"
	"tp-db-project/pkg/queries"
	"tp-db-project/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatusRepository struct {
	conn *pgxpool.Pool
}

func InitStatusRepository(db *pgxpool.Pool) repository.StatusRepo {
	return &StatusRepository{conn: db}
}

func (s *StatusRepository) GetStatus(cont context.Context) (*models.Status, error) {
	var result models.Status
	_ = s.conn.QueryRow(cont, queries.GetCountRecordsCommand).Scan(&result.User, &result.Forum, &result.Thread, &result.Post)
	return &result, nil
}

func (s *StatusRepository) Clear(cont context.Context) error {
	_, _ = s.conn.Exec(cont, queries.DeleteTablesCommand)
	return nil
}
