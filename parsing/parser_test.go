package parsing

import "testing"

func TestParser(t *testing.T) {
	q := "\"a"
	if _, err := Parse(q); err == nil {
		t.Fail()
	}

	q = "john OR"
	if _, err := Parse(q); err == nil {
		t.Fail()
	}

	q = "name:\"john smith\"[1 TO 9]^2"
	if _, err := Parse(q); err == nil {
		t.Fail()
	}

	q = "AND john"
	if _, err := Parse(q); err == nil {
		t.Fail()
	}

	q = "OR john"
	if _, err := Parse(q); err == nil {
		t.Fail()
	}

	q = "+name:john^2"
	if _, err := Parse(q); err != nil {
		t.Fail()
		t.Error(err)
	}

	q = "+name:\"john smith\"^2"
	if _, err := Parse(q); err != nil {
		t.Fail()
		t.Error(err)
	}

	q = "\"john smith\"\"yur met\"^2"
	if _, err := Parse(q); err == nil {
		t.Fail()
	}

	q = "\"john smith\"name\"yur met\""
	if _, err := Parse(q); err == nil {
		t.Fail()
	}

	q = "name:\"john smith\"\"yur met\""
	if _, err := Parse(q); err == nil {
		t.Fail()
	}

	q = "+name:\"john smith\""
	if _, err := Parse(q); err != nil {
		t.Fail()
		t.Error(err)
	}

	q = "\"john smith\"^2"
	if _, err := Parse(q); err != nil {
		t.Fail()
		t.Error(err)
	}

	q = "\"john smith\"\"yur met\""
	if _, err := Parse(q); err == nil {
		t.Fail()
	}
}
