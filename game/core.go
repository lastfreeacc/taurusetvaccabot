package game

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/lastfreeacc/taurusetvaccabot/teleapi"
)

// Game ...
type Game interface {
	Play()
}

type game struct {
	bot          teleapi.Bot
	ownerID      int64
	callerID     int64
	ownerNumber  string
	callerNumber string
	ownerCh      chan *teleapi.Update
	callerCh     chan *teleapi.Update
}

func (g *game) Play() {
	ownerCh := make(chan string, 10)
	callerCh := make(chan string, 10)
	go g.toOwnerSender(ownerCh)
	go g.toCallerSender(callerCh)

	sendGameRequest(g.bot, g.callerID)
	for u := range g.callerCh {
		if u.CallbackQuery.Data == "" {
			continue
		}
		if u.CallbackQuery.Data == "No" {
			ownerCh <- "u friend decline call"
			return
		}
		if u.CallbackQuery.Data == "Yes" {
			break
		}
	}

	// TODO: need to process yes/no answer

	ownerCh <- "загадай четырехзначное число"
	callerCh <- "загадай четырехзначное число"
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for u := range g.ownerCh {
			if isValidNumber(u.Message.Text) {
				g.ownerNumber = u.Message.Text
				wg.Done()
				break
			}
			ownerCh <- "wrong number"
			continue
		}
	}()
	go func() {
		for u := range g.callerCh {
			if isValidNumber(u.Message.Text) {
				g.callerNumber = u.Message.Text
				wg.Done()
				break
			}
			callerCh <- "wrong number"
			continue
		}
	}()
	wg.Wait()
	// all ok!
	// lets game starts!!!
	// owner goes first
	// TODO: throw dice

	isOwnerMove := true
	// true - owners move
	// false - caller move
	ownerCh <- "your move"
	callerCh <- "opponents move"
	timeout := time.After(20 * time.Minute)
	for {
		select {
		case update := <-g.ownerCh:
			if !isOwnerMove {
				ownerCh <- "not your turn"
				continue
			}
			if !isValidNumber(update.Message.Text) {
				ownerCh <- "not valid numbrer\ntry again"
				continue
			}
			t, c := countTandC(g.ownerNumber, update.Message.Text)
			if t == 4 {
				ownerCh <- "you win"
				callerCh <- "you lose"
				break
			}
			msg := fmt.Sprintf("t: %d, c: %d", t, c)
			isOwnerMove = !isOwnerMove
			ownerCh <- msg
		case update := <-g.callerCh:
			if isOwnerMove {
				callerCh <- "not your turn"
				continue
			}
			if !isValidNumber(update.Message.Text) {
				callerCh <- "not valid numbrer\ntry again"
				continue
			}
			t, c := countTandC(g.callerNumber, update.Message.Text)
			if t == 4 {
				callerCh <- "you win"
				ownerCh <- "you lose"
				break
			}
			msg := fmt.Sprintf("t: %d, c: %d", t, c)
			isOwnerMove = !isOwnerMove
			callerCh <- msg
		case <-timeout:
			ownerCh <- "timeout"
			callerCh <- "timeout"
			break
		}
	}

}

var (
	// ErrBadUserID means Bad user id
	ErrBadUserID = errors.New("Bad user id")
	// ErrBadNumber ...
	// TODO: never used
	ErrBadNumber = errors.New("Bad number")
)

// New creates Game
func New(bot teleapi.Bot, ownerID, callerID int64, ownerCh chan *teleapi.Update, callerCh chan *teleapi.Update) (Game, error) {
	if ownerID == 0 {
		log.Printf("[Warning] ownerID == 0")
		return nil, ErrBadUserID
	}
	if callerID == 0 {
		log.Printf("[Warning] callerID == 0")
		return nil, ErrBadUserID
	}
	return &game{bot: bot, ownerID: ownerID, callerID: callerID, ownerCh: ownerCh, callerCh: callerCh}, nil
}

func countTandC(n1, n2 string) (t, c int) {
	for i, d1 := range n1 {
		for j, d2 := range n2 {
			if d1 == d2 {
				if i == j {
					t++
				} else {
					c++
				}
				continue
			}
		}
	}
	return
}

func isValidNumber(str string) bool {
	if len(str) != 4 {
		log.Printf("[Warning] len(%s) != 4\n", str)
		return false
	}
	_, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("[Warning] %s is not number\n", str)
		return false
	}
	for i, s := range str {
		for _, ss := range str[i+1:] {
			if s == ss {
				log.Printf("[Warning] %s has same digit '%c'\n", str, s)
				return false
			}
		}
	}
	return true
}

// func (g game) diceFirstMove() bool {
// }

func sendToPleer(bot teleapi.Bot, pleerID int64, msg string) {
	if pleerID == 0 {
		return
	}
	if msg == "" {
		return
	}
	req := teleapi.SendMessageReq{ChatID: pleerID, Text: msg}
	bot.SendMessage(req)
}

func (g *game) toOwnerSender(c chan string) {
	for {
		msg := <-c
		sendToPleer(g.bot, g.ownerID, msg)
	}
}

func (g *game) toCallerSender(c chan string) {
	for {
		msg := <-c
		sendToPleer(g.bot, g.callerID, msg)
	}
}

func sendGameRequest(bot teleapi.Bot, callerID int64) {
	keyboard := yesNoKeyboard()
	req := teleapi.SendMessageReq{
		ChatID:      callerID,
		Text:        "You friend invites you in bulls and cows game\nGo?",
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
