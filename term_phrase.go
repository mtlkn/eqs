package eqs

import (
	"errors"
	"strconv"
	"strings"
)

type PhraseTerm struct {
	Query     string
	Proximity int
}

func NewPhraseTerm(q string) *PhraseTerm {
	return &PhraseTerm{
		Query: q,
	}
}

func (term *PhraseTerm) Validate() error {
	if term == nil || term.Query == "" {
		return ErrEmpty
	}

	if term.Proximity < 0 {
		return errors.New("negative proximity")
	}

	return nil
}

func (term *PhraseTerm) String() string {
	q := term.Query
	if q == "" || q == "\"\"" {
		return ""
	}

	if q[0] != '"' {
		q = strconv.Quote(q)
	}

	if term.Proximity == 0 {
		return q
	}

	var sb strings.Builder

	sb.WriteString(q)
	sb.WriteByte('~')
	sb.Write([]byte(strconv.Itoa(term.Proximity)))

	return sb.String()
}
