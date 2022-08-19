package parsing

import "testing"

func TestBoolean(t *testing.T) {
	toks, err := booleanize(nil)
	if err != nil || toks != nil {
		t.Fail()
	}

	tks, _ := tokenize("NOT john AND smith AND (text OR video) AND \"yur met\" AND NOT age:[50 TO 80} AND name:/joh?n(ath[oa]n)/")
	toks, err = booleanize(tks)
	if err != nil || len(toks) != 6 {
		t.Fail()
	}
	if !toks[0].IsNot {
		t.Fail()
	}
	if toks[2].Op != "AND" {
		t.Fail()
	}
	if !toks[4].IsNot || toks[4].Op != "AND" {
		t.Fail()
	}

	tks, _ = tokenize("john AND AND smith")
	if _, err := booleanize(tks); err == nil {
		t.Fail()
	}

	tks, _ = tokenize("john OR AND smith")
	if _, err := booleanize(tks); err == nil {
		t.Fail()
	}

	tks, _ = tokenize("john OR OR smith")
	if _, err := booleanize(tks); err == nil {
		t.Fail()
	}

	tks, _ = tokenize("john AND OR smith")
	if _, err := booleanize(tks); err == nil {
		t.Fail()
	}

	tks, _ = tokenize("john AND NOT NOT smith")
	if _, err := booleanize(tks); err == nil {
		t.Fail()
	}

	tks, _ = tokenize("john AND")
	if _, err := booleanize(tks); err == nil {
		t.Fail()
	}

	tks, _ = tokenize("(john smith)AND text")
	if _, err := booleanize(tks); err != nil {
		t.Fail()
	}

	tks, _ = tokenize("(john smith)OR text")
	if _, err := booleanize(tks); err != nil {
		t.Fail()
	}

	tks, _ = tokenize("(john smith)NOT text")
	if _, err := booleanize(tks); err != nil {
		t.Fail()
	}

}
