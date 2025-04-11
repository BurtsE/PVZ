package users

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	domain "pvz/internal/domain/users"
	"pvz/internal/service/users"
)

var _ users.UserRepository = (*PostgresRepo)(nil)

type userDB struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	Role         string    `db:"role"`
	PasswordHash []byte    `db:"password_hash"`
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

func (r *PostgresRepo) GetUser(ctx context.Context, email string) (domain.User, error) {
	var user userDB
	query := `
		SELECT id, role, password_hash
		FROM users
		WHERE email = $1
	`
	row := r.conn.QueryRow(ctx, query, email)
	if err := row.Scan(&user.ID, &user.Role, &user.PasswordHash); err != nil {
		return domain.User{}, fmt.Errorf("failed to select user: %w", err)
	}
	entity, err := domain.NewUser(user.ID, email, user.Role, user.PasswordHash)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to init user entity: %w", err)
	}
	return entity, nil
}

func (r *PostgresRepo) GetUserById(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var user userDB
	query := `SELECT email, role, password_hash FROM users WHERE id = $1`
	row := r.conn.QueryRow(ctx, query, id)
	if err := row.Scan(&user.Email, &user.Role, &user.PasswordHash); err != nil {
		return domain.User{}, fmt.Errorf("failed to select user: %w", err)
	}
	entity, err := domain.NewUser(id, user.Email, user.Role, user.PasswordHash)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to init user entity: %w", err)
	}
	return entity, nil
}

func (r *PostgresRepo) CreateUser(ctx context.Context, user domain.User) error {
	query := `
		INSERT INTO users (id, email, role, password_hash)
		VALUES ($1, $2, $3, $4)
	`
	if _, err := r.conn.Exec(ctx, query, user.ID(), user.Email(), user.Role(), user.PasswordHash()); err != nil {
		return fmt.Errorf("failed to exec insert user query: %w", err)
	}
	return nil
}
