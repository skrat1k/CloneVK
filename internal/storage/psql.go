package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ConnectionInfo struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

func CreatePostgresConnection(cfg string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func MakeURL(info ConnectionInfo) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", info.Username, info.Password, info.Host, info.Port, info.DBName, info.SSLMode)
}
