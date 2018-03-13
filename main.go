package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/lastfreeacc/taurusetvaccabot/teleapi"
)

type cmd string

const (
	confFilename = "taurusetvaccabot.conf.json"
)
const ( // commands
	startCmd  cmd = "/start"
	rulesCmd  cmd = "/r"
	playCmd   cmd = "/p"
	cancelCmd cmd = "/c"
)

var (
	conf     = make(map[string]interface{})
	botToken string
	bot      teleapi.Bot
)

func main() {
	myInit()
	startc := make(chan *teleapi.Update, 10)
	go startWorker(startc)
	
	rulesc := make(chan *teleapi.Update, 10)
	go rulesWorker(rulesc)

	gamec := make(chan *teleapi.Update, 10)
	go gameWorker(gamec)

	
	upCh := bot.Listen()
	for update := range upCh {
		log.Printf("[Debug] update message is: %#v\n", update)
		if update.Message.Contact != (teleapi.Contact{}) {
			doGame(update)
			continue
		}
		cmd := cmd(update.Message.Text)
		switch cmd {
		case startCmd:
			startc <- update
		case rulesCmd:
			rulesc <- update
		}
	}
}

func startWorker(updatec chan *teleapi.Update) {
	for update := range updatec {
		doStart(update)
	}
}

func doStart(update *teleapi.Update) {
	msg := fmt.Sprint(
		`Start...`)
	req := teleapi.SendMessageReq{ChatID: update.Message.Chat.ID, Text: msg}
	bot.SendMessage(req)
}

func rulesWorker(updatec chan *teleapi.Update) {
	for update := range updatec {
		doRules(update)
	}
}

func doRules(update *teleapi.Update) {
	msg := fmt.Sprint(
		`Rules...`)
	req := teleapi.SendMessageReq{ChatID: update.Message.Chat.ID, Text: msg}
	bot.SendMessage(req)
}

func gameWorker(updatec chan *teleapi.Update) {

}
func doGame(update *teleapi.Update) {
	contact := update.Message.Contact
	if contact.UserID == 0 {
		sendContactIsNotTelegramUser(update.Message.Chat.ID, contact)
		return
	}
	sendGameRequest(update)
}

func sendContactIsNotTelegramUser(chatID int64, contact teleapi.Contact) {
	userFullName := ""
	if contact.LastName != "" {
		userFullName = contact.LastName
	}
	if contact.FirstName != "" {
		if userFullName != "" {
			userFullName = userFullName + " "
		}
		userFullName = userFullName + contact.FirstName
	}
	if contact.PhoneNumber != "" {
		userFullName = userFullName + "(" + contact.PhoneNumber + ")"
	}
	if userFullName == "" {
		userFullName = "unknown"
	}
	msg := fmt.Sprintf("%s is not telegram user\n", userFullName)
	req := teleapi.SendMessageReq{ChatID: chatID, Text: msg}
	bot.SendMessage(req)
}

func sendGameRequest(update *teleapi.Update) {
	keyboard := yesNoKeyboard()
	req := teleapi.SendMessageReq{
		ChatID:      update.Message.Contact.UserID,
		Text:        update.Message.From.FirstName + " invites in bulls and cows game",
		ReplyMarkup: keyboard,
	}
	bot.SendMessage(req)
}

func yesNoKeyboard() *teleapi.InlineKeyboardMarkup {
	yes := teleapi.InlineKeyboardButton{
		Text:         "Yes",
		CallbackData: "Yes",
	}
	no := teleapi.InlineKeyboardButton{
		Text:         "No",
		CallbackData: "No",
	}
	keyboard := [][]teleapi.InlineKeyboardButton{{yes, no}}
	return &teleapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}
func myInit() {
	readMapFromJSON(confFilename, &conf)
	botToken, ok := conf["botToken"]
	if !ok || botToken == "" {
		log.Fatalf("[Error] can not find botToken in config file: %s\n", confFilename)
	}
	bot = teleapi.NewBot(botToken.(string))
}

func readMapFromJSON(filename string, mapVar *map[string]interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("[Error] can not read file '%s'\n", filename)
	}
	if err := json.Unmarshal(data, mapVar); err != nil {
		log.Fatalf("[Error] can not unmarshal json from file '%s'\n", filename)
	}
	log.Printf("[Info] read data from file: %s:\n%v\n", filename, mapVar)
}
