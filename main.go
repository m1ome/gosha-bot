package main

import (
	"fmt"
	"time"
	"flag"
	"strings"
	"os"

	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"github.com/m1ome/gosha-bot/library/draw"
	"github.com/satori/go.uuid"
)

var (
	logger *zap.SugaredLogger
	telegramKey string
	imagePath string
	fontPath string
)


var (
	Version = "dev"
)

func init() {
	// Initialize logger
	zlog, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Error initializing logger: %v", err))
	}

	logger = zlog.Sugar()

	// Reading flags
	flag.StringVar(&telegramKey, "key", "", "Telegram API key")
	flag.StringVar(&imagePath, "image", "assets/image.png", "Image path")
	flag.StringVar(&fontPath, "font", "assets/font.ttf", "Font path")
	flag.Parse()
}

func main() {
	if telegramKey == "" {
		logger.Fatalf("Please provide key")
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  telegramKey,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		logger.Fatalf("Error connecting bot", err)
	}

	bot.Handle("/version", func(m *tb.Message) {
		if m.Private() {
			_, err := bot.Send(m.Sender, fmt.Sprintf("Version: %s", Version))
			if err != nil {
				logger.Error("Error sending version", err)
			}
		}
	})

	bot.Handle("/text", func(m *tb.Message) {
		output := uuid.NewV4().String()
		message := m.Text
		message = strings.Replace(message, "/text", "", 1)
		message = strings.Replace(message, "@george_text_bot", "", 1)
		message = strings.TrimSpace(message)

		if len(message) == 0 {
			return
		}

		if err := draw.Text(imagePath, fontPath, message, output) ; err != nil {
			logger.Error("Error creating image", err)

			bot.Send(m.Sender, "Sorry there was an error creating your image")
			return
		}

		image := &tb.Photo{File: tb.FromDisk(output)}
		defer os.Remove(output)

		var receiver tb.Recipient
		if m.FromGroup() {
			logger.Info("Got request from group", m.Sender, message)

			receiver = m.Chat
		} else {
			logger.Info("Got request from private", m.Sender, message)

			receiver = m.Sender
		}

		sent, err := bot.Send(receiver, image)
		if err != nil {
			logger.Error("Error sending message", err)
		}

		logger.Infof("Send message back to %s: #%v", receiver.Recipient(), sent.ID)
	})

	logger.Info("Starting bot")
	bot.Start()
}