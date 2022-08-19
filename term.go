package eqs

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

type Term struct {
	Query Query
	Field string
	Boost float64
}

func NewTerm(q Query, field string) *Term {
	return &Term{
		Query: q,
		Field: field,
	}
}

func (term *Term) Validate() error {
	if term == nil {
		return ErrEmpty
	}

	err := term.Query.Validate()
	if err != nil {
		return err
	}

	if term.Boost < 0 {
		return errors.New("negative boost")
	}

	return nil
}

func (term *Term) String() string {
	var (
		sb strings.Builder
		q  = queryString(term.Query)
	)

	if q == "*" || q == "?" {
		if term.Field == "" {
			return ""
		}
		q = term.Field
		term.Field = "_exists_"
		term.Query = NewSimpleTerm(q)
	}

	if term.Field != "" {
		sb.WriteString(term.Field)
		sb.WriteByte(':')
	}

	sb.WriteString(q)

	if term.Boost > 0 && term.Boost != 1 {
		sb.WriteByte('^')

		if math.Ceil(term.Boost) == term.Boost {
			sb.WriteString(strconv.Itoa(int(term.Boost)))
		} else {
			sb.WriteString(strconv.FormatFloat(term.Boost, 'f', 1, 64))
		}
	}

	return sb.String()
}
