package store

type User struct {
	telegramProfile TelegramProfile
	gameProfile     GameProfile
}

type TelegramProfile struct {
	UserID    int64
	FirstName string
	LastName  string
	UserName  string
}

type GameProfile struct {
	BotTelegramChatID int64
}
