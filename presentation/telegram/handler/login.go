package handler

import (
	"context"
	"log"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (h *Handler) Login(b *gotgbot.Bot, ctx *ext.Context) error {
	send := sender(b, ctx)
	if ctx.EffectiveChat.Type != gotgbot.ChatTypePrivate {
		return send("Please use this command in a private chat with me.")
	}

	token, err := h.userUseCase.PrepareLinkTelegram(context.Background(), ctx.EffectiveUser.Id)
	if err != nil {
		return send("Failed to prepare login link. Please try again later or contact support.")
	}

	link := h.cfg.Base + "/auth/telegram?token=" + token
	log.Println(link)

	return send(`Click the [%s](%s) to login`, escapeMarkdownV2(link), link)
}

func escapeMarkdownV2(text string) string {
	// Escape special characters for MarkdownV2
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, `\`+char)
	}
	return text
}
