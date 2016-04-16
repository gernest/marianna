package magic

import (
	"fmt"
	"io"
)

type TokenKind int

const (
	Break TokenKind = iota
	ATXHeading
	SelectHeading
	IndentCode
	FencedCode
	HTML
	LinkReference
	Paragraph
	Blankline
	Quote
	List
	ListItem
	Escape
	EntityReference
	CodeSpan
	Emphasis
	StrongEmphasis
	Link
	Image
	HardLineBreak
	SoftLineBreak
	Text
)

func (k TokenKind) String() string {
	switch k {
	case Break:
		return "Break"
	case ATXHeading:
		return "ATXHeading"
	case SelectHeading:
		return "SelectHeading"
	case IndentCode:
		return "IndentCode"
	case FencedCode:
		return "FencedCode"
	case HTML:
		return "HTML"
	case LinkReference:
		return "LinkReference"
	case Paragraph:
		return "Paragraph"
	case Blankline:
		return "Blankline"
	case Quote:
		return "Quote"
	case List:
		return "List"
	case ListItem:
		return "ListItem"
	case Escape:
		return "Escape"
	case EntityReference:
		return "EntityReference"
	case CodeSpan:
		return "CodeSpan"
	case Emphasis:
		return "Emphasis"
	case StrongEmphasis:
		return "StrongEmphasis"
	case Link:
		return "Link"
	case Image:
		return "Image"
	case HardLineBreak:
		return "HardLineBreak"
	case SoftLineBreak:
		return "SoftLineBreak"
	case Text:
		return "Text"
	}
	return "unkown"
}

type Token struct {
	Kind  TokenKind
	Text  string
	Begin int
	End   int
}

func (t *Token) String() string {
	return fmt.Sprintf(" %s [ %d : %d]", t.Kind, t.Begin, t.End)
}

type LexFunc func(data []byte, currPos int, atEOF bool) (advanceAt int, tok *Token, err error)

type Lexer struct {
	IsBlock func(TokenKind) bool
	LFunc   LexFunc
}

func (l *Lexer) Lex(src []byte) ([]*Token, error) {
	var (
		currPos = 0
		tokens  []*Token
		atEOF   = false
		lerr    error
	)
STOP:
	for {
		if atEOF {
			lerr = io.EOF
			break STOP
		}
		if currPos > len(src)-1 {
			atEOF = true
		}
		a, t, err := l.LFunc(src[currPos:], currPos, atEOF)
		if err != nil {
			lerr = err
			break STOP
		}
		if t == nil {
			break STOP
		}
		currPos = a
		tokens = append(tokens, t)
	}
	if lerr != nil {
		if lerr.Error() == io.EOF.Error() && tokens != nil {
			return tokens, nil
		}
		return nil, lerr
	}
	return tokens, nil
}
func updatePositions(tokens []*Token, begin int) []*Token {
	return tokens
}
