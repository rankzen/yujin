package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/KnutZuidema/golio"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	repository       = "postgres"
	tag              = "alpine"
	databaseName     = "test_db"
	databaseUser     = "test_user"
	databasePassword = "testing123"
	riotApiKey       = ""
)

// HELPER FOR TEST SUITES.
// Should be initialized in TestMain
type TestInstance interface {
	NewDatabase(tb testing.TB) DB
        GetGolioClient() *golio.Client
	Close() error
	MustClose()
}

type testInstance struct {
	skipDB       bool
	skipDBReason string

	pool        *dockertest.Pool
	container   *dockertest.Resource
	golioClient *golio.Client
	conn        *pgx.Conn
	url         *url.URL
	mu          sync.Mutex
}

func NewTestInstance() (TestInstance, error) {
	if !flag.Parsed() {
		flag.Parse()
	}

	if testing.Short() {
		return &testInstance{
			skipDB:       true,
			skipDBReason: "-short flag provided",
		}, nil
	}

	// now we are for sure making the db
	log.Printf("Connecting to Docker")
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Docker: %w", err)
	}

	log.Printf("Running Docker Container: %s:%s", repository, tag)
	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: repository,
		Tag:        tag,
		Env: []string{
			"POSTGRES_DB=" + databaseName,
			"POSTGRES_USER=" + databaseUser,
			"POSTGRES_PASSWORD=" + databasePassword,
		},
	},
		func(c *docker.HostConfig) {
			c.AutoRemove = true
			c.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	err = container.Expire(120)
	if err != nil {
		return nil, fmt.Errorf("failed to expire database: %w", err)
	}

	connUrl := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(databaseUser, databasePassword),
		Host:     container.GetHostPort("5432/tcp"),
		Path:     databaseName,
		RawQuery: "sslmode=disable",
	}

	time.Sleep(5 * time.Second)

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, connUrl.String())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to databse: %w", err)
	}

	err = dbMigrate(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate databse: %w", err)
	}

        log.Printf("Starting Golio Client")
        gc := NewGolioClient(riotApiKey)

	return &testInstance{
		skipDB:      false,
		pool:        pool,
		container:   container,
		conn:        conn,
		url:         connUrl,
		golioClient: gc,
	}, nil
}

func MustTestInstance() TestInstance {
	db, err := NewTestInstance()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return db
}

// Clones a new databse from the template, test will fatal on error
func (t *testInstance) NewDatabase(tb testing.TB) DB {
	tb.Helper()

	if t.skipDB {
		tb.Skipf("skipping db test: %s", t.skipDBReason)
	}

	ctx := context.Background()

	name, err := randomDatabaseName()
	if err != nil {
		tb.Fatalf("no dataname: %s", err)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	q := fmt.Sprintf(`CREATE DATABASE "%s" WITH TEMPLATE "%s";`, name, databaseName)
	_, err = t.conn.Exec(ctx, q)
	if err != nil {
		tb.Fatalf("failed to clone template database: %s", err)
	}

	connUrl := url.URL{
		Scheme:   t.url.Scheme,
		User:     t.url.User,
		Host:     t.url.Host,
		Path:     name,
		RawQuery: t.url.RawQuery,
	}

	pool, err := pgxpool.New(ctx, connUrl.String())
	if err != nil {
		tb.Fatalf("failed to connect: %s", err)
	}

	tb.Cleanup(func() {
		pool.Close()

		t.mu.Lock()
		defer t.mu.Unlock()

		q := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE);`, name)
		_, err := t.conn.Exec(ctx, q)
		if err != nil {
			tb.Errorf("failed to drop database %q: %s", name, err)
		}
	})

	db, err := NewDB(ctx, connUrl.String())
	if err != nil {
                tb.Fatalf("failed to create db :%s", err)
	}

	return db
}

func (t *testInstance) GetGolioClient() *golio.Client {
        return t.golioClient
}

func (t *testInstance) MustClose() {
	err := t.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func (t *testInstance) Close() error {
	err := t.conn.Close(context.Background())
	if err != nil {
		return err
	}

	err = t.pool.Purge(t.container)
	if err != nil {
		return err
	}

	return nil
}

func randomDatabaseName() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}