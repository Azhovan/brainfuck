package brainfuck

import (
	"io"
)

// RuneParser will parse tokens and pack them in instructions
// initial state of the RuneParser is Parse method
type RuneParser interface {
	Parse() []*inst
}

// inst is an abstraction for an operation which machine can understand
// i is one single instruction
// c is complementary information about instruction like position or counts of occurrence
type inst struct {
	t Token
	c int
}

// Parser builds AST (abstract structure tree).
// Parser uses Stack to keep track of loops
// it contains the Scanner to tokenize the data from input
// buf is an internal struct to process input at a time of scan
// inst is an slice, which every member is one single instruction
type Parser struct {
	s    *Scanner
	inst []*inst
	buf  struct {
		tok     Token // last read token
		tokbufn bool  // whether the token buffer is in use.
	}
	stack Stack
}

// NewParser creates new Parser from input r.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) Parse() []*inst {
	for {
		tok := p.scan()
		if tok.Tok == IllegalToken {
			break
		}
		switch tok.Tok {
		case
			RightToken,
			LeftToken,
			PlusToken,
			MinusToken,
			PrintToken,
			ReadToken:
			p.addInst(tok)
		case LeftBracketToken:
			openLoop := p.buildInst(tok, 0)
			p.stack.Push(openLoop)
		case RightBracketToken:
			openLoop := p.stack.Pop().(int)
			closeLoop := p.buildInst(tok, openLoop)
			p.inst[openLoop].c = closeLoop
		}

	}
	return p.inst
}

// scan returns next token unit.
func (p *Parser) scan() Token {
	// there is a token on the buffer
	if p.buf.tokbufn {
		p.buf.tokbufn = false
		return p.buf.tok
	}
	// read the next token from s
	tok := p.s.Scan()
	p.buf.tok = tok
	return tok
}

// unscan sends the already consumed token back to buff.
func (p *Parser) unscan() {
	p.buf.tokbufn = true
}

// addInst adds instructions to []*inst of Parser
// for efficiency, if there are multiple occurrences of the
// same token consecutively, we will fold it.
func (p *Parser) addInst(t Token) int {
	// token occurrence count
	c := 1
	for {
		next := p.scan()
		if next.Tok != t.Tok {
			p.unscan()
			break
		}
		c++
	}
	return p.buildInst(t, c)
}

// buildInst creates a instruction from the given literals.
func (p *Parser) buildInst(t Token, c int) int {
	// build instruction
	inst := &inst{
		t: t,
		c: c,
	}
	// add inst to instruction list
	p.inst = append(p.inst, inst)
	return len(p.inst) - 1
}
