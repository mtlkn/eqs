package eqs

import "testing"

func TestTerm(t *testing.T) {
	var term *Term
	if err := term.Validate(); err != ErrEmpty {
		t.Fail()
	}

	st := NewSimpleTerm("quick")
	term = NewTerm(st, "name")
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if term.String() != "name:quick" {
		t.Fail()
	}

	term.Boost = 2
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if term.String() != "name:quick^2" {
		t.Fail()
	}

	st.Fuzzy()
	if term.String() != "name:quick~^2" {
		t.Fail()
	}

	st.Fuzziness = 0
	term.Field = ""
	if term.String() != "quick^2" {
		t.Fail()
	}

	term.Boost = 0.25
	if term.String() != "quick^0.2" {
		t.Fail()
	}

	st.Query = ""
	if err := term.Validate(); err != ErrEmpty {
		t.Fail()
	}

	st.Query = "quick brown fox"
	if err := term.Validate(); err == nil {
		t.Fail()
	}

	st.Query = "quick"
	term.Boost = -1
	if err := term.Validate(); err == nil {
		t.Fail()
	}

	pt := NewPhraseTerm("quick brown fox")
	pt.Proximity = 2
	term = NewTerm(pt, "name")
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if term.String() != `name:"quick brown fox"~2` {
		t.Fail()
	}

	term = NewTerm(NewSimpleTerm("*"), "")
	if term.String() != "" {
		t.Fail()
	}
}
