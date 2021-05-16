package main

import (
	"fmt"
	"strings"
)

// ----------------------------------------------------------------------------
// CONSTANTS
// ----------------------------------------------------------------------------

const DIGITS = "0123456789"

// ----------------------------------------------------------------------------
// POSITION
// ----------------------------------------------------------------------------

type position struct {
	Idx  int
	Ln   int
	Col  int
	Fn   string
	Ftxt string
}

func NewPosition(idx int, ln int, col int, fn string, ftxt string) position {
	return position{idx, ln, col, fn, ftxt}
}

func (p *position) Advance(currentChar rune) {
	p.Idx++
	p.Col++

	if currentChar == '\n' {
		p.Ln++
		p.Col = 0
	}
}

// ----------------------------------------------------------------------------
// TOKENS
// ----------------------------------------------------------------------------

const TT_INT = "TT_INT"
const TT_FLOAT = "TT_FLOAT"
const TT_PLUS = "TT_PLUS"
const TT_MINUS = "TT_MINUS"
const TT_MUL = "TT_MUL"
const TT_DIV = "TT_DIV"
const TT_LPAREN = "TT_LPAREN"
const TT_RPAREN = "TT_RPAREN"

type token struct {
	Type, Value string
}

func (t token) String() string {
	if t.Value != "" {
		return fmt.Sprintf("%v:%v", t.Type, t.Value)
	} else {
		return t.Type
	}
}

func NewToken(type_ string, value string) token {
	t := token{type_, value}
	return t
}

// ----------------------------------------------------------------------------
// LEXER
// ----------------------------------------------------------------------------

type lexer struct {
	Fn          string
	Text        []rune
	Pos         position
	CurrentChar rune
}

func NewLexer(fn string, text string) lexer {
	l := lexer{
		Fn:          fn,
		Text:        []rune(text),
		Pos:         NewPosition(-1, 0, 1, fn, text),
		CurrentChar: 0,
	}
	l.Advance()
	return l
}

func (l *lexer) Advance() {
	l.Pos.Advance(l.CurrentChar)
	if l.Pos.Idx < len(l.Text) {
		l.CurrentChar = l.Text[l.Pos.Idx]
	} else {
		l.CurrentChar = 0
	}
}

func (l *lexer) MakeTokens() ([]token, error) {
	var tokens []token

	for l.CurrentChar != 0 {
		if strings.Contains(" \t\r\n", string(l.CurrentChar)) {
			// pass
		} else if strings.Contains(DIGITS, string(l.CurrentChar)) {
			tokens = append(tokens, l.MakeNumber())
			// Do not advance lexer
			continue
		} else if l.CurrentChar == '+' {
			tokens = append(tokens, NewToken(TT_PLUS, ""))
		} else if l.CurrentChar == '-' {
			tokens = append(tokens, NewToken(TT_MINUS, ""))
		} else if l.CurrentChar == '*' {
			tokens = append(tokens, NewToken(TT_MUL, ""))
		} else if l.CurrentChar == '/' {
			tokens = append(tokens, NewToken(TT_DIV, ""))
		} else if l.CurrentChar == '(' {
			tokens = append(tokens, NewToken(TT_LPAREN, ""))
		} else if l.CurrentChar == ')' {
			tokens = append(tokens, NewToken(TT_RPAREN, ""))
		} else {
			char := l.CurrentChar
			posStart := l.Pos
			l.Advance()
			return []token{}, fmt.Errorf(
				"illegal character %q in file %q, line %d",
				char,
				posStart.Fn,
				posStart.Ln+1,
			)
		}

		l.Advance()
	}

	return tokens, nil
}

func (l *lexer) MakeNumber() token {
	var numStr strings.Builder
	dotCount := 0

	for ; (l.CurrentChar != 0) && strings.Contains(DIGITS+".", string(l.CurrentChar)); l.Advance() {
		if l.CurrentChar == '.' {
			if dotCount == 1 {
				break
			}
			dotCount += 1
		}
		numStr.WriteRune(l.CurrentChar)
	}

	if dotCount == 0 {
		return NewToken(TT_INT, numStr.String())
	} else {
		return NewToken(TT_FLOAT, numStr.String())
	}
}

// ----------------------------------------------------------------------------
// RUN
// ----------------------------------------------------------------------------

func Run(fn string, text string) ([]token, error) {
	lex := NewLexer(fn, text)
	tokens, err := lex.MakeTokens()

	return tokens, err
}
