package telegram

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"1litw/application"
	"1litw/domain"
	"1litw/sqlc"
	"1litw/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	bot         *tgbotapi.BotAPI
	urlUseCase  *application.URLUseCase
	userUseCase *application.UserUseCase
	queries     *sqlc.Queries // For token handling
	appBaseURL  string
}

func NewBotHandler(bot *tgbotapi.BotAPI, urlUC *application.URLUseCase, userUC *application.UserUseCase, db *sql.DB, baseURL string) *BotHandler {
	return &BotHandler{
		bot:         bot,
		urlUseCase:  urlUC,
		userUseCase: userUC,
		queries:     sqlc.New(db),
		appBaseURL:  baseURL,
	}
}

func (h *BotHandler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		h.handleCommand(update.Message)
	}
}

func (h *BotHandler) handleCommand(message *tgbotapi.Message) {
	var err error
	switch message.Command() {
	case "start":
		err = h.handleStart(message)
	case "auth":
		err = h.handleAuth(message)
	case "whoami":
		err = h.handleWhoami(message)
	case "new":
		err = h.handleNew(message)
	case "delete":
		err = h.handleDelete(message)
	case "list":
		err = h.handleList(message)
	default:
		h.sendMessage(message.Chat.ID, "Unknown command.")
	}

	if err != nil {
		log.Printf("Error handling command '%s': %v", message.Command(), err)
		h.sendMessage(message.Chat.ID, fmt.Sprintf("An error occurred: %s", err.Error()))
	}
}

func (h *BotHandler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	h.bot.Send(msg)
}

// Command Handlers

func (h *BotHandler) handleStart(message *tgbotapi.Message) error {
	text := "Welcome to the 1litw URL Shortener Bot!\n\n" +
		"Available commands:\n" +
		"/auth - Link your web account to Telegram\n" +
		"/new <url> [custom_path] - Create a new short link\n" +
		"/list - List your links\n" +
		"/delete <short_path> - Delete a link\n" +
		"/whoami - Check your linked account"
	h.sendMessage(message.Chat.ID, text)
	return nil
}

func (h *BotHandler) handleAuth(message *tgbotapi.Message) error {
	token := utils.GenerateRandomString(32)
	expiresAt := time.Now().Add(5 * time.Minute)

	_, err := h.queries.CreateTelegramAuthToken(context.Background(), sqlc.CreateTelegramAuthTokenParams{
		Token:          token,
		TelegramChatID: message.Chat.ID,
		ExpiresAt:      expiresAt,
	})
	if err != nil {
		return fmt.Errorf("could not create auth token: %w", err)
	}

	authURL := fmt.Sprintf("%s/auth/telegram?token=%s", h.appBaseURL, token)
	text := "To link your account, please open the following link in your browser:\n\n" + authURL
	h.sendMessage(message.Chat.ID, text)
	return nil
}

func (h *BotHandler) handleWhoami(message *tgbotapi.Message) error {
	user, err := h.getUserFromTelegram(context.Background(), message.Chat.ID)
	if err != nil {
		return err
	}

	text := fmt.Sprintf("You are linked as: %s", user.Username)
	h.sendMessage(message.Chat.ID, text)
	return nil
}

func (h *BotHandler) handleNew(message *tgbotapi.Message) error {
	user, err := h.getUserFromTelegram(context.Background(), message.Chat.ID)
	if err != nil {
		return err
	}

	args := strings.Fields(message.CommandArguments())
	if len(args) < 1 {
		h.sendMessage(message.Chat.ID, "Usage: /new <original_url> [custom_path]")
		return nil
	}

	originalURL := args[0]
	customPath := ""
	if len(args) > 1 {
		customPath = args[1]
	}

	shortURL, err := h.urlUseCase.CreateShortURL(context.Background(), user, originalURL, customPath)
	if err != nil {
		return err
	}

	fullShortURL := fmt.Sprintf("%s/%s", h.appBaseURL, shortURL.ShortPath)
	text := fmt.Sprintf("Success! Your short link is: %s", fullShortURL)
	h.sendMessage(message.Chat.ID, text)
	return nil
}

func (h *BotHandler) handleDelete(message *tgbotapi.Message) error {
	user, err := h.getUserFromTelegram(context.Background(), message.Chat.ID)
	if err != nil {
		return err
	}

	shortPath := message.CommandArguments()
	if shortPath == "" {
		h.sendMessage(message.Chat.ID, "Usage: /delete <short_path>")
		return nil
	}

	// Need to get the ID from the path first
	shortURL, err := h.urlUseCase.GetByPath(context.Background(), shortPath)
	if err != nil {
		return err
	}
	if shortURL == nil {
		h.sendMessage(message.Chat.ID, "Short URL not found.")
		return nil
	}

	err = h.urlUseCase.DeleteShortURLByID(context.Background(), user, shortURL.ID)
	if err != nil {
		return err
	}

	h.sendMessage(message.Chat.ID, fmt.Sprintf("Successfully deleted link for: %s", shortPath))
	return nil
}

func (h *BotHandler) handleList(message *tgbotapi.Message) error {
	user, err := h.getUserFromTelegram(context.Background(), message.Chat.ID)
	if err != nil {
		return err
	}

	urls, err := h.urlUseCase.ListByUser(context.Background(), user)
	if err != nil {
		return err
	}

	if len(urls) == 0 {
		h.sendMessage(message.Chat.ID, "You haven't created any links yet.")
		return nil
	}

	var builder strings.Builder
	builder.WriteString("Your links:\n")
	for _, u := range urls {
		fullShortURL := fmt.Sprintf("%s/%s", h.appBaseURL, u.ShortPath)
		builder.WriteString(fmt.Sprintf("\n- %s -> %s (%d clicks)", fullShortURL, u.OriginalURL, u.TotalClicks))
	}

	h.sendMessage(message.Chat.ID, builder.String())
	return nil
}

// getUserFromTelegram authenticates a user based on their Telegram chat ID.
func (h *BotHandler) getUserFromTelegram(ctx context.Context, chatID int64) (*domain.User, error) {
	user, err := h.userUseCase.GetUserByTelegramID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if user == nil {
		return nil, errors.New("account not linked. Please use /auth first")
	}
	return user, nil
}
