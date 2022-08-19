package parsing

type byteReader struct {
	buf []byte
	i   int
	b   byte
	r   int
}

func newReader(data []byte) *byteReader {
	if len(data) == 0 {
		return nil
	}

	return &byteReader{
		buf: data,
		r:   len(data) - 1,
	}
}

func (rd *byteReader) Read() bool {
	if rd.i > rd.r {
		return false
	}

	rd.b = rd.buf[rd.i]
	rd.i++ // move to next index
	return true
}

func (rd *byteReader) IsWS() bool {
	wc := rd.b == ' ' || rd.b == '\n' || rd.b == '\r' || rd.b == '\t'
	if wc {
		return !rd.IsEscaped()
	}
	return false
}

func (rd *byteReader) IsEscaped() bool {
	return rd.i > 1 && rd.buf[rd.i-2] == '\\'
}

func (rd *byteReader) Slice(r ...byte) ([]byte, bool) {
	var (
		q   bool // quote
		i   = rd.i - 1
		two = len(r) == 2
	)

	for rd.Read() {
		if rd.b == '"' {
			q = !q
			continue
		}

		ok := rd.b == r[0] && !q
		if !ok && two {
			ok = rd.b == r[1]
		}

		if ok {
			ok = !rd.IsEscaped()
		}

		if ok {
			return rd.buf[i:rd.i], true
		}
	}

	return nil, false
}

func (rd *byteReader) SliceParentheses() ([]byte, bool) {
	var (
		q bool // quote
		i = rd.i
		p = 1 // open parentheses
	)

	for rd.Read() {
		switch rd.b {
		case ')':
			if rd.IsEscaped() || q {
				continue
			}
			p--
			if p == 0 {
				return rd.buf[i : rd.i-1], true
			}
		case '(':
			if !rd.IsEscaped() && !q {
				p++
			}
		case '"':
			if !rd.IsEscaped() {
				q = !q
			}
		}
	}

	return nil, false
}

func (rd *byteReader) SlicePhrase() ([]byte, bool) {
	var (
		i = rd.i
	)

	for rd.Read() {
		if rd.b == '"' && !rd.IsEscaped() {
			return rd.buf[i : rd.i-1], true
		}
	}

	return nil, false
}
