package token

// Token types.
type TokenType string

// List of token types.
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	// Identifier and literals
	IDENT = "IDENT"
	INT   = "INT"
	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	GT       = ">"
	LT       = "<"
	BANG     = "!"
	EQ       = "=="
	NEQ      = "!="
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
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"
)

type Token struct {
	Type    TokenType
	Literal string
}

// Constructor function for tokens.
func New(typ TokenType, literal byte) Token {
	return Token{typ, string(literal)}
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
