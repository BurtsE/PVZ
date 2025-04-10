package users

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	domain "pvz/internal/domain/users"
)

//var _ receptions.UserRepository = (*PostgresRepo)(nil)

type userDB struct {
	ID    uuid.UUID `db:"id"`
	Name  string    `db:"name"`
	Email string    `db:"email"`
	Role  string    `db:"role"`
}

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

func (r *PostgresRepo) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var user userDB
	query := `SELECT name, email, role FROM users WHERE id = $1`
	row := r.conn.QueryRow(ctx, query, id)

	if err := row.Scan(&user.Name, &user.Email, &user.Role); err != nil {
		return domain.User{}, fmt.Errorf("failed to select user: %w", err)
	}
	entity, err := domain.NewUser(id, user.Name, user.Email, user.Role)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to init user entity: %w", err)
	}
	return entity, nil
}

func (r *PostgresRepo) CreateUser(ctx context.Context, user domain.User) error {
	query := `
		INSERT INTO users (id, name, email, role)
		VALUES ($1, $2, $3, $4)
	`
	if _, err := r.conn.Exec(ctx, query, user.ID(), user.Name(), user.Email(), user.Role()); err != nil {
		return fmt.Errorf("failed to exec insert user query: %w", err)
	}
	return nil
}
