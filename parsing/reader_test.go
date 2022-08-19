package parsing

import "testing"

func TestReader(t *testing.T) {
	r := newReader(nil)
	if r != nil {
		t.Fail()
	}

	r = newReader([]byte("x \\y"))
	if r == nil {
		t.Fail()
	}
	if !r.Read() || r.b != 'x' || r.i != 1 {
		t.Fail()
	}
	if !r.Read() || !r.IsWS() || r.i != 2 {
		t.Fail()
	}
	if !r.Read() || r.b != '\\' || r.i != 3 {
		t.Fail()
	}
	if !r.Read() || r.b != 'y' || r.i != 4 || !r.IsEscaped() {
		t.Fail()
	}
	if r.Read() {
		t.Fail()
	}

	s := "[\"[a]\" TO \\] xyz]"
	r = newReader([]byte(s))
	r.Read()
	bs, ok := r.Slice(']')
	if !ok || string(bs) != s {
		t.Fail()
	}

	s = "[a TO z}"
	r = newReader([]byte(s))
	r.Read()
	bs, ok = r.Slice(']', '}')
	if !ok || string(bs) != s {
		t.Fail()
	}

	s = "[a TO z}"
	r = newReader([]byte(s))
	r.Read()
	if _, ok := r.Slice(']'); ok {
		t.Fail()
	}

	s = "(a OR (\"()\") OR \\) AND NOT x)"
	r = newReader([]byte(s))
	r.Read()
	bs, ok = r.SliceParentheses()
	if !ok || string(bs) != s[1:len(s)-1] {
		t.Fail()
	}

	s = "(a OR (b OR c)"
	r = newReader([]byte(s))
	r.Read()
	if _, ok := r.SliceParentheses(); ok {
		t.Fail()
	}

	s = "\"joseph \\\"john\\\" smith\""
	r = newReader([]byte(s))
	r.Read()
	bs, ok = r.SlicePhrase()
	if !ok || string(bs) != s[1:len(s)-1] {
		t.Fail()
	}

	s = "\"john smith"
	r = newReader([]byte(s))
	r.Read()
	if _, ok := r.SlicePhrase(); ok {
		t.Fail()
	}
}
