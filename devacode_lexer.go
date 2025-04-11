// devacode_lexer.go
package devacode

import (
    "unicode"
    "unicode/utf8"
    "strings"
    "fmt"
)

type TokenType string

type Token struct {
    Type    TokenType
    Literal string
}

const (
    IDENTIFIER TokenType = "IDENTIFIER"
    NUMBER     TokenType = "NUMBER"
    STRING     TokenType = "STRING"
    KEYWORD    TokenType = "KEYWORD"
    OPERATOR   TokenType = "OPERATOR"
    DANDA      TokenType = "DANDA"
    EOF        TokenType = "EOF"
)

var keywords = map[string]TokenType{
    "यदि": KEYWORD,
    "अन्यथा": KEYWORD,
    "विधि": KEYWORD,
    "लेख": KEYWORD,
    "यावत्": KEYWORD,
}

func isDevanagari(r rune) bool {
    return unicode.In(r, unicode.Devanagari)
}

type Lexer struct {
    input string
    pos   int
}

func NewLexer(input string) *Lexer {
    return &Lexer{input: input}
}

func (l *Lexer) NextToken() Token {
    for l.pos < len(l.input) {
        r, size := utf8.DecodeRuneInString(l.input[l.pos:])

        if unicode.IsSpace(r) {
            l.pos += size
            continue
        }

        if r == '।' {
            l.pos += size
            return Token{Type: DANDA, Literal: "।"}
        }

        if isDevanagari(r) {
            start := l.pos
            for l.pos < len(l.input) {
                r2, size2 := utf8.DecodeRuneInString(l.input[l.pos:])
                if !isDevanagari(r2) && !unicode.IsDigit(r2) {
                    break
                }
                l.pos += size2
            }
            word := l.input[start:l.pos]
            if tokType, ok := keywords[word]; ok {
                return Token{Type: tokType, Literal: word}
            }
            return Token{Type: IDENTIFIER, Literal: word}
        }

        if unicode.IsDigit(r) {
            start := l.pos
            for l.pos < len(l.input) {
                r2, size2 := utf8.DecodeRuneInString(l.input[l.pos:])
                if !unicode.IsDigit(r2) {
                    break
                }
                l.pos += size2
            }
            return Token{Type: NUMBER, Literal: l.input[start:l.pos]}
        }

        l.pos += size
        return Token{Type: OPERATOR, Literal: string(r)}
    }
    return Token{Type: EOF, Literal: ""}
}

// Simple test function to show tokens
func Tokenize(input string) {
    lexer := NewLexer(input)
    for tok := lexer.NextToken(); tok.Type != EOF; tok = lexer.NextToken() {
        fmt.Printf("Token: Type=%s, Literal=%s\n", tok.Type, tok.Literal)
    }
}
