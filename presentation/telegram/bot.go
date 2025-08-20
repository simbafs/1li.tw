package telegram

import (
	"context"
	"fmt"
	"log"

	"1litw/presentation/telegram/handler"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type Bot struct {
	b       *gotgbot.Bot
	handler *handler.Handler
}

func NewBot(token string, handler *handler.Handler) (*Bot, error) {
	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}
	return &Bot{
		b:       b,
		handler: handler,
	}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("bot dispatcher: ", err)
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	// register handlers here
	dispatcher.AddHandler(handlers.NewCommand("login", b.handler.Login))

	updater := ext.NewUpdater(dispatcher, nil)

	if err := updater.StartPolling(b.b, &ext.PollingOpts{}); err != nil {
		return err
	}

	log.Println("Telegram bot started successfully")

	// Handle graceful shutdown
	go func() {
		<-ctx.Done()
		if err := updater.Stop(); err != nil {
			log.Println("Error stopping updater:", err)
		} else {
			log.Println("Telegram bot stopped successfully")
		}
	}()

	updater.Idle()

	return nil
}
