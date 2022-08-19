package eqs

import (
	"strings"
)

type GroupTerm struct {
	Query *Bool
}

func NewGroupTerm(op string, q ...string) *GroupTerm {
	b := new(Bool)

	for i, s := range q {
		tok := NewBoolToken(NewSimpleTerm(s))
		if op != "" && i > 0 {
			if op == "OR" {
				tok.OR = true
			} else if op == "AND" {
				tok.AND = true
			}
		}
		b.Tokens = append(b.Tokens, tok)
	}

	return &GroupTerm{
		Query: b,
	}
}

func (term *GroupTerm) Validate() error {
	if term == nil {
		return ErrEmpty
	}

	err := term.Query.Validate()
	if err != nil {
		return err
	}

	for _, tok := range term.Query.Tokens {
		switch tok.Query.(type) {
		case *Term:
			q := tok.Query.(*Term)
			if q.Field != "" {
				return ErrGroupValue
			}

			switch q.Query.(type) {
			case *SimpleTerm, *PhraseTerm:
				// ok
			default:
				return ErrGroupValue
			}
		case *SimpleTerm, *PhraseTerm:
			// ok
		default:
			return ErrGroupValue
		}
	}

	return nil
}

func (term *GroupTerm) String() string {
	var sb strings.Builder
	sb.WriteByte('(')
	sb.WriteString(term.Query.String())
	sb.WriteByte(')')
	return sb.String()
}
