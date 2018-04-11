package game

import "testing"

func Test_isValidNumber(t *testing.T) {
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

func Test_countTandV(t *testing.T) {
	type args struct {
		n1 string
		n2 string
	}
	tests := []struct {
		name  string
		n1    string
		n2    string
		wantT int
		wantC int
	}{
		{"t1",
			"1234",
			"1243",
			2,
			2},
		{"t2",
			"1357",
			"1357",
			4,
			0},
		{"t3",
			"1357",
			"2468",
			0,
			0},
		{"t4",
			"1245",
			"3278",
			1,
			0},
		{"t5",
			"12344",
			"4320",
			0,
			3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, gotC := countTandV(tt.n1, tt.n2)
			if gotT != tt.wantT {
				t.Errorf("countTandC() gotT = %v, want %v", gotT, tt.wantT)
			}
			if gotC != tt.wantC {
				t.Errorf("countTandC() gotC = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
