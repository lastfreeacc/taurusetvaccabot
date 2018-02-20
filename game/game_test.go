package game

import "testing"

func TestIsValidNumber(t *testing.T) {
	valid := [...]string{"1234", "5678", "1357"}
	invalid := [...]string{"abc", "12345", "123d", "1233", "1232", "1111"}
	for _, v := range valid {
		if !isValidNumber(v) {
			t.Errorf("%s must be valid number", v)
		}
	}
	for _, inval := range invalid {
		if isValidNumber(inval) {
			t.Errorf("%s must be invalid number", inval)
		}
	}
}
