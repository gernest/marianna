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

func TestCommon_LexParagraph(t *testing.T) {
	naormal := `aaa

	aaa`
	c := &Common{}
	l := Lexer{}
	l.IsBlock = c.IsBlock
	l.LFunc = c.Lex
	tk, err := l.Lex([]byte(naormal))
	if err != nil {
		t.Error(err)
	}
	if len(tk) != 3 {
		t.Errorf("expected 3 got %d", len(tk))
	}

	multi := `aaa
	bbb

	ccc
	ddd`
	tk, err = l.Lex([]byte(multi))
	if err != nil {
		t.Error(err)
	}
	if len(tk) != 3 {
		t.Errorf("expected 3 got %d", len(tk))
	}
	multiBlank := `aaa


	bbb`
	tk, err = l.Lex([]byte(multiBlank))
	if err != nil {
		t.Error(err)
	}
	if len(tk) != 4 {
		t.Errorf("expected 4 got %d", len(tk))
	}

	leadSpace := `   aaa
	bbb`
	tk, err = l.Lex([]byte(leadSpace))
	if err != nil {
		t.Error(err)
	}
	if len(tk) != 2 {
		t.Errorf("expected 2 got %d", len(tk))
	}
	indent := `aaa
	         bbb
		              			 
ccc`
	tk, err = l.Lex([]byte(indent))
	if err != nil {
		t.Error(err)
	}
	if len(tk) != 1 {
		t.Errorf("expected 1 got %d", len(tk))
	}
}
