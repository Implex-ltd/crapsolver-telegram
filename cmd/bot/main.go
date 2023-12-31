package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"bot/internal/invoice"
	"bot/internal/wrapper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

)

var botToken = "5953767481:AAHh6ZIpGj7fd-IL9U9T-CzEXbRjeUb10jg"

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "help":
			sendMessage(bot, update.Message.Chat.ID, "Everything is here: https://crapsolver.gitbook.io/readme/account/introduction")
		case "refill":
			handleRefillCommand(bot, update.Message)
		case "balance":
			handleGetUserCommand(bot, update.Message)
		default:
			sendMessage(bot, update.Message.Chat.ID, "I don't know that command, use /help")
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, messageText string) {
	if _, err := bot.Send(tgbotapi.NewMessage(chatID, messageText)); err != nil {
		log.Println("Error sending message:", err)
	}
}

func handleRefillCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	split := strings.Fields(message.Text)

	if len(split) < 2 {
		sendMessage(bot, message.Chat.ID, "Invalid command")
		return
	}

	account := ""
	if len(split) == 3 {
		account = split[2]
	}

	amount, err := strconv.Atoi(split[1])
	if err != nil {
		sendMessage(bot, message.Chat.ID, "Amount is invalid")
		return
	}

	if amount < 2500 {
		sendMessage(bot, message.Chat.ID, "Minimum order is 2500 credit (5 USD)")
		return
	}

	go func() {
		if account != "" {
			w := &wrapper.PrivateWrapper{}
			u, err := w.GetUser(account)
			if err != nil {
				sendMessage(bot, message.Chat.ID, "This api-key is invalid")
				return
			}

			sendMessage(bot, message.Chat.ID, fmt.Sprintf("⚠️ Following account gonna be refilled\n\n- ID %v\n- Current balance: %v", u.Data.ID, u.Data.Balance))
		} else {
			sendMessage(bot, message.Chat.ID, "A new account will be created and the api key sent after confirmation by the blockchain.")
		}

		inv, err := invoice.NewInvoice(message.From.UserName, amount)
		if err != nil {
			log.Println(err)
			return
		}

		sendMessage(bot, message.Chat.ID, fmt.Sprintf(`
		🕯️ %s
		
		- Invoice: %s
		- To pay: %v %v
		
		⚠️ This invoice is valid for 1 hours.
		`, inv.Data.Name, inv.Data.HostedURL, inv.Data.Pricing.Local.Amount, inv.Data.Pricing.Local.Currency))

		if success := invoice.WaitForOrder(inv.Data.ID, inv.Data.Code); success {
			handleSuccessfulOrder(bot, message, account, amount)
		} else {
			sendMessage(bot, message.Chat.ID, "Order failed")
		}
	}()
}

func handleSuccessfulOrder(bot *tgbotapi.BotAPI, message *tgbotapi.Message, account string, amount int) {
	w := &wrapper.PrivateWrapper{}
	uid := account

	if uid == "" {
		uid, err := w.CreateUser()
		if err != nil {
			sendMessage(bot, message.Chat.ID, "Failed to create api-key.")
			return
		}

		if uid == "" {
			sendMessage(bot, message.Chat.ID, "Error creating your API key")
			return
		}
	}

	balanceResponse, err := w.AddBalance(uid, amount)
	if err != nil {
		log.Printf("Error adding balance: %v\n", err)
		sendMessage(bot, message.Chat.ID, "Error adding balance")
	}

	if !balanceResponse.Success {
		sendMessage(bot, message.Chat.ID, fmt.Sprintf("Can't add credit to API key: %s, please contact support", uid))
	}

	u, err := w.GetUser(account)
	if err != nil {
		return
	}

	sendMessage(bot, message.Chat.ID, fmt.Sprintf("✅ Thanks !, account: %v, balance: %v", uid, u.Data.Balance))
}

func handleGetUserCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	split := strings.Fields(message.Text)

	if len(split) < 2 {
		sendMessage(bot, message.Chat.ID, "Invalid command")
		return
	}

	w := &wrapper.PrivateWrapper{}
	u, err := w.GetUser(split[1])
	if err != nil {
		sendMessage(bot, message.Chat.ID, "This api-key is invalid")
		return
	}

	sendMessage(bot, message.Chat.ID, fmt.Sprintf("- ID %v\n- Current balance: %v", u.Data.ID, u.Data.Balance))
}
