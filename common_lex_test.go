package magic

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestCommon_LexATXHeading(t *testing.T) {
	sample :=
		[]struct {
			src string
		}{
			{`# foo`},
			{`## foo`},
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
		tkns, err := l.Lex([]byte(v.src))
		if err != nil {
			t.Error(err)
		}
		if tkns[0].Text != v.src {
			t.Errorf(" expected %s fot %s", v.src, tkns[0].Text)
		}
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
	fmt.Println(tk)
}
