package postgres

import "time"

type Option func(*Postgres)

func ConnsMaxLifeTime(time time.Duration) Option {
	return func(m *Postgres) {
		m.ConnMaxLifeTime = time
	}
}

func MaxOpenConns(size int) Option {
	return func(m *Postgres) {
		m.MaxOPenConns = size
	}
}

func MaxIdleConns(size int) Option {
	return func(m *Postgres) {
		m.MaxIdleConns = size
	}
}
