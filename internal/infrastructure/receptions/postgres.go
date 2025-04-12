package receptions

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	domain "pvz/internal/domain/receptions"
	reseptionService "pvz/internal/service/receptions"
	"time"
)

var _ reseptionService.ReceptionRepository = (*PostgresRepo)(nil)

var ErrNoOpenReception = fmt.Errorf("no open receptions exists")

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

type pointDB struct {
	id               uuid.UUID
	registrationDate time.Time
	city             string
}
type receptionDB struct {
	id               uuid.UUID
	status           string
	registrationDate time.Time
}

type productDb struct {
	id          uuid.UUID
	arrivalTime time.Time
	productType string
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

func (r *PostgresRepo) GetPointList(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*domain.PickUpPoint, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	query := `
		SELECT points.id, points.registration_date, points.city 
		FROM points
		LIMIT $1 OFFSET $2
	`
	pvzRows, _ := tx.Query(ctx, query, limit, offset)
	defer pvzRows.Close()
	var (
		tempPoints []pointDB
		tempPoint  pointDB
	)
	_, err = pgx.ForEachRow(pvzRows, []any{&tempPoint.id, &tempPoint.registrationDate, &tempPoint.city}, func() error {
		tempPoints = append(tempPoints, tempPoint)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to exec select points: %w", err)
	}
	points := make([]*domain.PickUpPoint, 0, len(tempPoints))
	for _, p := range tempPoints {
		query = `
			SELECT receptions.id, registration_date, status, products.id, arrival_date, type
			FROM receptions INNER JOIN products ON (receptions.id = products.reception_id)
			WHERE point_id = $1
				AND registration_date BETWEEN $2 AND $3
			`
		productRows, _ := tx.Query(ctx, query, p.id, startDate, endDate)
		defer productRows.Close()

		var (
			receptionsMap = map[receptionDB][]productDb{}
			reception     receptionDB
			product       productDb
		)
		_, err = pgx.ForEachRow(productRows, []any{&reception.id, &reception.registrationDate, &reception.status, &product.id,
			&product.arrivalTime, &product.productType},
			func() error {
				log.Println(reception.registrationDate, reception.id)
				receptionsMap[reception] = append(receptionsMap[reception], product)
				return nil
			})
		if err != nil {
			return nil, fmt.Errorf("failed to exec select products: %w", err)
		}
		receptions, err := convertReceptionsMapToDomain(receptionsMap)
		if err != nil {
			return nil, err
		}
		point, err := domain.NewPickUpPoint(tempPoint.id, tempPoint.registrationDate, p.city, receptions)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
		log.Println(points)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return points, nil
}

func (r *PostgresRepo) CreateReception(ctx context.Context, reception domain.Reception, pvzID uuid.UUID) error {
	lastReception, err := r.lastReception(ctx, pvzID)
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
		reception.InitTime(), reception.Status(), pvzID); err != nil {
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
		return domain.Reception{}, ErrNoOpenReception
	}
	query := `
		UPDATE receptions SET status = $1 WHERE (point_id = $2 AND status = 'in_progress')
	`

	_, err = r.conn.Exec(ctx, query, domain.CLOSED, pvzID)
	if err != nil {
		return domain.Reception{}, fmt.Errorf("failed to close reception query: %w", err)
	}
	reception, err := domain.NewReception(lastReception.id, domain.CLOSED, lastReception.registrationDate, nil)
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
		return uuid.Nil, ErrNoOpenReception
	}
	query := `
		INSERT INTO products (id, arrival_date, type, reception_id)
		VALUES ($1, $2, $3, $4)
	`
	_, err = r.conn.Exec(ctx, query, product.ID(), product.ArrivalTime(), product.ProductType(), reception.id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to add product query: %w", err)
	}
	return reception.id, nil
}

func (r *PostgresRepo) DeleteLastAddedProduct(ctx context.Context, pvzID uuid.UUID) error {
	reception, err := r.lastReception(ctx, pvzID)
	if err != nil {
		return err
	}
	if reception.status == domain.CLOSED {
		return ErrNoOpenReception
	}
	query := `
		DELETE FROM products
		WHERE ctid = (
			SELECT ctid
			FROM products
			WHERE reception_id = $1
			ORDER BY arrival_date DESC, id DESC
			LIMIT 1
		)
	`
	_, err = r.conn.Exec(ctx, query, reception.id)
	if err != nil {
		return fmt.Errorf("failed to delete product query: %w", err)
	}

	return nil

}

func (r *PostgresRepo) lastReception(ctx context.Context, pvzID uuid.UUID) (receptionDB, error) {
	query := `
		SELECT id, status, registration_date  FROM receptions 
		WHERE point_id = $1
        ORDER BY registration_date DESC
        LIMIT 1
        `
	var reception receptionDB
	err := r.conn.QueryRow(ctx, query, pvzID).Scan(&reception.id, &reception.status, &reception.registrationDate)
	if err == pgx.ErrNoRows {
		reception.status = domain.CLOSED
	} else if err != nil {
		return receptionDB{}, fmt.Errorf("failed to query reception status: %w", err)
	}
	return reception, nil
}
