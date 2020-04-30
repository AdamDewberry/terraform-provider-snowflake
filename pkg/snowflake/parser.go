package snowflake

import (
	"fmt"
	"unicode"
)

// ViewSelectStatementExtractor is a simplistic parser that only exists to extract the select
// statement from a create view statement
//
// The implementation is optimized for undertandable and predictable behavior. So far we only seek
// to support queries of the sort that are generated by this project.
type ViewSelectStatementExtractor struct {
	input []rune
	pos   int
}

func NewViewSelectStatementExtractor(input string) *ViewSelectStatementExtractor {
	return &ViewSelectStatementExtractor{
		input: []rune(input),
	}
}

func (e *ViewSelectStatementExtractor) Extract() (string, error) {
	e.consumeSpace()
	e.consumeToken("create")
	e.consumeSpace()
	e.consumeToken("or replace")
	e.consumeSpace()
	e.consumeToken("secure")
	e.consumeSpace()
	e.consumeToken("recursive")
	e.consumeSpace()
	e.consumeToken("view")
	e.consumeSpace()
	e.consumeToken("if not exists")
	e.consumeSpace()
	e.consumeIdentifier()
	// TODO column list
	// TODO copy grants
	e.consumeComment()
	e.consumeSpace()
	e.consumeComment()
	e.consumeSpace()
	e.consumeToken("as")
	e.consumeSpace()

	return string(e.input[e.pos:]), nil
}

// consumeToken will move e.pos forward iff the token is the next part of the input. Comparison is
// case-insensitive. Will return true if consumed.
func (e *ViewSelectStatementExtractor) consumeToken(t string) bool {
	fmt.Printf("consume token %s\n", t)
	found := 0
	for i, r := range t {
		fmt.Printf("e.pos %d r %s\n", e.pos, string(r))
		if e.pos+i > len(e.input) || r != e.input[e.pos+i] {
			break
		}
		found += 1
	}
	fmt.Printf("found %d\n", found)

	if found == len(t) {
		e.pos += len(t)
		return true
	}
	return false
}

func (e *ViewSelectStatementExtractor) consumeSpace() {
	found := 0
	for {
		if e.pos+found > len(e.input)-1 || !unicode.IsSpace(e.input[e.pos+found]) {
			break
		}
		found += 1
	}
	e.pos += found
}

func (e *ViewSelectStatementExtractor) consumeIdentifier() {
	// TODO quoted identifiers
	e.consumeNonSpace()
}

func (e *ViewSelectStatementExtractor) consumeNonSpace() {
	found := 0
	for {
		if e.pos+found > len(e.input)-1 || unicode.IsSpace(e.input[e.pos+found]) {
			break
		}
		found += 1
	}
	e.pos += found
}
func (e *ViewSelectStatementExtractor) consumeComment() {
	if c := e.consumeToken("comment"); !c {
		return
	}

	e.consumeSpace()

	if c := e.consumeToken("="); !c {
		return
	}

	e.consumeSpace()

	if c := e.consumeToken("'"); !c {
		return
	}

	found := 0
	escaped := false
	for {
		if e.pos+found > len(e.input)-1 {
			break
		}
		fmt.Printf("e.pos %d found %d escaped %t r %s\n", e.pos, found, escaped, string(e.input[e.pos+found]))

		if escaped {
			escaped = false
		} else if e.input[e.pos+found] == '\\' {
			escaped = true
		} else if e.input[e.pos+found] == '\'' {
			break
		}
		found += 1
	}
	e.pos += found

	if c := e.consumeToken("'"); !c {
		return
	}
}
