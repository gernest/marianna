package magic

import (
	"fmt"
	"io"
)

//TokenKind defines type of the token
type TokenKind int

// common token kinds
const (
	Break TokenKind = iota
	ATXHeading
	SetextHeading
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
	Whitespace
)

func (k TokenKind) String() string {
	switch k {
	case Break:
		return "Break"
	case ATXHeading:
		return "ATXHeading"
	case SetextHeading:
		return "SetextHeading"
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
	case Whitespace:
		return "Whitespace"
	}
	return "unkown"
}

//Token identifies a chunk of tex or a character with important meaning
type Token struct {
	Kind  TokenKind
	Text  []byte
	Begin int
	End   int
}

func (t *Token) String() string {
	return fmt.Sprintf(" %s [ %d : %d]", t.Kind, t.Begin, t.End)
}

//LexFunc a function interface for lexing text inputs
type LexFunc func(data []byte, currPos int) (advanceAt int, tok *Token, err error)

//Lexer text tokenizer
type Lexer struct {
	IsBlock func(TokenKind) bool
	LFunc   LexFunc
}

//Lex returns a slice of tokens recognized from src
func (l *Lexer) Lex(src []byte) ([]*Token, error) {
	var (
		currPos = 0
		tokens  []*Token
		lerr    error
	)
STOP:
	for {
		if currPos > len(src)-1 {
			lerr = io.EOF
			break STOP
		}
		a, t, err := l.LFunc(src, currPos)
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
