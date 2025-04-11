package receptions

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	domain "pvz/internal/domain/receptions"
	reseptionService "pvz/internal/service/receptions"
)

var _ reseptionService.ReceptionRepository = (*PostgresRepo)(nil)

type PostgresRepo struct {
	conn *pgxpool.Pool
}

func NewPostgresRepo(conn *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{
		conn: conn,
	}
}
func (r *PostgresRepo) Close() {
	r.conn.Close()
}

func (r *PostgresRepo) CreatePoint(ctx context.Context, point domain.PickUpPoint) error {
	query := `
		INSERT INTO points (id, registration_date, city)
		VALUES ($1, $2, $3)
	`
	if _, err := r.conn.Exec(ctx, query, point.ID(), point.RegistrationDate(), point.City()); err != nil {
		return fmt.Errorf("failed to exec insert user query: %w", err)
	}
	return nil

}
