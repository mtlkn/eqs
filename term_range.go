package eqs

import "strings"

type RangeTerm struct {
	Left  *RangeTermValue
	Right *RangeTermValue
	one   bool // gt, gte, lt, lte not [left TO right]
}

func NewRangeTerm(left, right *RangeTermValue) *RangeTerm {
	return &RangeTerm{
		Left:  left,
		Right: right,
	}
}

func (term *RangeTerm) Validate() error {
	if term == nil {
		return ErrEmpty
	}

	err := term.Left.Validate()
	if err != nil {
		return err
	}

	err = term.Right.Validate()
	if err != nil {
		return err
	}

	if term.Left.Open && term.Right.Open {
		return ErrEmpty
	}

	return nil
}

func (term *RangeTerm) String() string {
	var sb strings.Builder

	if !term.one {
		if term.Left.Include {
			sb.WriteByte('[')
		} else {
			sb.WriteByte('{')
		}

		sb.WriteString(term.Left.String())
		sb.WriteString(" TO ")
		sb.WriteString(term.Right.String())

		if term.Right.Include {
			sb.WriteByte(']')
		} else {
			sb.WriteByte('}')
		}

		return sb.String()
	}

	if !term.Left.Open {
		sb.WriteByte('>')
		if term.Left.Include {
			sb.WriteByte('=')
		}
		sb.WriteString(term.Left.String())
		return sb.String()
	}

	sb.WriteByte('<')
	if term.Right.Include {
		sb.WriteByte('=')
	}
	sb.WriteString(term.Right.String())

	return sb.String()
}

type RangeTermValue struct {
	Query   Query
	Include bool
	Open    bool // *
}

func NewRangeTermValue(q Query, eq bool) *RangeTermValue {
	return &RangeTermValue{
		Query:   q,
		Include: eq,
	}
}

func NewRangeStringValue(q string, eq bool) *RangeTermValue {
	return &RangeTermValue{
		Query:   NewSimpleTerm(q),
		Include: eq,
	}
}

func NewRangePhraseValue(q string, eq bool) *RangeTermValue {
	return &RangeTermValue{
		Query:   NewPhraseTerm(q),
		Include: eq,
	}
}

func NewOpenRange() *RangeTermValue {
	return &RangeTermValue{
		Open: true,
	}
}

func (v *RangeTermValue) Validate() error {
	if v == nil {
		return ErrEmpty
	}

	if v.Open {
		v.Include = true
		return nil
	}

	err := v.Query.Validate()
	if err != nil {
		return err
	}

	switch v.Query.(type) {
	case *SimpleTerm:
		q := v.Query.(*SimpleTerm)
		if (q.HasWildcards() && q.Query != "*") || q.IsRegex() || q.Fuzziness != 0 {
			return ErrRangeValue
		}
	case *PhraseTerm:
		// ok
	default:
		return ErrRangeValue
	}

	return nil
}

func (v *RangeTermValue) String() string {
	if v.Open {
		return "*"
	}
	return v.Query.String()
}

func GT(q string) *RangeTerm {
	return &RangeTerm{
		Left: &RangeTermValue{
			Query: NewSimpleTerm(q),
		},
		Right: NewOpenRange(),
		one:   true,
	}
}

func GTE(q string) *RangeTerm {
	return &RangeTerm{
		Left: &RangeTermValue{
			Query:   NewSimpleTerm(q),
			Include: true,
		},
		Right: NewOpenRange(),
		one:   true,
	}
}

func LT(q string) *RangeTerm {
	return &RangeTerm{
		Left: NewOpenRange(),
		Right: &RangeTermValue{
			Query: NewSimpleTerm(q),
		},
		one: true,
	}
}

func LTE(q string) *RangeTerm {
	return &RangeTerm{
		Left: NewOpenRange(),
		Right: &RangeTermValue{
			Query:   NewSimpleTerm(q),
			Include: true,
		},
		one: true,
	}
}
