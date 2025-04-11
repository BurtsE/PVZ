package receptions

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	domain "pvz/internal/domain/receptions"
	reseptionService "pvz/internal/service/receptions"
	"time"
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
		return fmt.Errorf("failed to exec insert point query: %w", err)
	}
	return nil
}

func (r *PostgresRepo) CreateReception(ctx context.Context, reception domain.Reception) error {
	status, err := r.lastReceptionStatus(ctx, reception.PvzID())
	if err != nil {
		return err
	}
	if status == domain.IN_PROGRESS {
		return fmt.Errorf("opened reception already exists")
	}
	query := `
		INSERT INTO receptions (id, registration_date, status, point_id)
		VALUES ($1, $2, $3, $4)
	`
	if _, err = r.conn.Exec(ctx, query, reception.ID(),
		reception.InitTime(), reception.Status(), reception.PvzID()); err != nil {
		return fmt.Errorf("failed to exec insert reception query: %w", err)
	}
	return nil
}

func (r *PostgresRepo) CloseReception(ctx context.Context, pvzID uuid.UUID) (domain.Reception, error) {
	status, err := r.lastReceptionStatus(ctx, pvzID)
	if err != nil {
		return domain.Reception{}, err
	}
	if status == domain.CLOSED {
		return domain.Reception{}, fmt.Errorf("no open receptions exists")
	}
	query := `
		UPDATE receptions SET status = $1 WHERE (point_id = $2 AND status = 'in_progress')
		RETURNING id, registration_date
	`
	var (
		id               uuid.UUID
		registrationDate time.Time
	)
	err = r.conn.QueryRow(ctx, query, domain.CLOSED, pvzID).Scan(&id, &registrationDate)
	if err != nil {
		return domain.Reception{}, fmt.Errorf("failed to close reception query: %w", err)
	}
	reception, err := domain.NewReception(id, pvzID, status, registrationDate, nil)
	if err != nil {
		return domain.Reception{}, err
	}
	return *reception, nil

}

func (r *PostgresRepo) lastReceptionStatus(ctx context.Context, pvzID uuid.UUID) (string, error) {
	query := `
		SELECT status FROM receptions 
		WHERE point_id = $1
        ORDER BY registration_date DESC
        LIMIT 1
        `
	var status string
	err := r.conn.QueryRow(ctx, query, pvzID).Scan(&status)
	if err == pgx.ErrNoRows {
		status = domain.CLOSED
	} else if err != nil {
		return "", fmt.Errorf("failed to query reception status: %w", err)
	}
	return status, nil

}
