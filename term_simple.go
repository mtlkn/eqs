package eqs

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type SimpleTerm struct {
	Query     string
	Fuzziness int
	wc        bool // wildcards
	rgx       bool // regex
}

func NewSimpleTerm(q string) *SimpleTerm {
	return &SimpleTerm{
		Query: q,
	}
}

// sets to default fuzziness of 2
func (term *SimpleTerm) Fuzzy() {
	term.Fuzziness = 2
}

// needs to call Validate()
func (term *SimpleTerm) HasWildcards() bool {
	return term.wc
}

// needs to call Validate()
func (term *SimpleTerm) IsRegex() bool {
	return term.rgx
}

func (term *SimpleTerm) Validate() error {
	if term == nil || term.Query == "" {
		return ErrEmpty
	}

	var (
		wc, ws, esc bool
		rgx         int
	)

	for _, r := range term.Query {
		if unicode.IsSpace(r) && !esc {
			ws = true
			break
		}

		if (r == '*' || r == '?') && !esc {
			wc = true
			continue
		}

		if r == '/' && !esc {
			rgx++
		}

		esc = r == '\\'
	}

	if ws {
		return errors.New("multiple terms")
	}

	if wc && term.Query[0] == '*' && len(term.Query) > 1 {
		return ErrLeadWildcard
	}

	if rgx > 0 {
		if rgx != 2 || term.Query[0] != '/' || term.Query[len(term.Query)-1] != '/' {
			return errors.New("regex")
		}
		term.rgx = true
	}

	if term.Fuzziness < 0 {
		return errors.New("negative fuzziness")
	}

	term.wc = wc

	return nil
}

func (term *SimpleTerm) String() string {
	if term.Query == "" {
		return ""
	}

	if term.Fuzziness == 0 {
		return term.Query
	}

	var sb strings.Builder

	sb.WriteString(term.Query)
	sb.WriteByte('~')

	if term.Fuzziness != 2 {
		sb.Write([]byte(strconv.Itoa(term.Fuzziness)))
	}

	return sb.String()
}
