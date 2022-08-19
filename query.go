package eqs

import (
	"errors"
	"strings"
)

var (
	ErrEmpty        = errors.New("empty")
	ErrLeadWildcard = errors.New("leading wildcards")
	ErrRangeValue   = errors.New("range value")
	ErrGroupValue   = errors.New("group value")
	ErrBadField     = errors.New("bad field")
	ErrBadPrefix    = errors.New("bad prefix")
	ErrBadSuffix    = errors.New("bad suffix")
)

type Query interface {
	Validate() error
	String() string
}

func queryString(q Query) string {
	var s string

	switch q.(type) {
	case *Bool:
		var sb strings.Builder
		sb.WriteByte('(')
		sb.WriteString(q.String())
		sb.WriteByte(')')
		s = sb.String()
	default:
		s = q.String()
	}

	return s
}
