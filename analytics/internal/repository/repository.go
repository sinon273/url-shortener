package analytics_repository

import (
	"github.com/ClickHouse/clickhouse-go/v2"
)

type Repository struct {
	conn clickhouse.Conn
}

func NewRepository(conn clickhouse.Conn) *Repository {
	return &Repository{conn: conn}
}
