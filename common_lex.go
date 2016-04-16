package magic

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type Common struct {
}

func (c *Common) IsBlock(k TokenKind) bool {
	return k == Quote || k == ListItem
}

func (c *Common) Lex(data []byte, currPos int, atEOF bool) (int, *Token, error) {
	ch, _ := utf8.DecodeRune(data)
	switch ch {
	case '#':
		return c.LexATXHeading(data, currPos, atEOF)
	}
	return len(data), nil, nil
}

func (c *Common) LexParagraph(data []byte, currPos int, atEOF bool) (int, *Token, error) {
	ch, _ := utf8.DecodeRune(data)
	if IsLiteral(ch) {
	}
	return len(data), nil, nil
}

func (c *Common) LexATXHeading(data []byte, currPos int, atEOF bool) (int, *Token, error) {
	start := currPos
	end := currPos
	var head []rune
	var body []rune
	isHead := true
	var combineHeadWithBody = func(h []rune, b []rune) string {
		rst := ""
		for _, v := range h {
			rst += string(v)
		}
		for _, v := range b {
			rst += string(v)
		}
		return rst
	}
STOP:
	for {
	repeat:
		if end > len(data)-1 {
			break STOP
		}
		ch, size := utf8.DecodeRune(data[end:])
		if isHead {
			switch ch {
			case '#':
				head = append(head, ch)
				end += size
				goto repeat
			case ' ':
				isHead = false
				end += size
				body = append(body, ch)
				goto repeat
			default:
				if len(head) == 1 && ch == '\n' {
					break STOP
				}
				return len(data), nil, fmt.Errorf(" at %d wrong characher for heading", end)
			}
		}
		switch ch {
		case '\n', '\r':
			end += size
			break STOP
		default:
			end += size
			body = append(body, ch)
		}
	}
	t := &Token{Kind: ATXHeading, Begin: start, End: end}
	if len(head) > 6 {
		// This is a paragraph instead
		t.Kind = Paragraph
		t.Text = combineHeadWithBody(head, body)
		return end, t, nil
	}
	t.Text = combineHeadWithBody(head, body)
	return end, t, nil
}

func IsLiteral(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch)
}
