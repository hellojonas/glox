package scanner

import "errors"

var errUnterminatedString = errors.New("unterminated string.")

const (
	// single character token
	LEFT_PAREN = iota + 1
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

var tokenMap = map[int]string {
	LEFT_PAREN: "LEFT_PAREN",
	RIGHT_PAREN: "RIGHT_PAREN",
	LEFT_BRACE: "LEFT_BRACE",
	RIGHT_BRACE: "RIGHT_BRACE",
	COMMA: "COMMA",
	DOT: "DOT",
	MINUS: "MINUS",
	PLUS: "PLUS",
	SEMICOLON: "SEMICOLON",
	STAR: "STAR",
	BANG: "BANG",
	BANG_EQUAL: "BANG_EQUAL",
	EQUAL: "EQUAL",
	EQUAL_EQUAL: "EQUAL_EQUAL",
	GREATER: "GREATER",
	GREATER_EQUAL: "GREATER_EQUAL",
	LESS: "LESS",
	LESS_EQUAL: "LESS_EQUAL",
	SLASH: "SLASH",
	IDENTIFIER: "IDENTIFIER",
	STRING: "STRING",
	NUMBER: "NUMBER",
	AND: "AND",
	CLASS: "CLASS",
	ELSE: "ELSE",
	FALSE: "FALSE",
	FUN: "FUN",
	FOR: "FOR",
	IF: "IF",
	NIL: "NIL",
	OR: "OR",
	PRINT: "PRINT",
	RETURN: "RETURN",
	SUPER: "SUPER",
	THIS: "THIS",
	TRUE: "TRUE",
	VAR: "VAR",
	WHILE: "WHILE",
	EOF: "EOF",
}
