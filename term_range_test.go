package eqs

import "testing"

func TestRangeTerm(t *testing.T) {
	var term *RangeTerm
	if err := term.Validate(); err != ErrEmpty {
		t.Fail()
	}

	term = NewRangeTerm(nil, nil)
	if err := term.Validate(); err != ErrEmpty {
		t.Fail()
	}

	term = NewRangeTerm(NewRangeStringValue("26", true), nil)
	if err := term.Validate(); err != ErrEmpty {
		t.Fail()
	}

	term = NewRangeTerm(NewOpenRange(), NewOpenRange())
	if err := term.Validate(); err != ErrEmpty {
		t.Fail()
	}

	term = NewRangeTerm(NewRangeStringValue("26", true), NewOpenRange())
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if term.String() != "[26 TO *]" {
		t.Fail()
	}

	term = NewRangeTerm(NewRangeStringValue("26", false), NewRangeStringValue("50", true))
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if term.String() != "{26 TO 50]" {
		t.Fail()
	}

	term = NewRangeTerm(NewRangeStringValue("26", false), NewRangeStringValue("50", false))
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if term.String() != "{26 TO 50}" {
		t.Fail()
	}

	q := GT("2020")
	if err := q.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if q.String() != ">2020" {
		t.Fail()
	}

	q = GTE("2020")
	if err := q.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if q.String() != ">=2020" {
		t.Fail()
	}

	q = LT("2020")
	if err := q.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if q.String() != "<2020" {
		t.Fail()
	}

	q = LTE("2020")
	if err := q.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if q.String() != "<=2020" {
		t.Fail()
	}

	v := NewRangeTermValue(NewPhraseTerm("quick brown fox"), true)
	if err := v.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}

	v = &RangeTermValue{
		Query: &SimpleTerm{
			Query: "quick brown fox",
		},
	}
	if err := v.Validate(); err == nil {
		t.Fail()
	}

	v = &RangeTermValue{
		Query: &SimpleTerm{
			Query: "/f[o]x/",
		},
	}
	if err := v.Validate(); err == nil {
		t.Fail()
	}

	v = &RangeTermValue{
		Query: NewRangeTerm(NewRangeStringValue("26", false), NewRangeStringValue("26", false)),
	}
	if err := v.Validate(); err == nil {
		t.Fail()
	}
}
