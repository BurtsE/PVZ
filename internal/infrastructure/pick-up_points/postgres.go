package pick_up_points

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type PostgresRepo struct {
	conn *pgx.Conn
}

func NewPostgresRepo(conn *pgx.Conn) *PostgresRepo {
	return &PostgresRepo{
		conn: conn,
	}
}
func (r *PostgresRepo) Close() error {
	return r.conn.Close(context.Background())
}
