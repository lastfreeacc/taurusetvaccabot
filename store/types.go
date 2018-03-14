package store

type User struct {
	telegramProfile TelegramProfile
	// gameProfile     GameProfile
}

type TelegramProfile struct {
	UserID    int64
	FirstName string
	LastName  string
	UserName  string
}

// type GameProfile struct {
// 	BotTelegramChatID int64
// }

type Game struct {
	ID       int64
	OwnerID  int64
	CallerID int64
}

type Store interface {
	SaveGame(g *Game) *Game
	SaveUser(u *User)
}
