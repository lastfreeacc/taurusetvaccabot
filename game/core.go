package game

import (
	"errors"
	"log"
	"strconv"
)

type Game interface {
	Play()
}

type game struct {
	ownerID  int64
	callerID int64
}

func (g game) Play() {

}

var (
	ErrBadUserID = errors.New("Bad user id")
	ErrBadNumber = errors.New("Bad number")
)

func New(ownerID, callerID int64) (Game, error) {
	if ownerID == 0 || callerID == 0 {
		return nil, ErrBadUserID
	}
	return game{ownerID: ownerID, callerID: callerID}, nil
}

func countTandC(str string) (t, c int, err error) {
	_, err = strconv.Atoi(str)
	if err != nil {
		return 0, 0, err
	}
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
