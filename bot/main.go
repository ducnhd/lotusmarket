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
		paypalLink = "https://paypal.me/lotusmarket"
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized as @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// Handle pre-checkout queries
		if update.PreCheckoutQuery != nil {
			bot.Send(tgbotapi.PreCheckoutConfig{
				PreCheckoutQueryID: update.PreCheckoutQuery.ID,
				OK:                 true,
			})
			continue
		}

		// Handle callback queries (inline button presses)
		if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			bot.Send(callback)

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
				BaseChat:    tgbotapi.BaseChat{ChatID: chatID},
				Title:       "Ủng hộ lotusmarket",
				Description: fmt.Sprintf("Ủng hộ %d Stars cho dự án lotusmarket", amount),
				Payload:     fmt.Sprintf("donate_%d", amount),
				Currency:    "XTR",
				Prices:      []tgbotapi.LabeledPrice{{Label: "Donate", Amount: amount}},
			}
			bot.Send(invoice)
			continue
		}

		// Handle messages
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID

		// Handle successful payment
		if update.Message.SuccessfulPayment != nil {
			amount := update.Message.SuccessfulPayment.TotalAmount
			msg := tgbotapi.NewMessage(chatID,
				fmt.Sprintf("Cảm ơn bạn đã ủng hộ %d ⭐!\nlotusmarket sẽ tiếp tục phát triển. 🪷", amount))
			bot.Send(msg)
			continue
		}

		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(chatID,
				"🪷 *lotusmarket* — Vietnamese Stock Market Toolkit\n\n"+
					"Thư viện mã nguồn mở cho Go & Python.\n"+
					"Dùng /donate để ủng hộ dự án!")
			msg.ParseMode = "Markdown"
			bot.Send(msg)
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
			bot.Send(msg)
		case "help":
			msg := tgbotapi.NewMessage(chatID,
				"/start — Giới thiệu\n"+
					"/donate — Ủng hộ dự án\n"+
					"/help — Trợ giúp")
			bot.Send(msg)
		}
	}
}
