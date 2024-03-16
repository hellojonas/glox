package scanner

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type Token struct {
	tokenType int
	lexeme    string
	literal   interface{}
	line      int
}

type Scanner struct {
	source string
	line   int
	start  int
	cursor int
	tokens []Token
	Errors []error
}

func NewScanner(src string) (*Scanner, error) {
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

	return &Scanner{
		source: string(source),
		line:   1,
		start:  0,
		cursor: 0,
		Errors: make([]error, 0),
	}, nil
}

func (t Token) String() string {
	return fmt.Sprintf("{type: %s, lexeme: %s, literal: %v, line: %d}",
		tokenMap[t.tokenType], t.lexeme, t.literal, t.line)
}

func (sc *Scanner) Scan() ([]Token, []error) {
	for {
		sc.start = sc.cursor
		c, err := sc.advance()
		if err != nil {
			sc.addToken(EOF, "")
			break
		}
		switch c {
		case ' ', '\t', '\r':
			continue
		case '\n':
			sc.line++
		case '(':
			sc.addToken(LEFT_PAREN, "")
		case ')':
			sc.addToken(RIGHT_PAREN, "")
		case '{':
			sc.addToken(LEFT_BRACE, "")
		case '}':
			sc.addToken(RIGHT_BRACE, "")
		case ';':
			sc.addToken(SEMICOLON, "")
		case '.':
			sc.addToken(DOT, "")
		case '+':
			sc.addToken(PLUS, "")
		case '-':
			sc.addToken(MINUS, "")
		case '*':
			sc.addToken(STAR, "")
		case '=':
			if sc.matchAdvance('=') {
				sc.addToken(EQUAL_EQUAL, "")
				continue
			}
			sc.addToken(EQUAL, "")
		case '!':
			if sc.matchAdvance('=') {
				sc.addToken(BANG_EQUAL, "")
				continue
			}
			sc.addToken(BANG, "")
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
		case '"':
			sc.consumeString()
		case '/':
			if sc.matchAdvance('/') {
				sc.comment()
				continue
			} else if sc.matchAdvance('*') {
				// TODO: handle block comment
				continue
			}
			sc.addToken(SLASH, "")
		default:
			if isAlpha(c) {
				sc.identifier()
			} else if isDigit(c) {
				sc.digit()
			}
		}
	}
	return sc.tokens, sc.Errors
}

func (sc *Scanner) comment() {
	for {
		c, _ := sc.advance() 
		if c == '\n' {
			sc.line++
			break
		}
	}
}

func (sc *Scanner) digit() error {
	for {
		if sc.isAtEnd() {
			break
		}
		c := sc.peek()
		if isDigit(c) {
			sc.advance()
			continue
		}
		if c == '.' && isDigit(sc.peekNext()) {
			sc.advance()
			for isDigit(sc.peek()) {
				sc.advance()
			}
		}
		break
	}
	text := sc.source[sc.start:sc.cursor]
	n, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return err
	}
	sc.addToken(NUMBER, n)
	return nil
}

func (sc *Scanner) consumeString() error {
	for {
		c, err := sc.advance()
		if err != nil {
			return errUnterminatedString
		}
		if c == '"' {
			break
		}
		if c == '\n' {
			sc.line++
		}
	}
	text := sc.source[sc.start+1:sc.cursor-1]
	sc.addToken(STRING, text)
	return nil
}

func (sc *Scanner) identifier() {
	for {
		c := sc.peek()
		if sc.isAtEnd() || !isAlphaNum(c) {
			break
		}
		sc.advance()
	}
	text := sc.source[sc.start:sc.cursor]
	tokenType := IDENTIFIER
	key := keywords[text]
	if key != 0 {
		tokenType = keywords[text]
		text = ""
	}
	sc.addToken(tokenType, text)
}

func (sc *Scanner) advance() (rune, error) {
	if sc.isAtEnd() {
		return rune(0), io.EOF
	}
	c := rune(sc.source[sc.cursor])
	sc.cursor++
	return c, nil
}

func (sc *Scanner) matchAdvance(c rune) bool {
	if sc.isAtEnd() {
		return false
	}
	if rune(sc.source[sc.cursor]) != c {
		return false
	}
	sc.cursor++
	return true
}

func (sc *Scanner) addToken(tokenType int, literal interface{}) {
	lexeme := sc.source[sc.start:sc.cursor]
	sc.tokens = append(sc.tokens, Token{
		tokenType: tokenType,
		lexeme: lexeme,
		literal: literal,
		line: sc.line,
	})
}

func (sc *Scanner) addError(err error) {
	sc.Errors = append(sc.Errors, err)
}

func (sc Scanner) hasError() bool {
	return len(sc.Errors) > 0
}

func (sc Scanner) peek() rune {
	if sc.isAtEnd() {
		return rune(0)
	}
	return rune(sc.source[sc.cursor])
}

func (sc Scanner) peekNext() rune {
	if sc.isAtEnd() {
		return rune(0)
	}
	return rune(sc.source[sc.cursor+1])
}

func isAlpha(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlphaNum(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func (sc Scanner) isAtEnd() bool {
	return sc.cursor >= len(sc.source)
}
