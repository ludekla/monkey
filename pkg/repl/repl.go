package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/pkg/lexer"
	"monkey/pkg/token"
)

const PROMPT = "> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	// REPL: Read Eval Print Loop
	fmt.Fprint(out, PROMPT)
	for scanner.Scan() {
		line := scanner.Text()
		lx := lexer.New(line)
		for tok := lx.Next(); tok.Type != token.EOF; tok = lx.Next() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
		fmt.Fprint(out, PROMPT)
	}
}
