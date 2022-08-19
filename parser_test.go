package eqs

import (
	"fmt"
	"os"
	"testing"

	js "github.com/mtlkn/json"
)

var aliases = map[string][]string{
	"head":   {"headline", "title"},
	"person": {"person.name"},
}

func TestBoolean(t *testing.T) {
	testFile("boolean.json", t)
}

func TestRange(t *testing.T) {
	testFile("range.json", t)
}

func TestGroup(t *testing.T) {
	testFile("group.json", t)
}

func TestAliases(t *testing.T) {
	testFile("aliases.json", t)
}

func TestRandom(t *testing.T) {
	testFile("random.json", t)
}

func TestBase(t *testing.T) {
	testFile("base.json", t)
}

func testFile(f string, t *testing.T) {
	bs, err := os.ReadFile("testdata/" + f)
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}

	jo, err := js.ParseObject(bs)
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}

	fmt.Println()
	fmt.Println("Testing " + f)
	fmt.Println("==========================")

	for _, jp := range jo.Properties {
		switch jp.Name {
		case "direct":
			fmt.Println("DIRECT")
			ss, _ := jp.GetStrings()
			for _, s := range ss {
				fmt.Print("\t", s, " ->")
				q, err := Parse(s, aliases)
				if err != nil {
					fmt.Println(" ERROR")
					Parse(s, aliases)
					t.Fail()
					continue
				}

				if q.String() != s {
					fmt.Println(" MISMATCH")
					q, _ = Parse(s, aliases)
					fmt.Println(q.String())
					t.Fail()
					continue
				}

				fmt.Println(" OK")
			}
		case "transform":
			fmt.Println("TRANSFORM")
			ja, _ := jp.GetArray()
			for _, v := range ja.Values {
				a, _ := js.ArrayValue(v)
				ss, _ := a.GetStrings()

				fmt.Print("\t", ss[0], " ->")
				q, err := Parse(ss[0], aliases)
				if err != nil {
					fmt.Println(" ERROR")
					Parse(ss[0], aliases)
					t.Fail()
					continue
				}

				if q.String() != ss[1] {
					fmt.Println(" MISMATCH")
					q, _ = Parse(ss[0], aliases)
					fmt.Println(q.String())
					t.Fail()
					continue
				}

				fmt.Println(" ", q.String())
			}
		case "errors":
			fmt.Println("ERRORS")
			ss, _ := jp.GetStrings()
			for _, s := range ss {
				fmt.Print("\t", s, " ->")
				if _, err := Parse(s, aliases); err == nil {
					fmt.Println(" ERROR")
					Parse(s, aliases)
					t.Fail()
					continue
				}

				fmt.Println(" OK")
			}
		}
	}
}
