package eqs

import "testing"

func TestPhraseTerm(t *testing.T) {
	term := NewPhraseTerm("quick brown fox")
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	s := term.String()
	if s != `"quick brown fox"` {
		t.Fail()
	}

	term.Proximity = 5
	s = term.String()
	if s != `"quick brown fox"~5` {
		t.Fail()
	}

	term.Query = ""
	if err := term.Validate(); err != ErrEmpty {
		t.Fail()
	}
	s = term.String()
	if s != "" {
		t.Fail()
	}

	term.Query = "fox"
	term.Proximity = -1
	if err := term.Validate(); err.Error() != "negative proximity" {
		t.Fail()
	}
}
