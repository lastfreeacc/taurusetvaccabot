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
	startCmd cmd = "/start"
	rulesCmd cmd = "/r"
)

var (
	conf     = make(map[string]interface{})
	botToken string
	bot      teleapi.Bot
)

func main() {
	myInit()
	upCh := bot.Listen()
	for update := range upCh {
		cmd := cmd(update.Message.Text)
		switch cmd {
		case startCmd:
			doStart(update)
		case rulesCmd:
			doRules(update)

		}
	}
}

func doStart(update *teleapi.Update) {
	msg := fmt.Sprint(
		`Start...`)
	req := teleapi.SendMessageReq{ChatID: update.Message.Chat.ID, Text: msg}
	bot.SendMessage(req)
}

func doRules(update *teleapi.Update) {
	msg := fmt.Sprint(
		`Rules...`)
	req := teleapi.SendMessageReq{ChatID: update.Message.Chat.ID, Text: msg}
	bot.SendMessage(req)
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
		log.Fatalf("[Warning] can not read file '%s'\n", filename)
	}
	if err := json.Unmarshal(data, mapVar); err != nil {
		log.Fatalf("[Warning] can not unmarshal json from file '%s'\n", filename)
	}
	log.Printf("[Info] read data from file: %s:\n%v\n", filename, mapVar)
}
