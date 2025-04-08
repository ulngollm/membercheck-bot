package main

import (
	"github.com/ulngollm/membercheckbot/middleware"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("godotenv.Load: %s", err)
		return
	}
}

func main() {
	t, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatalf("bot token is empty")
		return
	}

	if !ok {
		log.Fatalf("chatID is empty")
		return
	}

	chID, ok := os.LookupEnv("SUBSCRIPTION_CHAT_ID")
	if !ok {
		log.Fatalf("SUBSCRIPTION_CHAT_ID is empty")
		return
	}
	chatID, err := strconv.Atoi(chID)
	if err != nil {
		log.Fatalf("strconv.Atoi: %s", err)
		return
	}
	message := os.Getenv("SUBSCRIPTION_ASK_MSG")
	if message == "" {
		log.Fatalf("SUBSCRIPTION_ASK_MSG is empty")
		return
	}

	pref := tele.Settings{
		Token:  t,
		Poller: &tele.LongPoller{Timeout: time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("tele.NewBot: %s", err)
		return
	}

	chat, err := bot.ChatByID(int64(chatID))
	if err != nil {
		log.Fatalf("tele.ChatByID: %s", err)
		return
	}

	g := bot.Group()
	g.Use(middleware.NewSubscriptionMiddleware(message, chat).CheckSubscription)
	g.Handle("/start", handle)
	g.Handle(tele.OnText, handle)

	bot.Start()
}

func handle(c tele.Context) error {
	return c.Send("Hello!")
}
