package app

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lmittmann/tint"

	// Register PostgreSQL driver.
	_ "github.com/lib/pq"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "Subscription_Service/docs"
	"Subscription_Service/internal/application/service"
	"Subscription_Service/internal/config"
	httpHandler "Subscription_Service/internal/infrastructure/controllers/http"
	"Subscription_Service/internal/infrastructure/repository"
	httpServer "Subscription_Service/pkg/http_server"
)

type App struct {
	config   *config.Config
	logger   *slog.Logger
	db       *sqlx.DB
	server   *httpServer.HTTPServer
	services service.Service
}

func New() (*App, error) {
	logger := setupLogger()

	cfg, err := config.LoadConfig("config/config.yaml", ".env")
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		return nil, err
	}

	db, err := initDatabase(cfg, logger)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		return nil, err
	}

	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, logger)
	services := service.NewService(subscriptionService)
	router := initRouter(services, logger)
	serverConfig := &httpServer.Config{
		Host:              cfg.Service.Host,
		Port:              cfg.Service.Port,
		StartMsg:          "Subscription service started!",
		Handler:           router,
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		ShutdownTimeout:   30 * time.Second,
	}
	server := httpServer.NewServer(logger, serverConfig)

	return &App{
		config:   cfg,
		logger:   logger,
		db:       db,
		server:   server,
		services: services,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	if err := a.server.Start(ctx); err != nil {
		a.logger.Error("Failed to start server", "error", err)
		return err
	}

	if err := a.db.Close(); err != nil {
		a.logger.Error("Failed to close database connection", "error", err)
		return err
	}

	a.logger.Info("Application shutdown complete")
	return nil
}

func setupLogger() *slog.Logger {
	return slog.New(tint.NewHandler(os.Stdout, nil))
}

func initDatabase(cfg *config.Config, logger *slog.Logger) (*sqlx.DB, error) {
	logger.Info("Connecting to database", "dsn", cfg.Database.GetDSN())

	db, err := sqlx.Connect("postgres", cfg.Database.GetDSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	logger.Info("Database connected successfully")
	return db, nil
}

func initRouter(service service.Service, logger *slog.Logger) *gin.Engine {
	handler := httpHandler.NewHandler(service)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/health", healthCheck)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	handler.RegisterRoutes(router)

	logger.Info("Router initialized")
	return router
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
