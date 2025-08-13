package main

import (
	"database/sql"
	"embed"
	"log"

	"1litw/application"
	"1litw/config"
	"1litw/infrastructure/repository"
	"1litw/infrastructure/telegram"
	"1litw/presentation/gin"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed all:web/dist
var webDist embed.FS

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", cfg.DBPath+"?_foreign_keys=on")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

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
	userUseCase := application.NewUserUseCase(userRepo, cfg.JWTSecret)
	urlUseCase := application.NewURLUseCase(urlRepo, userRepo)

	// The base URL for links needs to be configured.
	// For now, we'll construct it from the server port.
	baseURL := "http://localhost:" + cfg.ServerPort

	botHandler := telegram.NewBotHandler(bot, urlUseCase, userUseCase, db, baseURL)
	botHandler.Start()
}