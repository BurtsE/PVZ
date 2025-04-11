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

type receptionDB struct {
	id       uuid.UUID
	status   string
	pvzId    uuid.UUID
	initTime time.Time
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
	lastReception, err := r.lastReception(ctx, reception.PvzID())
	if err != nil {
		return err
	}
	if lastReception.status == domain.IN_PROGRESS {
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
	lastReception, err := r.lastReception(ctx, pvzID)
	if err != nil {
		return domain.Reception{}, err
	}
	if lastReception.status == domain.CLOSED {
		return domain.Reception{}, fmt.Errorf("no open receptions exists")
	}
	query := `
		UPDATE receptions SET status = $1 WHERE (point_id = $2 AND status = 'in_progress')
	`

	_, err = r.conn.Exec(ctx, query, domain.CLOSED, pvzID)
	if err != nil {
		return domain.Reception{}, fmt.Errorf("failed to close reception query: %w", err)
	}
	reception, err := domain.NewReception(lastReception.id, pvzID, domain.CLOSED, lastReception.initTime, nil)
	if err != nil {
		return domain.Reception{}, err
	}
	return *reception, nil
}

func (r *PostgresRepo) AddProduct(ctx context.Context, product domain.Product, pvzID uuid.UUID) (uuid.UUID, error) {
	reception, err := r.lastReception(ctx, pvzID)
	if err != nil {
		return uuid.Nil, err
	}
	if reception.status == domain.CLOSED {
		return uuid.Nil, fmt.Errorf("no open receptions exists")
	}
	query := `
		INSERT INTO products (id, arrival_date, type, reception_id)
		VALUES ($1, $2, $3, $4)
	`
	_, err = r.conn.Exec(ctx, query, product.ID(), product.ArrivalTime(), product.ProductType(), reception.id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to close reception query: %w", err)
	}
	return reception.id, nil
}

func (r *PostgresRepo) lastReception(ctx context.Context, pvzID uuid.UUID) (receptionDB, error) {
	query := `
		SELECT id, status, registration_date  FROM receptions 
		WHERE point_id = $1
        ORDER BY registration_date DESC
        LIMIT 1
        `
	var reception receptionDB
	err := r.conn.QueryRow(ctx, query, pvzID).Scan(&reception.id, &reception.status, &reception.initTime)
	if err == pgx.ErrNoRows {
		reception.status = domain.CLOSED
	} else if err != nil {
		return receptionDB{}, fmt.Errorf("failed to query reception status: %w", err)
	}
	return reception, nil
}
