package magic

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"
)

//Common implements lexer for commonmark
type Common struct {
}

//IsBlock returns true if the token kind is Block element and false otherwise
func (c *Common) IsBlock(k TokenKind) bool {
	return k == Quote || k == ListItem
}

//Lex implements LexFunc for common mark
func (c *Common) Lex(data []byte, currPos int) (int, *Token, error) {
	ch, _ := utf8.DecodeRune(data[currPos:])
	switch ch {
	case '#':
		return c.LexATXHeading(data, currPos)
	case '\r', '\n':
		return c.LexBlankline(data, currPos)
	case ' ':
		return c.LexWHitespace(data, currPos)
	}

	return c.LexParagraph(data, currPos)
}

//LexParagraph lexes commonmark paragraph
func (c *Common) LexParagraph(data []byte, currPos int) (int, *Token, error) {
	end := currPos
	txt := &bytes.Buffer{}
STOP:
	for {
		if end > len(data)-1 {
			break STOP
		}
		ch, size := utf8.DecodeRune(data[end:])
		switch ch {
		case '\n', '\r':
			txt.WriteRune(ch)
			end += size
			if end > len(data)-1 {
				break STOP
			}
			nch, nsize := utf8.DecodeRune(data[end:])
			switch nch {
			case '\r', '\n':
				break STOP
			default:
				txt.WriteRune(ch)
				end += nsize
			}

		default:
			txt.WriteRune(ch)
			end += size
		}
	}
	t :=
		&Token{Kind: Paragraph, Begin: currPos, End: end, Text: txt.Bytes()}
	return end, t, nil
}

//LexATXHeading lexes commonmark ATXHeading
func (c *Common) LexATXHeading(data []byte, currPos int) (int, *Token, error) {
	end := currPos
	txt := &bytes.Buffer{}
	ch, size := utf8.DecodeRune(data[end:])
	if ch == '#' {

		end += size
		txt.WriteRune(ch)
		next, nsize := utf8.DecodeRune(data[end:])
		switch next {
		case ' ':
			end += nsize
			txt.WriteRune(next)
			goto STOP
		case '#':
			match := 2
		HSTOP:
			for {
				hch, hsize := utf8.DecodeRune(data[end:])
				switch hch {
				case ' ':
					txt.WriteRune(hch)
					end += hsize
				case '#':
					txt.WriteRune(hch)
					end += hsize
					match++
				default:
					break HSTOP
				}
			}
			if match > 6 {
				return c.LexParagraph(data, currPos)
			}
		}
	STOP:
		for {
			if end > len(data)-1 {
				break
			}
			tch, tsize := utf8.DecodeRune(data[end:])
			switch tch {
			case '\r', '\n':
				end += tsize
				txt.WriteRune(tch)
				break STOP
			default:
				end += tsize
				txt.WriteRune(tch)
			}
		}
		t :=
			&Token{Kind: ATXHeading, Begin: currPos, End: end, Text: txt.Bytes()}
		return end, t, nil
	}
	//fmt.Println("HERE", string(data[currPos:]))
	return c.LexParagraph(data, currPos)
}

//IsLiteral checks rune ch if it is a commonmark literal
func IsLiteral(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch)
}

//LexBlankline lexes blank line
func (c *Common) LexBlankline(data []byte, currPos int) (int, *Token, error) {
	end := currPos
	if currPos > len(data)-1 {
		return len(data), nil, nil
	}
	ch, size := utf8.DecodeRune(data[end:])
	if ch == '\r' || ch == '\n' {
		end += size
		t :=
			&Token{Kind: Blankline, Begin: currPos, End: end, Text: []byte(string(ch))}
		return end, t, nil
	}
	return len(data), nil, fmt.Errorf(" at %d txt: %s  failed to lex blankline", currPos, string(ch))
}

//LexWHitespace lexes begin of the the line spaces. More than four white spaces
//signifies an indented code block
func (c *Common) LexWHitespace(data []byte, currPos int) (int, *Token, error) {
	if currPos > len(data)-1 {
		return len(data), nil, nil
	}
	end := currPos
	var chars []rune
	for {
		ch, size := utf8.DecodeRune(data[end:])
		if ch == ' ' {
			chars = append(chars, ch)
			end += size
			if len(chars) == 4 {
				break
			}
			continue
		}
		break
	}
	if len(chars) > 3 {
		// whatever foolows is a indented code blocko
		return c.LexIndentCode(data, end)
	}
	t :=
		&Token{Kind: Whitespace, Begin: currPos, End: end}
	for _, v := range chars {
		t.Text = append(t.Text, []byte(string(v))...)
	}
	return end, t, nil
}

func (c *Common) LexIndentCode(data []byte, currPos int) (int, *Token, error) {
	end := currPos
	txt := &bytes.Buffer{}
STOP:
	for {
		if end > len(data)-1 {
			break STOP
		}
		ch, size := utf8.DecodeRune(data[end:])
		switch ch {
		case '\n', '\r':
			txt.WriteRune(ch)
			end += size
			break STOP
		default:
			txt.WriteRune(ch)
			end += size
		}
	}
	t :=
		&Token{Kind: IndentCode, Begin: currPos, End: end, Text: txt.Bytes()}
	return end, t, nil
}
