package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN required")
	}
	paypalLink := os.Getenv("PAYPAL_LINK")
	if paypalLink == "" {
		paypalLink = "https://www.paypal.com/paypalme/ducnhd"
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized as @%s", bot.Self.UserName)

	// Set bot menu commands
	commands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{Command: "start", Description: "Giới thiệu lotusmarket"},
		tgbotapi.BotCommand{Command: "donate", Description: "Ủng hộ dự án (Stars / PayPal)"},
		tgbotapi.BotCommand{Command: "help", Description: "Trợ giúp"},
	)
	if _, err := bot.Request(commands); err != nil {
		log.Printf("Failed to set commands: %v", err)
	} else {
		log.Println("Bot menu commands set")
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	log.Println("Listening for updates...")

	for update := range updates {
		// Handle pre-checkout queries
		if update.PreCheckoutQuery != nil {
			log.Printf("PreCheckout from %d", update.PreCheckoutQuery.From.ID)
			_, err := bot.Request(tgbotapi.PreCheckoutConfig{
				PreCheckoutQueryID: update.PreCheckoutQuery.ID,
				OK:                 true,
			})
			if err != nil {
				log.Printf("PreCheckout error: %v", err)
			}
			continue
		}

		// Handle callback queries (inline button presses)
		if update.CallbackQuery != nil {
			log.Printf("Callback: %s from %d", update.CallbackQuery.Data, update.CallbackQuery.From.ID)
			bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))

			chatID := update.CallbackQuery.Message.Chat.ID
			var amount int
			switch update.CallbackQuery.Data {
			case "donate_50":
				amount = 50
			case "donate_100":
				amount = 100
			case "donate_500":
				amount = 500
			default:
				continue
			}

			invoice := tgbotapi.InvoiceConfig{
				BaseChat:            tgbotapi.BaseChat{ChatID: chatID},
				Title:               fmt.Sprintf("Ủng hộ lotusmarket %d Stars", amount),
				Description:         fmt.Sprintf("Ủng hộ %d Stars cho dự án lotusmarket", amount),
				Payload:             fmt.Sprintf("donate_%d", amount),
				Currency:            "XTR",
				Prices:              []tgbotapi.LabeledPrice{{Label: "Donate", Amount: amount}},
				SuggestedTipAmounts: []int{},
			}
			if _, err := bot.Send(invoice); err != nil {
				log.Printf("Invoice send error: %v", err)
			}
			continue
		}

		// Handle messages
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		log.Printf("Msg from %s (%d): %s", update.Message.From.UserName, chatID, update.Message.Text)

		// Handle successful payment
		if update.Message.SuccessfulPayment != nil {
			amount := update.Message.SuccessfulPayment.TotalAmount
			send(bot, chatID, fmt.Sprintf("Cảm ơn bạn đã ủng hộ %d ⭐!\nlotusmarket sẽ tiếp tục phát triển. 🪷", amount))
			continue
		}

		switch update.Message.Command() {
		case "start":
			send(bot, chatID,
				"🪷 lotusmarket — Vietnamese Stock Market Toolkit\n\n"+
					"Thư viện mã nguồn mở cho thị trường chứng khoán Việt Nam.\n"+
					"Hỗ trợ Go (go get) và Python (pip install).\n\n"+
					"📦 GitHub: https://github.com/ducnhd/lotusmarket\n"+
					"🐍 PyPI: https://pypi.org/project/lotusmarket/\n\n"+
					"Dùng /donate để ủng hộ dự án!")

		case "donate":
			msg := tgbotapi.NewMessage(chatID,
				"Cảm ơn bạn muốn ủng hộ lotusmarket! 🪷\n\nChọn cách ủng hộ:")
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("⭐ 50 Stars", "donate_50"),
					tgbotapi.NewInlineKeyboardButtonData("⭐ 100 Stars", "donate_100"),
					tgbotapi.NewInlineKeyboardButtonData("⭐ 500 Stars", "donate_500"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("💳 PayPal", paypalLink),
				),
			)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Send error: %v", err)
			}

		case "help":
			send(bot, chatID,
				"📋 Các lệnh:\n\n"+
					"/start — Giới thiệu lotusmarket\n"+
					"/donate — Ủng hộ dự án (Stars / PayPal)\n"+
					"/help — Trợ giúp\n\n"+
					"GitHub: https://github.com/ducnhd/lotusmarket")

		default:
			if update.Message.Text != "" {
				send(bot, chatID, "Dùng /start để xem giới thiệu hoặc /donate để ủng hộ dự án 🪷")
			}
		}
	}
}

func send(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.DisableWebPagePreview = true
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Send error to %d: %v", chatID, err)
	} else {
		log.Printf("Sent to %d: %.50s...", chatID, text)
	}
}
