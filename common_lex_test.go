package magic

import (
	"io/ioutil"
	"testing"
)

func TestCommon_LexATXHeading(t *testing.T) {
	sample :=
		[]struct {
			src string
		}{
			{`# foo`},
			{`# foo`},
			{`### foo`},
			{`#### foo`},
			{`##### foo`},
			{`###### foo`},
			{`####### foo`},
		}
	c := &Common{}
	l := Lexer{}
	l.IsBlock = c.IsBlock
	l.LFunc = c.Lex

	for _, v := range sample {
		_, err := l.Lex([]byte(v.src))
		if err != nil {
			t.Error(err)
		}
		//fmt.Println(tkns)
	}
}

func TestCommon_LexBlankline(t *testing.T) {
	b, err := ioutil.ReadFile("fixture/lex/blankline.md")
	if err != nil {
		t.Fatal(err)
	}
	c := &Common{}
	l := Lexer{}
	l.IsBlock = c.IsBlock
	l.LFunc = c.Lex
	tk, err := l.Lex(b)
	if err != nil {
		t.Error(err)
	}
	if len(tk) != 5 {
		t.Errorf("expected %d got %d", 5, len(tk))
	}
}

func TestCommon_LexWHitespace(t *testing.T) {
	src := "   # bar"
	c := &Common{}
	l := Lexer{}
	l.IsBlock = c.IsBlock
	l.LFunc = c.Lex
	tk, err := l.Lex([]byte(src))
	if err != nil {
		t.Error(err)
	}
	if len(tk) != 2 {
		t.Errorf("expected 2 got %d", len(tk))
	}
}
