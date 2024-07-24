package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

const (
	_defaultConnMaxLifeTime = time.Minute * 3
	_defaultMaxOpenConns    = 10
	_defaultMaxIdleConns    = 10
)

type Config struct {
	DbUser string `mapstructure:"db_user"`
	DbPass string `mapstructure:"db_pass"`
	DbHost string `mapstructure:"db_host"`
	DbPort int    `mapstructure:"db_port"`
	DbName string `mapstructure:"db_name"`
}

type Postgres struct {
	ConnMaxLifeTime time.Duration
	MaxOPenConns    int
	MaxIdleConns    int
	Builder         squirrel.StatementBuilderType
	Db              *sql.DB
	Tx              *sql.Tx
}

func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		ConnMaxLifeTime: _defaultConnMaxLifeTime,
		MaxOPenConns:    _defaultMaxOpenConns,
		MaxIdleConns:    _defaultMaxIdleConns,
	}

	for _, opt := range opts {
		opt(pg)
	}

	// Use the dollar sign for PostgreSQL placeholders
	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Update the connection string for PostgreSQL
	db, err := sql.Open("postgres", url)
	if err != nil {
		return &Postgres{}, fmt.Errorf("postgres - new - sql.open: %v", err)
	}

	pg.Db = db
	pg.Db.SetConnMaxLifetime(pg.ConnMaxLifeTime)
	pg.Db.SetMaxOpenConns(pg.MaxOPenConns)
	pg.Db.SetMaxIdleConns(pg.MaxIdleConns)

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Db != nil {
		p.Db.Close()
	}
}
