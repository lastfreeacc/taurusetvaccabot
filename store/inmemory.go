package store

import (
	"sync"
)

type store struct {
	sync.Mutex
	gameNextVal int64
	users       map[int64]User
	games       map[int64]Game
}

func (s *store) SaveUser(u *User) {
	s.Lock()
	id := u.telegramProfile.UserID
	s.users[id] = *u
	s.Unlock()
}

func (s *store) SaveGame(g *Game) *Game {
	s.Lock()
	if g.ID == 0 {
		g.ID = s.gameNextVal
		s.gameNextVal++
	}
	s.games[g.ID] = *g
	s.Unlock()
	return g
}

// NewInMemory uses map for persistent objects
func NewInMemory() Store {
	store := store{
		gameNextVal: 1,
		users:       make(map[int64]User, 0),
		games:       make(map[int64]Game, 0),
	}
	return &store
}
