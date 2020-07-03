package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/proxy"
	tb "gopkg.in/tucnak/telebot.v2"
)

var bot *tb.Bot

func main() {
	fmt.Println("Make the bot exist. :)")

	var httpClient *http.Client
	httpClient, err := initSocks5Client()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Start creating bot ...")
	bot, err := tb.NewBot(tb.Settings{
		Token:  "1307117789:AAFqYg88j-R4dykmjY_J0FxAHKG8Cf306is",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Client: httpClient,
	})
	if err != nil {
		log.Fatalf("Cannot start bot. Error: %v\n", err)
	}

	bot.Handle("/start", func(m *tb.Message) {
		msg := "Testing the bot"
		if _, err = bot.Send(m.Chat, msg); err != nil {
			log.Println(err)
		}
		log.Println("Start request has been sent by: %v\n in chat: %v", m.Sender, m.Chat)
	})

	log.Println("Bot started")
	go func() {
		bot.Start()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	log.Println("Shutdown signal received, exiting...")
}

func initSocks5Client() (*http.Client, error) {
	addr := fmt.Sprintf("%s:%s", "127.0.0.1", "9350")
	dialer, err := proxy.SOCKS5("tcp", addr, &proxy.Auth{User: " ", Password: " "}, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("cannot init socks5 proxy client dialer: %w", err)
	}

	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial

	return httpClient, nil
}
