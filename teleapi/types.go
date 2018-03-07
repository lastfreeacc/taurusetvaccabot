package teleapi

// Update ...
type Update struct {
	UpdateID      int64         `json:"update_id"`
	Message       Message       `json:"message,omitempty"`
	CallbackQuery CallbackQuery `json:"callback_query,omitempty"`
}

// Message ...
type Message struct {
	MessageID int64    `json:"message_id"`
	From      User     `json:"from,omitempty"`
	Date      int64    `json:"date"`
	Chat      Chat     `json:"chat"`
	Text      string   `json:"text,omitempty"`
	Entities  []Entity `json:"entities,omitempty"`
	Contact   Contact  `json:"contact,omitempty"`
}

// CallbackQuery represents an incoming callback query from a callback button in an inline keyboard.
// If the button that originated the query was attached to a message sent by the bot, the field message will be present.
// If the button was attached to a message sent via the bot (in inline mode), the field inline_message_id will be present.
// Exactly one of the fields data or game_short_name will be present
type CallbackQuery struct {
	ID              string  `json:"id"`
	From            User    `json:"from"`
	Message         Message `json:"Message,omitempty"`
	InlineMessageID string  `json:"inline_message_id,omitempy"`
	ChatInstance    string  `json:"chat_instance,omitempty"`
	Data            string  `json:"data,omitempty"`
	GameShortName   string  `json:"game_short_name,omitempty"`
}

// User ...
type User struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

// Chat ...
type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// EntityType ...
type EntityType int

func (e EntityType) String() string {
	return entities[e]
}

var entities = [...]string{
	"mention",
	"hashtag",
	"bot_command",
	"url",
	"email",
	"bold",
	"italic",
	"code",
	"pre",
	"text_link",
	"text_mention",
}

// ...
const (
	MentionEntity EntityType = iota
	HashtagEntity
	BotCommandEntity
	URLEntity
	EmailEntity
	BoldEntity
	ItalicEntity
	CodeEntity
	PreEntity
	TextLinkEntity
	TextMentionEntity
)

// Entity ...
type Entity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	URL    string `json:"url"`
	User   User   `json:"user"`
}

type getUpdatesResp struct {
	Ok     bool      `json:"ok"`
	Result []*Update `json:"result"`
}

// SendMessageReq ...
type SendMessageReq struct {
	ChatID                interface{} `json:"chat_id"` // string || number
	Text                  string      `json:"text"`
	ParseMode             bool        `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool        `json:"disable_web_page_preview,omitempty"`
	DisableNotification   bool        `json:"disable_notification,omitempty"`
	ReplyToMessageID      int64       `json:"reply_to_message_id,omitempty"`
	ReplyMarkup           interface{} `json:"reply_markup,omitempty"`
}

// KeyboardButton represents one button of the reply keyboard.
// For simple text buttons String can be used instead of this object to specify text of the button.
// Optional fields are mutually exclusive.
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact,omitempty"`
	RequestLocation bool   `json:"request_location,omitempty"`
}

// ReplyKeyboardMarkup represents a custom keyboard with reply options (see Introduction to bots for details and examples).
type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool               `json:"one_time_keyboard,omitempty"`
	Selective       bool               `josn:"selective,omitempty"`
}

// CallbackGame is a placeholder, currently holds no information. Use BotFather to set up your game.
type CallbackGame struct{}

// InlineKeyboardButton represents one button of an inline keyboard. You must use exactly one of the optional fields.
type InlineKeyboardButton struct {
	Text                        string       `json:"text"`
	URL                         string       `json:"url,omitempty"`
	CallbackData                string       `json:"callback_data,omitempty"`
	SwitchInlineQuery           string       `json:"switch_inline_query,omitempty"`
	SwitchInlinQueryCurrentChat string       `json:"switch_inline_query_current_chat,omitempty"`
	CallbackGame                CallbackGame `json:"callback_game,omitempty"`
	Pay                         bool         `json:"pay,omitempty"`
}

// InlineKeyboardMarkup ...
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// Contact represents a phone contact.
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	UserID      int64  `json:"user_id,omitempty"`
}
