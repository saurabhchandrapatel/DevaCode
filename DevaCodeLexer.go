package devacode

import (
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"
)

// Token types
type TokenType string

const (
	TokenKeyword   TokenType = "KEYWORD"
	TokenIdentifier TokenType = "IDENTIFIER"
	TokenNumber     TokenType = "NUMBER"
	TokenOperator   TokenType = "OPERATOR"
	TokenString     TokenType = "STRING"
	TokenEOF        TokenType = "EOF"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Line    int
	Column  int
}

type Lexer struct {
	source  string
	pos     int
	line    int
	column  int
	tokens  []Token
}

// Keywords in Devanagari
var keywords = map[string]TokenType{
	"यदि":     TokenKeyword,
	"अन्यथा":  TokenKeyword,
	"पर्यन्तम्": TokenKeyword,
	"यावत्":   TokenKeyword,
	"लेख":     TokenKeyword,
	"विधि":     TokenKeyword,
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		source: input,
		line:   1,
		column: 0,
	}
}

func (l *Lexer) isDevanagari(r rune) bool {
	return unicode.In(r, unicode.Devanagari)
}

func (l *Lexer) NextToken() Token {
	if l.pos >= len(l.source) {
		return Token{Type: TokenEOF, Lexeme: "", Line: l.line, Column: l.column}
	}

	r, size := utf8.DecodeRuneInString(l.source[l.pos:])
	l.pos += size
	l.column++

	if unicode.IsSpace(r) {
		if r == '\n' {
			l.line++
			l.column = 0
		}
		return l.NextToken()
	}

	start := l.pos - size

	if l.isDevanagari(r) {
		for l.pos < len(l.source) {
			r2, s2 := utf8.DecodeRuneInString(l.source[l.pos:])
			if !l.isDevanagari(r2) {
				break
			}
			l.pos += s2
			l.column++
		}
		word := l.source[start:l.pos]
		if tokenType, found := keywords[word]; found {
			return Token{Type: tokenType, Lexeme: word, Line: l.line, Column: l.column}
		}
		return Token{Type: TokenIdentifier, Lexeme: word, Line: l.line, Column: l.column}
	}

	if unicode.IsDigit(r) {
		for l.pos < len(l.source) {
			r2, s2 := utf8.DecodeRuneInString(l.source[l.pos:])
			if !unicode.IsDigit(r2) {
				break
			}
			l.pos += s2
			l.column++
		}
		return Token{Type: TokenNumber, Lexeme: l.source[start:l.pos], Line: l.line, Column: l.column}
	}

	// Match operators like `=`, `+`, etc.
	opRegex := regexp.MustCompile(`[=+*/<>!]+`)
	if loc := opRegex.FindStringIndex(l.source[start:]); loc != nil && loc[0] == 0 {
		op := l.source[start : start+loc[1]]
		l.pos = start + loc[1]
		l.column += loc[1] - 1
		return Token{Type: TokenOperator, Lexeme: op, Line: l.line, Column: l.column}
	}

	return Token{Type: TokenIdentifier, Lexeme: string(r), Line: l.line, Column: l.column}
}

func (l *Lexer) Tokenize() []Token {
	for {
		token := l.NextToken()
		l.tokens = append(l.tokens, token)
		if token.Type == TokenEOF {
			break
		}
	}
	return l.tokens
}
