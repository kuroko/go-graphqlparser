package lexer

import (
	"bufio"
	"bytes"
	"fmt"
)

type Lexer struct {
	// TODO(elliot): Does the lexer contain things like line number? It will contain the position in
	// the input string. So, maybe we can also keep track of line number and column on that line? We
	// should be able to identify when a newline occurs after all. This would still need to work
	// when inside a string. This would mean that tokens could include their position by line and
	// column. It would be able to be worked out by the start and end position though later on, so
	// maybe it just makes sense to store that info, and worry about line numbers later on?
}

func (l *Lexer) Foo(input string) {
	rdr := bufio.NewReader(bytes.NewBufferString(input))

	// It's possible to unread runes:
	// rdr.UnreadRune()
	//

	for {
		r, _, err := rdr.ReadRune()
		if err != nil {
			break
		}

		fmt.Println(string(r))
	}
}
