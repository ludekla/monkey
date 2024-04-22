package token

// Token types.
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	// Identifier and literals
	IDENT = "IDENT"
	INT   = "INT"
	// Operators
	ASSIGN = "="
	PLUS   = "+"
	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

type Token struct {
	Type    string
	Literal string
}
