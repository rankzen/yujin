package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/pkg/logging"
)

// This is a wrapper for exclusivly pgx "QUERY" logic
type pgxDB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

func newPgxPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres connection string: %w", err)
	}

	pgxCfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		return conn.Ping(ctx) == nil
	}

	RegisterTracers(pgxCfg.ConnConfig)

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return pool, nil
}

func RegisterTracers(connCfg *pgx.ConnConfig) {
	connCfg.Tracer = &zapTracer{}
}

type zapTracer struct{}

func (z *zapTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	logger := logging.FromContext(ctx).Sugar()

	q := strings.Join(strings.Fields(data.SQL), " ")
	logger.Debugf("query start: %.20s", q)

	return ctx
}

func (z *zapTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	logger := logging.FromContext(ctx).Sugar()

	logger.Debugf("The flip is a command tag: %v", data.CommandTag)
}

func (z *zapTracer) TraceBatchStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchStartData) context.Context {
	logger := logging.FromContext(ctx).Sugar()

	logger.Debugf("batch start queued: %v", data.Batch.Len())

	return ctx
}

func (z *zapTracer) TraceBatchQuery(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchQueryData) {
	logger := logging.FromContext(ctx).Sugar()

	q := strings.Join(strings.Fields(data.SQL), " ")
	logger.Debugf("batch query done: %.20s", q)
}

func (z *zapTracer) TraceBatchEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchEndData) {
	logger := logging.FromContext(ctx).Sugar()

	logger.Debugf("batch done: %v", data.Err)
}