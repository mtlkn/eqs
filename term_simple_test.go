package eqs

import "testing"

func TestSimpleTerm(t *testing.T) {
	term := NewSimpleTerm("quick")
	if term.String() != "quick" {
		t.Fail()
	}

	term.Fuzzy()
	if term.String() != "quick~" {
		t.Fail()
	}

	term.Fuzziness = 1
	if term.String() != "quick~1" {
		t.Fail()
	}

	term.Query = ""
	if err := term.Validate(); err != ErrEmpty {
		t.Fail()
	}
	if term.String() != "" {
		t.Fail()
	}

	term.Query = "quick brown fox"
	if err := term.Validate(); err.Error() != "multiple terms" {
		t.Fail()
	}

	term.Query = "quick"
	term.Fuzziness = -1
	if err := term.Validate(); err.Error() != "negative fuzziness" {
		t.Fail()
	}

	term = NewSimpleTerm("qu?ck")
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if !term.HasWildcards() || term.String() != "qu?ck" {
		t.Fail()
	}

	term = NewSimpleTerm("qu*")
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if !term.HasWildcards() || term.String() != "qu*" {
		t.Fail()
	}

	term = NewSimpleTerm("*uick")
	if err := term.Validate(); err != ErrLeadWildcard {
		t.Fail()
		t.Error(err)
	}

	term = NewSimpleTerm("*")
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}

	term = NewSimpleTerm("/joh?n(ath[oa]n)/")
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if !term.IsRegex() {
		t.Fail()
	}

	term = NewSimpleTerm("/joh?n(ath[oa]n)")
	if err := term.Validate(); err == nil {
		t.Fail()
	}
}
