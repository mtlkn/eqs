package parsing

import "testing"

func TestTokenizer(t *testing.T) {
	toks, err := tokenize("")
	if err != nil || toks != nil {
		t.Fail()
	}

	toks, err = tokenize(" john smith AND (text OR video) AND \"yur met\" AND age:[50 TO 80} AND name:/joh?n(ath[oa]n)/")
	if err != nil || len(toks) != 22 {
		t.Fail()
	}
	if toks[0].Type != WHITE_SPACE {
		t.Fail()
	}
	if toks[1].Type != TEXT {
		t.Fail()
	}
	if toks[2].Type != WHITE_SPACE {
		t.Fail()
	}
	if toks[3].Type != TEXT {
		t.Fail()
	}
	if toks[4].Type != WHITE_SPACE {
		t.Fail()
	}
	if toks[5].Type != TEXT {
		t.Fail()
	}
	if toks[7].Type != GROUP {
		t.Fail()
	}
	if toks[11].Type != PHRASE {
		t.Fail()
	}
	if toks[15].Type != TEXT {
		t.Fail()
	}
	if toks[16].Type != RANGE {
		t.Fail()
	}
	if toks[21].Type != REGEX {
		t.Fail()
	}

	if _, err := tokenize("\"a"); err == nil {
		t.Fail()
	}

	if _, err := tokenize("(a OR (b OR c)"); err == nil {
		t.Fail()
	}

	if _, err := tokenize("age:[26 TO 50"); err == nil {
		t.Fail()
	}

	if _, err := tokenize("\\\" \""); err == nil {
		t.Fail()
	}

	if _, err := tokenize("\\) )"); err == nil {
		t.Fail()
	}

	if _, err := tokenize("\\] ]"); err == nil {
		t.Fail()
	}

	if _, err := tokenize("\\} }"); err == nil {
		t.Fail()
	}

	if _, err := tokenize("\\/ /"); err == nil {
		t.Fail()
	}

	if _, err := tokenize("\\\"a \\(b \\[c \\{d"); err != nil {
		t.Fail()
		t.Error(err)
	}

	toks, err = tokenize("john\\ smith")
	if err != nil || len(toks) != 1 {
		t.Fail()
	}
}
