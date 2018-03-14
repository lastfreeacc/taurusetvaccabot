package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/lastfreeacc/taurusetvaccabot/game"
	"github.com/lastfreeacc/taurusetvaccabot/store"
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
	merory   = store.NewInMemory()
	gamers   = make(map[int64]chan *teleapi.Update)
)

func main() {
	myInit()

	startCh := make(chan *teleapi.Update, 10)
	go startWorker(startCh)

	rulesCh := make(chan *teleapi.Update, 10)
	go rulesWorker(rulesCh)

	gameCh := make(chan *teleapi.Update, 10)
	go gameWorker(gameCh)

	upCh := bot.Listen()
	for update := range upCh {
		log.Printf("[Debug] update message is: %#v\n", update)
		switch update.Message.Command() {
		case "start":
			startCh <- update
		case "r":
			rulesCh <- update
		default:
			fromID := update.Message.From.ID
			if ch, ok := gamers[fromID]; ok {
				ch <- update
				continue
			}
			if update.Message.Contact != (teleapi.Contact{}) {
				doGame(update)
				continue
			}
		}
	}
}

func startWorker(updateCh chan *teleapi.Update) {
	for update := range updateCh {
		doStart(update)
	}
}

func doStart(update *teleapi.Update) {
	msg := fmt.Sprint(
		`Start...`)
	req := teleapi.SendMessageReq{ChatID: update.Message.Chat.ID, Text: msg}
	bot.SendMessage(req)
}

func rulesWorker(updateCh chan *teleapi.Update) {
	for update := range updateCh {
		doRules(update)
	}
}

func doRules(update *teleapi.Update) {
	msg := fmt.Sprint(
		`Rules...`)
	req := teleapi.SendMessageReq{ChatID: update.Message.Chat.ID, Text: msg}
	bot.SendMessage(req)
}

func gameWorker(updateCh chan *teleapi.Update) {

}
func doGame(update *teleapi.Update) {
	// game starts with sending Contact message
	// TODO: need to start with command with @username (need to test on telegram)

	contact := update.Message.Contact
	if contact.UserID == 0 {
		sendContactIsNotTelegramUser(update.Message.Chat.ID, contact)
		return
	}
	ownerID := update.Message.From.ID
	callerID := contact.UserID
	ownerCh := make(chan *teleapi.Update, 10)
	callerCh := make(chan *teleapi.Update, 10)
	gamers[ownerID] = ownerCh
	gamers[callerID] = callerCh
	// gameData := &store.Game{
	// 	OwnerID:  ownerID,
	// 	CallerID: callerID,
	// }
	// gameData = merory.SaveGame(gameData)
	game, err := game.New(bot, ownerID, callerID, ownerCh, callerCh)
	if err != nil {
		log.Printf("[Warning] can not create game for owner and caller (%d, %d)\n", ownerID, callerID)
		log.Printf("[Warning] error is: %s\n", err)
		// TODO: need to notify about error to users
	}
	go game.Play()
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
	msg := fmt.Sprintf("Sorry, %s is not telegram user", userFullName)
	req := teleapi.SendMessageReq{ChatID: chatID, Text: msg}
	bot.SendMessage(req)
}

func sendGameRequest(update *teleapi.Update) {
	keyboard := yesNoKeyboard()
	req := teleapi.SendMessageReq{
		ChatID:      update.Message.Contact.UserID,
		Text:        update.Message.From.FirstName + " invites you in bulls and cows game\nGo?",
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
