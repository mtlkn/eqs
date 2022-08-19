package eqs

import (
	"strconv"
	"strings"

	"git.ap.org/golang/eqs/parsing"
)

func Parse(q string, aliases map[string][]string) (*Bool, error) {
	toks, err := parsing.Parse(q)
	if err != nil {
		return nil, err
	}

	b := new(Bool)

	for i, tok := range toks {
		if i > 0 && len(b.Tokens) == 0 {
			// previous tokens are empty (*, +*, etc.), remove AND or OR
			tok.AND = false
			tok.OR = false
		}

		err := b.appendToken(tok, aliases)
		if err != nil {
			return nil, err
		}
	}

	if len(b.Tokens) == 0 {
		return b, nil
	}

	err = b.Validate()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Bool) appendToken(tok *parsing.Token, aliases map[string][]string) error {
	if len(tok.Value) == 1 && tok.Type == parsing.TEXT && (tok.Value[0] == '*' || tok.Value[0] == '?') {
		return nil // * or ? or AND * or AND ? or OR * or OR ?
	}

	token := &BoolToken{
		AND: tok.AND,
		OR:  tok.OR,
		NOT: tok.NOT,
	}

	p := new(textParser)
	if tok.Type == parsing.TEXT {
		err := p.ProcessPrefix(tok.Value, false)
		if err != nil {
			return err
		}
		tok.Value = p.Value
	} else {
		err := p.ProcessPrefix(tok.Prefix, true)
		if err != nil {
			return err
		}

		err = p.ParseSuffix(tok.Suffix)
		if err != nil {
			return err
		}
	}

	if p.Fuzziness > 0 {
		if p.Op != "" {
			return ErrBadSuffix
		}

		switch tok.Type {
		case parsing.TEXT:
			// ok
		case parsing.PHRASE:
			if p.Proximity == 0 { // cannot be "john smith"~, only "john smith"~2
				return ErrBadSuffix
			}
		default:
			return ErrBadSuffix
		}
	}

	var (
		q   Query
		err error
	)

	if p.Op == "" {
		switch tok.Type {
		case parsing.TEXT:
			q, err = processText(tok, p.Field != "", p.Fuzziness)
		case parsing.PHRASE:
			q, err = processPhrase(tok, p.Proximity)
		case parsing.RANGE:
			if p.Field == "" {
				return ErrRangeValue
			}
			q, err = processRange(tok)
		case parsing.GROUP:
			if p.Field == "" {
				q, err = processGroup(tok, false, aliases, b)
			} else {
				q, err = processGroup(tok, true, nil, b)
			}
		case parsing.REGEX:
			q = NewSimpleTerm(string(tok.Value))
		}
	} else {
		q, err = processRangeOp(tok, p.Op)
	}

	if err != nil || q == nil {
		return err
	}

	if p.Field != "" && len(aliases) > 0 {
		fs, ok := aliases[p.Field]
		if ok {
			if len(fs) == 1 {
				p.Field = fs[0]
			} else {
				bq := new(Bool)

				for i, f := range fs {
					bq.Tokens = append(bq.Tokens, &BoolToken{
						Query: &Term{
							Query: q,
							Field: f,
							Boost: p.Boost,
						},
						OR: i > 0,
					})
				}

				token.Query = bq
			}
		}
	}

	if token.Query == nil {
		token.Query = &Term{
			Query: q,
			Field: p.Field,
			Boost: p.Boost,
		}
	}

	token.Plus = p.Plus
	token.Minus = p.Minus

	b.Tokens = append(b.Tokens, token)

	if p.Field == "" && tok.Type != parsing.GROUP {
		b.IsFreeText = true
	}

	return nil
}

func processRangeOp(tok *parsing.Token, op string) (Query, error) {
	var (
		q  *RangeTermValue
		s  = safeString(tok.Value)
		eq = op == ">=" || op == "<="
	)

	switch tok.Type {
	case parsing.TEXT:
		q = NewRangeStringValue(s, eq)
	case parsing.PHRASE:
		q = NewRangePhraseValue(s, eq)
	default:
		return nil, ErrBadField
	}

	var rt *RangeTerm

	switch op {
	case ">", ">=":
		rt = NewRangeTerm(q, NewOpenRange())
	case "<", "<=":
		rt = NewRangeTerm(NewOpenRange(), q)
	}
	rt.one = true

	return rt, nil
}

func processText(tok *parsing.Token, field bool, fuzziness int) (Query, error) {
	s := safeString(tok.Value)
	if (s == "*" || s == "?") && !field {
		return nil, nil
	}
	q := NewSimpleTerm(s)
	q.Fuzziness = fuzziness
	return q, nil
}

func processPhrase(tok *parsing.Token, proximity int) (Query, error) {
	if len(tok.Value) == 0 {
		return nil, nil
	}

	v := string(tok.Value)

	if v[0] != '"' {
		var sb strings.Builder
		sb.WriteByte('"')
		sb.WriteString(v)
		sb.WriteByte('"')
		v = sb.String()

		var err error
		v, err = strconv.Unquote(v)
		if err != nil {
			return nil, err
		}
	}

	v = strings.TrimSpace(v)
	if v == "" {
		return nil, nil
	}

	q := NewPhraseTerm(v)
	if proximity > 0 {
		q.Proximity = proximity
	}
	return q, nil
}

func processGroup(tok *parsing.Token, value bool, aliases map[string][]string, parent *Bool) (Query, error) {
	q, err := Parse(string(tok.Value), aliases)
	if err != nil {
		return nil, err
	}

	if len(q.Tokens) == 0 {
		if value {
			return nil, ErrEmpty
		}
		return nil, nil
	}

	if !value { // pure bool query
		if q.IsFreeText {
			parent.IsFreeText = true
		}
		return q, nil
	}

	return &GroupTerm{Query: q}, nil
}

func processRange(tok *parsing.Token) (Query, error) {
	if len(tok.Value) < 8 { // [a TO z]
		return nil, ErrRangeValue
	}

	s := string(tok.Value[1 : len(tok.Value)-1])
	tks, err := parsing.Parse(s)
	if err != nil {
		return nil, err
	}

	if len(tks) != 3 || string(tks[1].Value) != "TO" {
		return nil, ErrRangeValue
	}

	if len(tks[0].Prefix) > 0 || len(tks[0].Suffix) > 0 || tks[0].AND || tks[0].OR || tks[0].NOT {
		return nil, ErrRangeValue
	}

	if len(tks[2].Prefix) > 0 || len(tks[2].Suffix) > 0 || tks[2].AND || tks[2].OR || tks[2].NOT {
		return nil, ErrRangeValue
	}

	q := new(RangeTerm)

	switch tks[0].Type {
	case parsing.TEXT:
		v := safeString(tks[0].Value)
		q.Left = NewRangeStringValue(v, tok.Value[0] == '[')
	case parsing.PHRASE:
		v := strconv.Quote(string(tks[0].Value))
		q.Left = NewRangePhraseValue(v, tok.Value[0] == '[')
	default:
		return nil, ErrRangeValue
	}

	switch tks[2].Type {
	case parsing.TEXT:
		v := safeString(tks[2].Value)
		q.Right = NewRangeStringValue(v, tok.Value[len(tok.Value)-1] == ']')
	case parsing.PHRASE:
		v := strconv.Quote(string(tks[2].Value))
		q.Right = NewRangePhraseValue(v, tok.Value[len(tok.Value)-1] == ']')
	default:
		return nil, ErrRangeValue
	}

	return q, nil
}

type textParser struct {
	Value     []byte
	Plus      bool
	Minus     bool
	Field     string
	Op        string
	Boost     float64
	Fuzziness int
	Proximity int
}

// strict = true for pure prefix (+title:, +, title:)
func (p *textParser) ProcessPrefix(buf []byte, strict bool) error {
	sz := len(buf)
	if sz == 0 {
		return nil
	}

	l := 0

	for i, c := range buf {
		var exit bool

		switch c {
		case '+':
			if i == 0 {
				p.Plus = true
				if strict && sz == 1 {
					return nil
				}
				l = 1
			}
		case '-':
			if i == 0 {
				p.Minus = true
				if strict && sz == 1 {
					return nil
				}
				l = 1
			}
		case ':':
			if l == i {
				return ErrBadPrefix
			}

			p.Field = string(buf[l:i])
			l = i + 1

			if l < sz {
				switch buf[l] {
				case '>':
					p.Op = ">"
					l++
				case '<':
					p.Op = "<"
					l++
				default:
					if strict {
						return ErrBadPrefix
					}
				}
			}

			if l < sz {
				switch buf[l] {
				case '=':
					p.Op += "="
					l++
				default:
					if strict {
						return ErrBadPrefix
					}
				}
			}

			exit = true
		}

		if exit {
			break
		}
	}

	if l == sz {
		if strict {
			return nil
		}
		return ErrBadPrefix
	}

	if strict {
		return ErrBadPrefix
	}
	buf = buf[l:]

	for i, c := range buf {
		if c == '^' || c == '~' {
			p.Value = buf[:i]
			return p.ParseSuffix(buf[i:])
		}
	}

	p.Value = buf
	return nil
}

func (p *textParser) ParseSuffix(buf []byte) error {
	sz := len(buf)
	if sz == 0 {
		return nil
	}

	if sz == 1 {
		if buf[0] == '~' {
			p.Fuzziness = 2 // quikc~ 2 is default fuzziness
			return nil
		}
		return ErrBadSuffix // anything of size 1 but ~ is bad
	}

	var (
		fuzzy string // ~ text
		boost string // ^ text
	)

	switch buf[0] {
	case '~':
		p.Fuzziness = 2 // it may change if there are digits after ~

		var r int
		for i := 1; i < len(buf); i++ {
			if buf[i] == '^' {
				r = i
				break
			}
		}

		if r == 0 { // no boost
			fuzzy = string(buf[1:])
		} else if r == len(buf)-1 { // quikc~^ // no ^ value
			return ErrBadSuffix
		} else {
			fuzzy = string(buf[1:r])
			boost = string(buf[r+1:])
		}
	case '^':
		boost = string(buf[1:]) // boost cannot precede ~
	default:
		return ErrBadSuffix // no ^ or ~
	}

	if fuzzy != "" {
		i, err := strconv.Atoi(fuzzy)
		if err != nil {
			return err
		}
		p.Fuzziness = i
		p.Proximity = i
	}

	if boost != "" {
		f, err := strconv.ParseFloat(boost, 64)
		if err != nil {
			return err
		}
		p.Boost = f
	}

	return nil
}

func safeString(buf []byte) string {
	// check for unescaped special characters
	// according to Elastic documentaion https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html#query-string-syntax:
	// The reserved characters are: + - = && || > < ! ( ) { } [ ] ^ " ~ * ? : \ /
	// Failing to escape these special characters correctly could lead to a syntax error which prevents your query from running.
	// < and > can’t be escaped at all. The only way to prevent them from attempting to create a range query is to remove them from the query string entirely.

	// "could" doesn't sound 100%, so I checked all
	// 100% failed characters are ( ) [ ] { }
	// : when first or last in the last term; this is too confusing, so I promoting : to 100% reserved
	// - and + only when first character
	// ! only when first or last character
	// ^ ~ will be converted to suffix
	// * will try to convert to phrase
	// / will try to convert to regex
	// * ? are wildcard characters
	// \\ is escape character
	// our parser will try to convert to other term types on ( [ { " ^ ~

	// let's escape first + and - and all = & | > < ! ) } ] : /

	var sb strings.Builder
	for i, c := range buf {
		switch c {
		case '-', '+':
			if i == 0 {
				sb.WriteByte('\\')
			}
		case ')', ']', '}', '=', '>', '<', '!', ':', '&', '|':
			if i == 0 || buf[i-1] != '\\' {
				sb.WriteByte('\\')
			}
		}
		sb.WriteByte(c)
	}

	return sb.String()
}
