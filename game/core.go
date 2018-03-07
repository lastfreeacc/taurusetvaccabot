package game

import (
	"errors"
	"log"
	"strconv"

	"github.com/lastfreeacc/taurusetvaccabot/teleapi"
)

// Game ...
type Game interface {
	Play()
}

type game struct {
	bot      teleapi.Bot
	ownerID  int64
	callerID int64
	gameCh   chan *teleapi.Update
}

func (g game) Play() {
	ownderCh := make(chan string)
	callerCh := make(chan string)
	go g.toOwnerSender(ownderCh)
	go g.toCallerSender(callerCh)

	ownderCh <- "загадай четырехзначное число"
	callerCh <- "загадай четырехзначное число"

}

var (
	// ErrBadUserID means Bad user id
	ErrBadUserID = errors.New("Bad user id")
	// ErrBadNumber ...
	// TODO: never used
	ErrBadNumber = errors.New("Bad number")
)

// New creates Game
func New(bot teleapi.Bot, ownerID, callerID int64, gameCh chan *teleapi.Update) (Game, error) {
	if ownerID == 0 {
		log.Printf("[Warning] ownerID == 0")
		return nil, ErrBadUserID
	}
	if callerID == 0 {
		log.Printf("[Warning] callerID == 0")
		return nil, ErrBadUserID
	}
	return game{bot: bot, ownerID: ownerID, callerID: callerID, gameCh: gameCh}, nil
}

func countTandC(str string) (t, c int, err error) {
	_, err = strconv.Atoi(str)
	if err != nil {
		return 0, 0, err
	}
	// TODO: need to implement
	return 0, 0, nil
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
// 	Rando
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

// func readFromPleer() {

// }

func (g game) toOwnerSender(c chan string) {
	for {
		msg := <-c
		sendToPleer(g.bot, g.ownerID, msg)
	}
}

func (g game) toCallerSender(c chan string) {
	for {
		msg := <-c
		sendToPleer(g.bot, g.callerID, msg)
	}
}

func (g game) Listen() {
	go func() {
		for {
			u := <-g.gameCh
			switch u.Message.From.ID {
			case g.ownerID:
				break
			case g.callerID:
				break
			default:
				log.Printf("[Warning] strange Message.From.ID:%d, it is not owner:%d or caller:%d\n", u.Message.From.ID, g.ownerID, g.callerID)
				log.Printf("[Data] original Update is:%v\n", u)
			}
		}
	}()
}
