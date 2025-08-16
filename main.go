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
	"1litw/infrastructure/processor"
	"1litw/infrastructure/repository"
	"1litw/infrastructure/telegram"
	"1litw/presentation/gin"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	// Setup repositories
	clickRepo := repository.NewClickRepository(db)

	// Start GeoIP Processor
	geoIPProcessor := processor.NewGeoIPProcessor(clickRepo)
	geoIPProcessor.Start()

	// Setup router
	router := gin.SetupRouter(db, cfg.JWTSecret, webDist)

	// Start Telegram Bot if token is provided
	if cfg.BotToken != "" {
		go startTelegramBot(cfg, db)
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

func startTelegramBot(cfg *config.Config, db *sql.DB) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Printf("Failed to create Telegram bot: %v", err)
		return
	}
	bot.Debug = cfg.Environment == "development"
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// We need to re-create dependencies here for the bot.
	// This suggests a dependency injection container would be beneficial for a larger app.
	userRepo := repository.NewUserRepository(db)
	urlRepo := repository.NewShortURLRepository(db)
	clickRepo := repository.NewClickRepository(db)
	uaParser := external.NewUAParserService()
	userUseCase := application.NewUserUseCase(userRepo, cfg.JWTSecret)
	urlUseCase := application.NewURLUseCase(urlRepo, userRepo, clickRepo, uaParser)

	// The base URL for links needs to be configured.
	// For now, we'll construct it from the server port.
	baseURL := "http://localhost:" + cfg.ServerPort

	botHandler := telegram.NewBotHandler(bot, urlUseCase, userUseCase, db, baseURL)
	botHandler.Start()
}
