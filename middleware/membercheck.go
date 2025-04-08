package middleware

import (
	"fmt"
	tele "gopkg.in/telebot.v4"
	"log"
)

type SubscriptionMiddleware struct {
	askMessage string
	chat       *tele.Chat
}

func NewSubscriptionMiddleware(message string, chat *tele.Chat) *SubscriptionMiddleware {
	return &SubscriptionMiddleware{
		askMessage: message,
		chat:       chat,
	}
}

func (s *SubscriptionMiddleware) CheckSubscription(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		member, err := c.Bot().ChatMemberOf(s.chat, c.Chat())
		if err != nil {
			return fmt.Errorf("chatMemberOf: %w", err)
		}
		if member.Member {
			return next(c)
		}

		link, err := c.Bot().InviteLink(s.chat)
		if err != nil {
			log.Printf("inviteLink: %s", err)
		}
		menu := &tele.ReplyMarkup{}
		if link != "" {
			menu.Inline(menu.Row(menu.URL("Subscribe", link)))
		} else {
			menu.Inline(menu.Row(menu.Data("Ok", "confirm")))
		}
		return c.Send(s.askMessage, menu)
	}
}
