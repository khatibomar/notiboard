package main

import (
	"context"
	"errors"
	"regexp"
	"sync"

	"github.com/jackc/pgx/v5"
)

const (
	dbConnStr = "postgresql://user:pass@localhost:5432/testdb?sslmode=disable"
)

func sensoredConnString() string {
	re := regexp.MustCompile(`postgresql://([^:]+):([^@]+)@(.+)`)
	return re.ReplaceAllString(dbConnStr, "postgresql://$1:****@$3")
}

type DatabaseManager struct {
	mu   sync.Mutex
	ctx  context.Context
	conn *pgx.Conn
}

func NewDatabaseManager(ctx context.Context) *DatabaseManager {
	return &DatabaseManager{
		ctx: ctx,
	}
}

func (dm *DatabaseManager) Connect() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.conn != nil {
		dm.conn.Close(dm.ctx)
	}

	conn, err := pgx.Connect(dm.ctx, dbConnStr)
	if err != nil {
		return err
	}

	dm.conn = conn
	return nil
}

func (dm *DatabaseManager) Ping() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.conn == nil {
		return errors.New("no database connection")
	}

	return dm.conn.Ping(dm.ctx)
}

func (dm *DatabaseManager) Close() {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.conn != nil {
		dm.conn.Close(dm.ctx)
		dm.conn = nil
	}
}
