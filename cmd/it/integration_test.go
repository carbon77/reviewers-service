package integration_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reviewers/internal/handler"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate"
	migratepg "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		Env:          map[string]string{"POSTGRES_PASSWORD": "secret", "POSTGRES_USER": "test", "POSTGRES_DB": "testdb"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	defer pgContainer.Terminate(ctx)

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("host=%s port=%s user=test password=secret dbname=testdb sslmode=disable", host, port.Port())
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Run migrations
	sqlDB, _ := db.DB()
	driver, _ := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	mig, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	mig.Up()

	os.Exit(m.Run())
}

func setupRouter(tx *gorm.DB) *gin.Engine {
	logger := slog.Default()
	router := gin.Default()
	handler.InitHandlers(logger, tx, router)
	return router
}

func runInTransaction(t *testing.T, testFunc func(tx *gorm.DB)) {
	tx := db.Begin()
	defer tx.Rollback()
	testFunc(tx)
}
