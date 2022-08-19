package eqs

import (
	"testing"
)

func TestBool(t *testing.T) {
	b := new(Bool)
	if err := b.Validate(); err != ErrEmpty {
		t.Fail()
	}

	tok1 := NewTerm(NewSimpleTerm("quick"), "name")

	b.Tokens = append(b.Tokens, NewBoolToken(tok1))
	if err := b.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if b.String() != "name:quick" {
		t.Fail()
	}

	tok1.Boost = 2
	if b.String() != "name:quick^2" {
		t.Fail()
	}

	b.Tokens[0].AND = true
	if err := b.Validate(); err == nil {
		t.Fail()
	}

	b.Tokens[0].AND = false
	b.Tokens[0].OR = true
	if err := b.Validate(); err == nil {
		t.Fail()
	}

	b.Tokens[0].OR = false
	b.Tokens[0].NOT = true
	if err := b.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if b.String() != "NOT name:quick^2" {
		t.Fail()
	}

	b.Tokens[0].NOT = false
	b.Tokens[0].Minus = true
	if err := b.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if b.String() != "-name:quick^2" {
		t.Fail()
	}

	b.Tokens[0].Plus = true
	if err := b.Validate(); err == nil {
		t.Fail()
	}

	b.Tokens[0].Minus = false
	if err := b.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if b.String() != "+name:quick^2" {
		t.Fail()
	}

	b.Tokens[0].Plus = false
	tok2 := NewTerm(NewSimpleTerm("brown"), "name")
	tok2.Boost = 0.5
	b.Tokens = append(b.Tokens, NewBoolToken(tok2))
	b.Tokens[1].AND = true
	if err := b.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if b.String() != "name:quick^2 AND name:brown^0.5" {
		t.Fail()
	}

	b.Tokens[1].NOT = true
	if err := b.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if b.String() != "name:quick^2 AND NOT name:brown^0.5" {
		t.Fail()
	}

	b.Tokens[1].AND = false
	b.Tokens[1].NOT = false
	if err := b.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if b.String() != "name:quick^2 name:brown^0.5" {
		t.Fail()
	}

	b.Tokens[1].OR = true
	if err := b.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if b.String() != "name:quick^2 OR name:brown^0.5" {
		t.Fail()
	}

	b.Tokens[0].Query = NewTerm(NewSimpleTerm(""), "")
	if err := b.Validate(); err == nil {
		t.Fail()
	}

	b.Tokens[0].Query = NewTerm(NewSimpleTerm("AND"), "")
	if err := b.Validate(); err == nil {
		t.Fail()
	}

	tok3 := NewTerm(NewPhraseTerm("quick brown fox"), "story")
	b.Tokens = append(b.Tokens, NewBoolToken(tok3))
	b.Tokens = b.Tokens[1:]
	b.Tokens[0].OR = false
	b.Tokens[1].AND = true
	if err := b.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if b.String() != "name:brown^0.5 AND story:\"quick brown fox\"" {
		t.Fail()
	}

	q1 := new(Bool)
	q1.Tokens = append(q1.Tokens, &BoolToken{
		Query: &Term{
			Query: GTE("26"),
			Field: "age",
		},
		Plus: true,
	})
	q1.Tokens = append(q1.Tokens, &BoolToken{
		Query: &Term{
			Query: LTE("50"),
			Field: "age",
		},
		Minus: true,
	})

	q2 := new(Bool)
	q2.Tokens = append(q2.Tokens, &BoolToken{
		Query: &Term{
			Query: NewSimpleTerm("quick"),
			Field: "name",
			Boost: 2,
		},
	})
	q2.Tokens = append(q2.Tokens, &BoolToken{
		Query: &Term{
			Query: NewPhraseTerm("brown fox"),
			Field: "name",
		},
		AND: true,
	})

	b = new(Bool)
	b.Tokens = append(b.Tokens, &BoolToken{
		Query: q1,
	})
	b.Tokens = append(b.Tokens, &BoolToken{
		Query: q2,
		OR:    true,
	})
	if err := b.Validate(); err != nil {
		t.Fail()
		t.Error(err)
	}
	if b.String() != `(+age:>=26 -age:<=50) OR (name:quick^2 AND name:"brown fox")` {
		t.Fail()
	}

	b, err := Parse("headline:smith", nil)
	if err != nil {
		t.Fail()
		t.Error(err)
	}
	if b.IsFreeText {
		t.Fail()
	}

	b, err = Parse("headline:smith AND headline:trump usa", nil)
	if err != nil {
		t.Fail()
		t.Error(err)
	}
	if !b.IsFreeText {
		t.Fail()
	}

	b, err = Parse("type:text AND (headline:smith OR john) AND headline:trump", nil)
	if err != nil {
		t.Fail()
		t.Error(err)
	}
	if !b.IsFreeText {
		t.Fail()
	}

	b, err = Parse("(type:text OR type:photo) AND (headline:john title:smith)", nil)
	if err != nil {
		t.Fail()
		t.Error(err)
	}
	if b.IsFreeText {
		t.Fail()
	}

}
