package scanner

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var errUnterminatedString = errors.New("unterminated string.")

const (
	// single character token
	LEFT_PAREN = iota << 1
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	STAR

	// one or two characters token
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL
	SLASH

	// identifiers
	IDENTIFIER
	STRING
	NUMBER

	// keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

var keywords = map[string]int{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"fun":    FUN,
	"for":    FOR,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type token struct {
	tokenType int
	lexeme    string
	literal   interface{}
	line      int
}

func (t token) String() string {
	return fmt.Sprintf("{type: %d, lexeme: %s, literal: %v}", t.tokenType,
		t.lexeme, t.literal)
}

type scanner struct {
	source string
	line   int
	start  int
	cursor int
	tokens []token
}

func NewScanner(src string) (*scanner, error) {
	path, err := filepath.Abs(src)

	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	source, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return &scanner{
		source: string(source),
		line:   0,
		start:  0,
		cursor: 0,
	}, nil
}

func (sc *scanner) ScanToken() []token {
	for {
		sc.start = sc.cursor
		c, err := sc.advance()

		if errors.Is(err, io.EOF) {
			sc.addToken(EOF, "")
			break
		}

		switch c {
		case ' ', '\t', '\r':
			continue
		case '\n':
			sc.line++
		case '{':
			sc.addToken(LEFT_BRACE, "")
		case '}':
			sc.addToken(RIGHT_BRACE, "")
		case ')':
			sc.addToken(RIGHT_PAREN, "")
		case '(':
			sc.addToken(LEFT_PAREN, "")
		case '*':
			sc.addToken(STAR, "")
		case '+':
			sc.addToken(PLUS, "")
		case '-':
			sc.addToken(MINUS, "")
		case ',':
			sc.addToken(COMMA, "")
		case '.':
			sc.addToken(DOT, "")
		case ';':
			sc.addToken(SEMICOLON, "")
		case '!':
			if sc.matchAdvance('=') {
				sc.addToken(BANG_EQUAL, "")
				continue
			}
			sc.addToken(BANG, "")
		case '=':
			if sc.matchAdvance('=') {
				sc.addToken(EQUAL_EQUAL, "")
				continue
			}
			sc.addToken(EQUAL, "")
		case '>':
			if sc.matchAdvance('=') {
				sc.addToken(GREATER_EQUAL, "")
				continue
			}
			sc.addToken(GREATER, "")
		case '<':
			if sc.matchAdvance('=') {
				sc.addToken(LESS_EQUAL, "")
				continue
			}
			sc.addToken(LESS, "")
		case '/':
			if !sc.matchAdvance('/') {
				sc.addToken(SLASH, "")
				continue
			}
			sc.comment()
		case '"':
			sc.consumeString()
		default:
			if isAlpha(c) {
                sc.identifier()
			}
		}
	}
	return sc.tokens
}

func (sc *scanner) identifier() error {
	for {
		c, err := sc.advance()
		if err != nil {
			return err
		}
		if !isAlphaNum(c) {
			break
		}
	}
	text := sc.source[sc.start : sc.cursor-1]
	tokenType := IDENTIFIER
	if t := keywords[text]; t != 0 {
		tokenType = t
	}
	sc.addToken(tokenType, text)
	return nil
}

func isAlpha(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'Z' && c <= 'Z' || c == '_'
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlphaNum(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func (sc *scanner) consumeString() error {
	for {
		c, err := sc.advance()
		if err != nil {
			return errUnterminatedString
		}
		if c == '"' {
			text := sc.source[sc.start : sc.cursor-1]
			sc.addToken(STRING, text)
			return nil
		}
	}
}

func (sc *scanner) comment() error {
	for {
		c, err := sc.advance()
		if err != nil {
			return err
		}
		if c != '\n' {
			break
		}
	}
	return nil
}

func (sc *scanner) matchAdvance(c rune) bool {
	if sc.isAtEnd() || sc.peek() != c {
		return false
	}
	sc.cursor++
	return true
}

func (sc scanner) peek() rune {
	if sc.isAtEnd() {
		return rune(0)
	}
	return rune(sc.source[sc.cursor])
}

func (sc scanner) isAtEnd() bool {
	return sc.cursor >= len(sc.source)
}

func (sc *scanner) addToken(tokenType int, literal string) {
	text := sc.source[sc.start:sc.cursor]
	sc.tokens = append(sc.tokens, token{
		tokenType: tokenType,
		lexeme:    text,
		literal:   literal,
		line:      sc.line,
	})
}

func (sc *scanner) advance() (rune, error) {
	if sc.isAtEnd() {
		return rune(0), io.EOF
	}
	r := rune(sc.source[sc.cursor])
	sc.cursor++
	return r, nil
}
