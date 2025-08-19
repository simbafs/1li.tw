package handler

import (
	"fmt"
	"log"

	"1litw/application"
	"1litw/config"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Handler struct {
	cfg         *config.Config
	userUseCase *application.UserUseCase
}

func New(cfg *config.Config, userUC *application.UserUseCase) *Handler {
	return &Handler{
		cfg:         cfg,
		userUseCase: userUC,
	}
}

func sendMessage(b *gotgbot.Bot, ctx *ext.Context, text string, opt *gotgbot.SendMessageOpts) (*gotgbot.Message, error) {
	log.Println("send message:", text)
	if ctx.EffectiveChat == nil {
		return nil, fmt.Errorf("no effective chat")
	}
	if opt == nil {
		opt = &gotgbot.SendMessageOpts{
			ParseMode: "MarkdownV2",
		}
	}
	if ctx.EffectiveChat.IsForum {
		opt.MessageThreadId = ctx.EffectiveMessage.MessageThreadId
	}
	m, err := ctx.EffectiveChat.SendMessage(b, text, opt)
	return m, err
}

func sender(b *gotgbot.Bot, ctx *ext.Context) func(string, ...any) error {
	return func(format string, args ...any) error {
		_, err := sendMessage(b, ctx, fmt.Sprintf(format, args...), nil)
		return err
	}
}
