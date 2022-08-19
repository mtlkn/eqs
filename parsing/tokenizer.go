package parsing

import (
	"errors"
)

const (
	WHITE_SPACE = iota
	TEXT
	PHRASE
	GROUP
	RANGE
	REGEX
)

type byteToken struct {
	Type  int
	Value []byte
}

func tokenize(q string) ([]*byteToken, error) {
	if q == "" {
		return nil, nil
	}

	var (
		tokens []*byteToken
		l      int // start index of token
		rd     = newReader([]byte(q))
	)

	for rd.Read() {
		if rd.IsWS() {
			tokens = rd.Flush(l, tokens)

			tokens = append(tokens, &byteToken{
				Type: WHITE_SPACE,
			})

			l = rd.i
			continue
		}

		switch rd.b {
		case '"':
			if rd.IsEscaped() {
				continue
			}

			tokens = rd.Flush(l, tokens)

			buf, ok := rd.SlicePhrase()
			if !ok {
				return nil, errors.New("tokenizer: missing closing quote")
			}

			tokens = append(tokens, &byteToken{
				Type:  PHRASE,
				Value: buf,
			})

			l = rd.i
		case '(':
			if rd.IsEscaped() {
				continue
			}

			tokens = rd.Flush(l, tokens)

			buf, ok := rd.SliceParentheses()
			if !ok {
				return nil, errors.New("tokenizer: missing closing parenthesis")
			}

			tokens = append(tokens, &byteToken{
				Type:  GROUP,
				Value: buf,
			})

			l = rd.i
		case '[', '{':
			if rd.IsEscaped() {
				continue
			}

			tokens = rd.Flush(l, tokens)

			buf, ok := rd.Slice(']', '}')
			if !ok {
				return nil, errors.New("tokenizer: missing closing range")
			}

			tokens = append(tokens, &byteToken{
				Type:  RANGE,
				Value: buf,
			})

			l = rd.i
		case '/':
			if rd.IsEscaped() {
				continue
			}

			tokens = rd.Flush(l, tokens)

			buf, ok := rd.Slice('/')
			if !ok {
				return nil, errors.New("tokenizer: missing closing regex")
			}

			tokens = append(tokens, &byteToken{
				Type:  REGEX,
				Value: buf,
			})
			l = rd.i
		case ')', ']', '}':
			if rd.IsEscaped() {
				continue
			}
			// we should not see those if they have openning counterparts
			return nil, errors.New("tokenizer: missing openning for " + string(rd.b))
		}
	}

	if l <= rd.r {
		tokens = append(tokens, &byteToken{
			Type:  TEXT,
			Value: rd.buf[l:],
		})
	}

	return tokens, nil
}

func (rd *byteReader) Flush(l int, tokens []*byteToken) []*byteToken {
	if rd.i-1 > l {
		tokens = append(tokens, &byteToken{
			Type:  TEXT,
			Value: rd.buf[l : rd.i-1],
		})
	}
	return tokens
}
