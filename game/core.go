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

const (
	ox  = rune(128003)
	cow = rune(128004)
)

// Game ...
type Game interface {
	Play(gamers map[int64]chan *teleapi.Update)
}

type game struct {
	bot          teleapi.Bot
	owner        teleapi.User
	caller       teleapi.User
	ownerNumber  string
	callerNumber string
	ownerCh      chan *teleapi.Update
	callerCh     chan *teleapi.Update
}

func (g *game) Play(gamers map[int64]chan *teleapi.Update) {
	ownerCh := make(chan string, 10)
	callerCh := make(chan string, 10)
	go g.toOwnerSender(ownerCh)
	go g.toCallerSender(callerCh)

	sendGameRequest(g.bot, g.caller.ID)
	for u := range g.callerCh {
		if u.CallbackQuery.Data == "" {
			continue
		}
		if u.CallbackQuery.Data == "No" {
			ownerCh <- locale["u_friend_decline_call"][g.owner.GetLanguage()]
			return
		}
		if u.CallbackQuery.Data == "Yes" {
			break
		}
	}

	// TODO: need to process yes/no answer

	ownerCh <- locale["wrote_ur_number_here"][g.owner.GetLanguage()]
	callerCh <- locale["wrote_ur_number_here"][g.caller.GetLanguage()]
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for u := range g.ownerCh {
			if isValidNumber(u.Message.Text) {
				g.ownerNumber = u.Message.Text
				wg.Done()
				break
			}
			ownerCh <- locale["wrong_number_try_again"][g.owner.GetLanguage()]
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
			callerCh <- locale["wrong_number_try_again"][g.caller.GetLanguage()]
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
	ownerCh <- locale["ur_move"][g.owner.GetLanguage()]
	callerCh <- locale["opponents_move"][g.caller.GetLanguage()]
	// timeout :=
	timer := time.NewTimer(20 * time.Minute)
	for {
		select {
		case update := <-g.ownerCh:
			if !isOwnerMove {
				ownerCh <- locale["not_ur_turn"][g.owner.GetLanguage()]
				continue
			}
			if !isValidNumber(update.Message.Text) {
				ownerCh <- locale["wrong_number_try_again"][g.owner.GetLanguage()]
				continue
			}
			t, c := countTandV(g.callerNumber, update.Message.Text)
			if t == 4 {
				ownerCh <- locale["u_win"][g.owner.GetLanguage()]
				callerCh <- locale["u_lose"][g.caller.GetLanguage()]
				callerCh <- locale["number_was"][g.caller.GetLanguage()] + g.ownerNumber
				delete(gamers, g.owner.ID)
				delete(gamers, g.caller.ID)
				timer.Stop()
				break
			}
			msg := fmt.Sprintf(string(ox)+" %d "+string(cow)+" %d", t, c)
			isOwnerMove = !isOwnerMove
			ownerCh <- msg
			callerCh <- locale["now_ur_turn"][g.caller.GetLanguage()]
		case update := <-g.callerCh:
			if isOwnerMove {
				callerCh <- locale["not_ur_turn"][g.owner.GetLanguage()]
				continue
			}
			if !isValidNumber(update.Message.Text) {
				callerCh <- locale["wrong_number_try_again"][g.caller.GetLanguage()]
				continue
			}
			t, c := countTandV(g.ownerNumber, update.Message.Text)
			if t == 4 {
				callerCh <- locale["u_win"][g.caller.GetLanguage()]
				ownerCh <- locale["u_lose"][g.owner.GetLanguage()]
				ownerCh <- locale["number_was"][g.owner.GetLanguage()] + g.callerNumber
				delete(gamers, g.owner.ID)
				delete(gamers, g.caller.ID)
				timer.Stop()
				break
			}
			msg := fmt.Sprintf(string(ox)+" %d "+string(cow)+" %d", t, c)
			isOwnerMove = !isOwnerMove
			callerCh <- msg
			ownerCh <- locale["now_ur_turn"][g.owner.GetLanguage()]
		case <-timer.C:
			ownerCh <- locale["timeout"][g.owner.GetLanguage()]
			callerCh <- locale["timeout"][g.caller.GetLanguage()]
			delete(gamers, g.owner.ID)
			delete(gamers, g.caller.ID)
			timer.Stop()
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
func New(bot teleapi.Bot, owner, caller teleapi.User, ownerCh chan *teleapi.Update, callerCh chan *teleapi.Update) (Game, error) {
	if owner.ID == 0 {
		log.Printf("[Warning] ownerID == 0")
		return nil, ErrBadUserID
	}
	if caller.ID == 0 {
		log.Printf("[Warning] callerID == 0")
		return nil, ErrBadUserID
	}
	return &game{bot: bot, owner: owner, caller: caller, ownerCh: ownerCh, callerCh: callerCh}, nil
}

func countTandV(n1, n2 string) (t, c int) {
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
		sendToPleer(g.bot, g.owner.ID, msg)
	}
}

func (g *game) toCallerSender(c chan string) {
	for {
		msg := <-c
		sendToPleer(g.bot, g.caller.ID, msg)
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

// code|locale|value tree
var locale = map[string]map[string]string{
	"ur_friend_decline_call": map[string]string{
		"en": "Ur friend decline call",
		"ru": "Друже не хочет играть",
	},
	"wrote_ur_number_here": map[string]string{
		"en": "Wrote ur number here",
		"ru": "Пиши свое число",
	},
	"wrong_number_try_again": map[string]string{
		"en": "Wrong number, try again",
		"ru": "Неправильное число, попробуй снова",
	},
	"ur_move": map[string]string{
		"en": "Ur_move",
		"ru": "Ходи!",
	},
	"opponents_move": map[string]string{
		"en": "Opponents move",
		"ru": "Ждем хода друже",
	},
	"not_ur_turn": map[string]string{
		"en": "Not ur turn",
		"ru": "Не спеши, сейчас ходит друже",
	},
	"u_win": map[string]string{
		"en": "U win",
		"ru": "ПОБЕДА!!! <3",
	},
	"u_lose": map[string]string{
		"en": "U lose",
		"ru": "Продул :(",
	},
	"number_was": map[string]string{
		"en": "Number was: ",
		"ru": "А число было: ",
	},
	"now_ur_turn": map[string]string{
		"en": "Now ur turn",
		"ru": "Твой ход!",
	},
	"timeout": map[string]string{
		"en": "Timeout",
		"ru": "Что-то очень долго, игра закончилась",
	},
}
