package magic

import "testing"

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
