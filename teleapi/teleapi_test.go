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
