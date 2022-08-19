package parsing

import "errors"

type Token struct {
	Type   int
	Value  []byte
	Prefix []byte
	Suffix []byte
	AND    bool // AND v
	OR     bool // OR v
	NOT    bool // NOT v
}

func Parse(q string) ([]*Token, error) {
	tks, err := tokenize(q)
	if err != nil {
		return nil, err
	}

	toks, err := booleanize(tks)
	if err != nil {
		return nil, err
	}

	var (
		tokens []*Token
	)

	for i, tok := range toks {
		sz := len(tok.Values)
		if sz > 3 {
			return nil, errors.New("too many tokens")
		}

		token := &Token{
			NOT: tok.IsNot,
		}

		if tok.Op == "AND" {
			token.AND = true
			if i == 0 {
				return nil, errors.New("AND")
			}
		} else if tok.Op == "OR" {
			token.OR = true
			if i == 0 {
				return nil, errors.New("OR")
			}
		}

		if sz == 1 {
			v := tok.Values[0]
			token.Type = v.Type
			token.Value = v.Value
			tokens = append(tokens, token)
			continue
		}

		if sz == 3 {
			v := tok.Values[0]
			if v.Type != TEXT {
				return nil, errors.New("bad prefix")
			}
			token.Prefix = v.Value

			v = tok.Values[1]
			token.Type = v.Type
			token.Value = v.Value

			v = tok.Values[2]
			if v.Type != TEXT {
				return nil, errors.New("bad suffix")
			}
			token.Suffix = v.Value

			tokens = append(tokens, token)
			continue
		}

		l := tok.Values[0]
		r := tok.Values[1]

		if l.Type == TEXT {
			token.Prefix = l.Value
			token.Type = r.Type
			token.Value = r.Value
		} else if r.Type == TEXT {
			token.Suffix = r.Value
			token.Type = l.Type
			token.Value = l.Value
		} else {
			return nil, errors.New("bad two values token")
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}
