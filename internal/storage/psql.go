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
}

func CreatePostgresConnection(cfg ConnectionInfo) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), cfg.makeURL())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c ConnectionInfo) makeURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.Username, c.Password, c.Host, c.Port, c.DBName)
}
