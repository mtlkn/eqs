package eqs

import "testing"

func TestGroupTerm(t *testing.T) {
	var term *GroupTerm
	if err := term.Validate(); err != ErrEmpty {
		t.Fail()
	}

	var b *Bool
	term = &GroupTerm{
		Query: b,
	}
	if err := term.Validate(); err != ErrEmpty {
		t.Fail()
	}

	b = new(Bool)
	b.Tokens = append(b.Tokens, &BoolToken{
		Query: NewPhraseTerm("quick brown fox"),
	})

	b.Tokens = append(b.Tokens, &BoolToken{
		Query: NewPhraseTerm("jumps over the fence"),
	})

	term.Query = b
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if term.String() != `("quick brown fox" "jumps over the fence")` {
		t.Fail()
	}

	term.Query.Tokens[1].OR = true
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if term.String() != `("quick brown fox" OR "jumps over the fence")` {
		t.Fail()
	}

	term.Query.Tokens[1].OR = false
	term.Query.Tokens[1].AND = true
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if term.String() != `("quick brown fox" AND "jumps over the fence")` {
		t.Fail()
	}

	b = new(Bool)
	b.Tokens = append(b.Tokens, &BoolToken{
		Query: NewSimpleTerm("quick"),
		Minus: true,
	})
	term = &GroupTerm{
		Query: b,
	}
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}

	b = new(Bool)
	b.Tokens = append(b.Tokens, &BoolToken{
		Query: &Term{
			Query: GT("quick"),
			Field: "name",
		},
	})
	term = &GroupTerm{
		Query: b,
	}
	if err := term.Validate(); err == nil {
		t.Fail()
	}

	b = new(Bool)
	b.Tokens = append(b.Tokens, &BoolToken{
		Query: GT("quick"),
	})
	term = &GroupTerm{
		Query: b,
	}
	if err := term.Validate(); err == nil {
		t.Fail()
	}

	term = NewGroupTerm("", "quick", "brown", "fox")
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if term.String() != "(quick brown fox)" {
		t.Fail()
	}

	term = NewGroupTerm("OR", "quick", "brown", "fox")
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if term.String() != "(quick OR brown OR fox)" {
		t.Fail()
	}

	term = NewGroupTerm("AND", "quick", "brown", "fox")
	if err := term.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if term.String() != "(quick AND brown AND fox)" {
		t.Fail()
	}

}
