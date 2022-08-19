package eqs

import (
	"errors"
	"strings"
)

type Bool struct {
	Tokens     []*BoolToken
	IsFreeText bool // has no field queries
}

func (b *Bool) Validate() error {
	if b == nil || len(b.Tokens) == 0 {
		return ErrEmpty
	}

	for i, tok := range b.Tokens {
		err := tok.Validate()
		if err != nil {
			return err
		}

		if tok.Plus && tok.Minus {
			return errors.New("+-")
		}

		if i == 0 {
			if tok.AND {
				return errors.New("AND")
			}

			if tok.OR {
				return errors.New("OR")
			}
		}
	}

	return nil
}

func (b *Bool) String() string {
	var sb strings.Builder

	for i, tok := range b.Tokens {
		if i > 0 {
			sb.WriteByte(' ')
		}

		if tok.AND {
			sb.WriteString("AND ")
		} else if tok.OR {
			sb.WriteString("OR ")
		}

		if tok.NOT {
			sb.WriteString("NOT ")
		}

		if tok.Plus {
			sb.WriteByte('+')
		} else if tok.Minus {
			sb.WriteByte('-')
		}

		sb.WriteString(tok.String())
	}

	return sb.String()
}

type BoolToken struct {
	Query Query
	AND   bool // AND v
	OR    bool // OR v
	NOT   bool // NOT v
	Plus  bool // +v +f:v
	Minus bool // -v -f:v
}

func NewBoolToken(q Query) *BoolToken {
	return &BoolToken{
		Query: q,
	}
}

func (tok *BoolToken) Validate() error {
	err := tok.Query.Validate()
	if err != nil {
		return err
	}

	q := tok.String()
	if q == "+" || q == "-" || q == "AND" || q == "OR" || q == "NOT" {
		return errors.New(q)
	}

	return nil
}

func (tok *BoolToken) String() string {
	return queryString(tok.Query)
}
