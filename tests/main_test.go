package tests

import (
	"context"
	"fmt"
	"moneytracker/internal/config"
	"moneytracker/internal/repository"
	"moneytracker/internal/service"
	"moneytracker/internal/transport/http/handler"
	"moneytracker/pkg/database"
	"moneytracker/pkg/logger"

	"os"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type TrackerTestSuite struct {
	suite.Suite
	config       *config.Config
	redis        *redis.Client
	db           *pgxpool.Pool
	spendingUrl  string
	app          *fiber.App
	serverCancel context.CancelFunc
}

func TestTrackerSuite(t *testing.T) {
	fmt.Println("SetupSuite called")
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(TrackerTestSuite))
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

// Вызывается один раз перед всеми тестами, где выполняются «тяжелые» операции (миграции, запуск сервера).
func (s *TrackerTestSuite) SetupSuite() {
	s.initConfig()
	s.spendingUrl = "http://localhost:8000/spending"
	s.initDBConnections()

	// Run migrations
	err := s.runMigrations()
	if err != nil {
		s.FailNow("migration error", err)
	}

	logrus.SetFormatter(new(logrus.JSONFormatter))
}

func (s *TrackerTestSuite) TearDownSuite() {
	s.redis.Close()
	s.dropTables()
	s.db.Close()
}

func (s *TrackerTestSuite) cleanupDb() {
	_, err := s.db.Exec(context.Background(), `
		TRUNCATE TABLE
			spendings
		CASCADE;
`)
	if err != nil {
		s.FailNow("Failed to truncate tables in db", err)
	}
}

func (s *TrackerTestSuite) dropTables() {
	_, err := s.db.Exec(context.Background(), `DROP TABLE IF EXISTS spendings CASCADE;`)
	if err != nil {
		s.FailNow("Failed to drop tables", err)
	}
}

// Вызывается перед каждым тестом — создание изолированного окружения.
func (s *TrackerTestSuite) SetupTest() {
	// s.populateRedis()

	ready := make(chan bool)
	s.startTestServer(ready)

	select {
	case <-ready:
		s.T().Log("Server started successfully")
	case <-time.After(30 * time.Second):
		s.FailNow("Timed out waiting for server to start")
	}
}

func (s *TrackerTestSuite) TearDownTest() {
	// Cleanup server
	if s.app != nil {
		s.app.Shutdown()
	}
	if s.serverCancel != nil {
		s.serverCancel()
	}

	s.cleanupDb()
	s.cleanupRedis()
}

type TestDeviceCloud struct {
}

func NewTestDeviceCloud() *TestDeviceCloud {
	return &TestDeviceCloud{}
}

func (s *TrackerTestSuite) startTestServer(ready chan<- bool) {
	ctx := context.Background()
	//	defer cancel()

	zapLogger, err := logger.NewLogger(zapcore.InfoLevel, "money_tracker.log")
	if err != nil {
		panic(fmt.Sprintf("error initializing logger: %v", err))
	}
	repo := repository.NewRepository(ctx, s.db, s.redis)
	service, err := service.NewService(ctx, repo, zapLogger)
	if err != nil {
		zapLogger.Fatal("Error creating service", zap.Error(err))
	}
	if err := service.Start(ctx); err != nil {
		zapLogger.Fatal("Error starting service", zap.Error(err))
	}
	handler := handler.NewHTTPHandler(ctx, zapLogger, service)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true, // Reduce noise in tests
	})

	handler.SetSpendingRoutes(app)

	s.app = app

	go func() {
		err := app.Listen(":8000")
		if err != nil {
			s.T().Logf("Server stopped: %v", err)
		}
	}()

	// Wait a moment for server to start, then signal ready
	go func() {
		time.Sleep(100 * time.Millisecond) // Give server time to start
		ready <- true
	}()

	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// <-sigChan
	// cancel()
	// app.Shutdown()

	// zapLogger.Info("Service stopped")
}

func (s *TrackerTestSuite) initConfig() {
	appConfig, err := config.NewConfig()
	if err != nil {
		s.FailNowf("error config init: %s", err.Error())
	}

	s.config = appConfig
}

func (s *TrackerTestSuite) initDBConnections() {
	dbConn, err := database.NewPGXPool(s.config.DBSource)
	if err != nil {
		s.FailNow("postgres initialization error", err)
	}
	s.db = dbConn

	redisConn, err := database.NewRedisConn(context.Background(), s.config.RedisAddress, s.config.RedisPassword)
	if err != nil {
		s.FailNow("redis initialization error", err)
	}
	s.redis = redisConn
}

func (s *TrackerTestSuite) runMigrations() error {
	// Read the migrations file
	migrationSQL, err := os.ReadFile("migrations.sql")
	if err != nil {
		return fmt.Errorf("failed to read migrations.sql: %w", err)
	}

	// Execute the migration
	_, err = s.db.Exec(context.Background(), string(migrationSQL))
	if err != nil {
		return fmt.Errorf("failed to execute migrations: %w", err)
	}

	return nil
}

func (s *TrackerTestSuite) cleanupRedis() {
	err := s.redis.FlushAll(context.Background()).Err()
	if err != nil {
		s.FailNow("Failed to clean redis", err)
	}
}
