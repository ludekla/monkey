package lexer

type Lexer struct {
	input    string
	position int  // position in input string
	readPos  int  // current reading position, one after ch
	ch       byte // current char under examination
}

// New is the Lexer factory.
func New(input string) *Lexer {
	return &Lexer{input: input}
}

// readChar advanced the reader pointer to the next char
// within the input string.
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // go to the beginning
	} else {
		l.ch = l.input[l.readPos]
	}
	l.position = l.readPos
	l.readPos++
}
