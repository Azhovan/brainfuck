package BrainFuck

import (
	"io"
)

// RuneParser will parse tokens and pack them in instructions
// initial state of the RuneParser is Parse method
type RuneParser interface {
	scan() (Token, string)
	Parse() []*inst
	NewParser(r io.Reader) *Parser
}

// inst is an abstraction for an operation which
// machine can understand
// i is one single instruction
// c is complementary information about instruction like position or counts of occurrence
type inst struct {
	i string
	c int
}

// Parser build AST (abstract structure tree)
// it contains the Scanner to tokenize the data from input
// buf is an internal struct to process input at a time of scan
// inst is an slice, which every member is one single instruction
type Parser struct {
	s    *Scanner
	inst []*inst
	buf  struct {
		t   Token  // last read token
		lit string // last read literal
		n   int    // buffer size
	}
}

// NewParser create new Parser from input r.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) Parse() []*inst {
	// fast way keep track of loops by saving
	// the position of the of the LBRACK(left bracket) and RBRACK(right bracket)
	// O(c) complexity
	stack := []int{}

	for {
		tok, lit := p.scan()
		if lit == "<nil>" {
			break
		}

		switch tok {
		case RIGHT, LEFT, PLUS, MINUS, PRINT, READ:
			p.addInst(tok, lit)
		case LBRACK:
			openLoop := p.buildInst(lit, 0)
			stack = append(stack, openLoop)
		case RBRACK:
			// pop off the position of the last [
			openLoop := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			closeLoop := p.buildInst(lit, openLoop)
			p.inst[openLoop].c = closeLoop
		}
	}
	return p.inst
}

// scan method ignore invalid token returned from the tscan
func (p *Parser) scan() (Token, string) {
	tok, lit := p.tscan()
	if tok == ILLEGAL || tok == WS {
		tok, lit = p.tscan()
	}
	return tok, lit

}

// tscan method (t stand for token) returns the next token from the scanner
// if the previous token was not consumed properly, it will be returned again and buffer will be cleared
func (p *Parser) tscan() (Token, string) {
	// there is a token on the buffer
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.t, p.buf.lit
	}

	// read the next token from the s
	tok, lit := p.s.Scan()

	// buffer it in buf
	p.buf.t, p.buf.lit = tok, lit

	return tok, lit
}

// unscan sends the already consumed token back to buff
func (p *Parser) unscan() {
	p.buf.n = 1
}

// addInst adds instructions to []*inst of Parser
// for efficiency, if there are multiple occurrences of the
// same token consecutively, we will fold it
func (p *Parser) addInst(token Token, literal string) int {
	// token occurrence count
	c := 1
	for {
		tok, _ := p.scan()
		if tok != token {
			p.unscan()
			break
		}
		c++
	}

	return p.buildInst(literal, c)
}

// buildInst create a instruction from the given literals
func (p *Parser) buildInst(literal string, c int) int {
	// build instruction
	inst := &inst{
		i: literal,
		c: c,
	}
	// add inst to instruction list
	p.inst = append(p.inst, inst)
	return len(p.inst) - 1
}
