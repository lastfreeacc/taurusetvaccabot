package teleapi

import (
	"testing"
)

func TestMessageCommand(t *testing.T) {
	var commandTable = []struct {
		in  string
		out string
	}{
		{"just common message", ""},
		{"message with @ and /", ""},
		{"/command", "command"},
		{"/command@botName", "command"},
		{"/command args", "command"},
		{"/command@botName args", "command"},
	}
	upd := Update{
		Message: Message{},
	}
	for _, v := range commandTable {
		upd.Message.Text = v.in
		cmd := upd.Message.Command()
		if cmd != v.out {
			t.Errorf("expected: \"%s\" for Message.Text: \"%s\", but found: \"%s\"", v.out, v.in, cmd)
		}
	}
}

func Test_get2letters(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name         string
		languageCode string
		want         string
	}{
		{
			name:         "t1",
			languageCode: "en-us",
			want:         "en",
		},
		{
			name:         "t2",
			languageCode: "",
			want:         "en",
		},
		{
			name:         "t3",
			languageCode: "EN-us",
			want:         "en",
		},
		{
			name:         "t4",
			languageCode: "Ru-ru",
			want:         "ru",
		},
		{
			name:         "t5",
			languageCode: "russian",
			want:         "ru",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := get2letters(tt.languageCode); got != tt.want {
				t.Errorf("get2letters() = %v, want %v", got, tt.want)
			}
		})
	}
}
