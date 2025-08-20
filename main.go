package main

import (
	"context"
	"database/sql"
	"embed"
	"log"

	"1litw/application"
	"1litw/config"
	"1litw/domain"
	"1litw/infrastructure/external"
	"1litw/infrastructure/repository"
	"1litw/presentation/gin"
	"1litw/presentation/telegram"
	"1litw/presentation/telegram/handler"

	_ "modernc.org/sqlite"
)

//go:embed all:web/dist
var webDist embed.FS

//go:embed sql/schema.sql
var schemaSQL string

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite", cfg.DBPath+"?_foreign_keys=on")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Ensure database schema and initial data are set up
	if err := ensureInitialData(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repositories
	clickRepo := repository.NewClickRepository(db)
	userRepo := repository.NewUserRepository(db)
	urlRepo := repository.NewShortURLRepository(db)
	analyticsRepo := repository.NewClickRepository(db)
	tgAuthTokenRepo := repository.NewTGAuthTokenRepository(db)

	// Initialize external services
	uaParser := external.NewUAParserService()
	geoIPProcessor := external.NewGeoIPProcessor(clickRepo)
	geoIPProcessor.Start()

	// Initialize use cases
	userUC := application.NewUserUseCase(cfg.JWTSecret, userRepo, tgAuthTokenRepo)
	urlUC := application.NewURLUseCase(urlRepo, userRepo, analyticsRepo, uaParser)
	analyticsUC := application.NewAnalyticsUseCase(analyticsRepo, urlRepo)

	// Setup router
	router := gin.SetupRouter(db, webDist, cfg.JWTSecret, userUC, urlUC, analyticsUC)

	// Start Telegram Bot if token is provided
	if cfg.BotToken != "" {
		h := handler.New(cfg, userUC)

		bot, err := telegram.NewBot(cfg.BotToken, h)
		if err != nil {
			panic(err)
		}

		go bot.Run(context.TODO())
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// ensureInitialData ensures that the database schema is created and the
// necessary initial data (like the 'anonymous' user) exists.
func ensureInitialData(db *sql.DB) error {
	log.Println("Initializing database...")

	// Execute the schema script to create tables if they don't exist.
	// The schema should use `CREATE TABLE IF NOT EXISTS` to be idempotent.
	if _, err := db.Exec(schemaSQL); err != nil {
		return err
	}

	// Check for and create the 'anonymous' user if it doesn't exist.
	userRepo := repository.NewUserRepository(db)
	ctx := context.Background()

	anonUser, err := userRepo.GetByUsername(ctx, "anonymous")
	if err != nil {
		// If the error is anything other than "not found", it's a real problem.
		if err != domain.ErrNotFound {
			return err
		}
		// The user was not found, so we create it.
		log.Println("Anonymous user not found, creating it...")
		_, err := userRepo.Create(ctx, &domain.User{
			Username:     "anonymous",
			PasswordHash: "*", // An invalid hash to prevent login
			Permissions:  domain.RoleGuest,
		})
		if err != nil {
			return err
		}
		log.Println("Anonymous user created successfully.")
	} else if anonUser != nil {
		log.Println("Anonymous user already exists.")
	}

	log.Println("Database initialization complete.")
	return nil
}
