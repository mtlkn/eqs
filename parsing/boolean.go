package parsing

import (
	"errors"
)

type boolToken struct {
	Values []*byteToken
	Op     string // operator AND or OR
	IsNot  bool   // has NOT prefix
}

func booleanize(toks []*byteToken) ([]*boolToken, error) {
	if len(toks) == 0 {
		return nil, nil
	}

	var (
		tokens []*boolToken
		token  *boolToken // current token
	)

	for _, tok := range toks {
		if tok.Type == WHITE_SPACE {
			if token != nil && len(token.Values) > 0 {
				tokens = append(tokens, token)
				token = nil
			}
			continue
		}

		if tok.Type == TEXT {
			switch string(tok.Value) {
			case "AND":
				if token != nil {
					if len(token.Values) > 0 {
						tokens = append(tokens, token)
					} else {
						return nil, errors.New("booleanizer: bad AND operator")
					}
				}

				token = &boolToken{
					Op: "AND",
				}
			case "OR":
				if token != nil {
					if len(token.Values) > 0 {
						tokens = append(tokens, token)
					} else {
						return nil, errors.New("booleanizer: bad OR operator")
					}
				}

				token = &boolToken{
					Op: "OR",
				}
			case "NOT":
				if token != nil {
					if token.IsNot {
						return nil, errors.New("booleanizer: bad NOT operator")
					}

					if len(token.Values) > 0 {
						tokens = append(tokens, token)
					} else {
						token.IsNot = true
						continue
					}
				}

				token = &boolToken{
					IsNot: true,
				}
			default:
				if token == nil {
					token = new(boolToken)
				}
				token.Values = append(token.Values, tok)
			}

			continue
		}

		if tok.Type == PHRASE || tok.Type == GROUP || tok.Type == RANGE || tok.Type == REGEX {
			if token == nil {
				token = new(boolToken)
			}

			token.Values = append(token.Values, tok)
		}
	}

	if token != nil {
		if len(token.Values) == 0 {
			return nil, errors.New("booleanizer: missing token value at closing")
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}
