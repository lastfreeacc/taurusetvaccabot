package game

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

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
	ownerCh  chan *teleapi.Update
	callerCh chan *teleapi.Update
}

func (g *game) Play() {
	ownerCh := make(chan string, 10)
	callerCh := make(chan string, 10)
	go g.toOwnerSender(ownerCh)
	go g.toCallerSender(callerCh)

	ownerCh <- "загадай четырехзначное число"
	callerCh <- "загадай четырехзначное число"

	go func() {
		for u := range g.ownerCh {
			if isValidNumber(u.Message.Text) {
				break
			}
			ownerCh <- "wrong number"
			continue
		}
	}()
	go func() {
		for u := range g.callerCh {
			if isValidNumber(u.Message.Text) {
				break
			}
			callerCh <- "wrong number"
			continue
		}
	}()

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
			t, c, err := countTandC(update.Message.Text)
			if err != nil {
				// TODO: ???
			}
			if t == 4 {
				ownerCh <- "you win"
				callerCh <- "you lose"
				break
			}
			msg := fmt.Sprintf("t: %d, c: %d", t, c)
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
			t, c, err := countTandC(update.Message.Text)
			if err != nil {
				// TODO: ???
			}
			if t == 4 {
				callerCh <- "you win"
				ownerCh <- "you lose"
				break
			}
			msg := fmt.Sprintf("t: %d, c: %d", t, c)
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

// func (g game) Listen() {
// 	go func() {
// 		for {
// 			u := <-g.gameCh
// 			switch u.Message.From.ID {
// 			case g.ownerID:
// 				break
// 			case g.callerID:
// 				break
// 			default:
// 				log.Printf("[Warning] strange Message.From.ID:%d, it is not owner:%d or caller:%d\n", u.Message.From.ID, g.ownerID, g.callerID)
// 				log.Printf("[Data] original Update is:%v\n", u)
// 			}
// 		}
// 	}()
// }
