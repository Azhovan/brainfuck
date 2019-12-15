package brainfuck

const (
	IllegalToken      Tok = iota
	LeftToken             // <
	RightToken            // >
	PlusToken             // +
	MinusToken            // -
	PrintToken            // .
	ReadToken             // ,
	LeftBracketToken      // [
	RightBracketToken     // ]
	WhitespaceToken
)

// special var to indicate end of stream.
var EOF = rune(-1)

// Tok represents a lexical token type.
type Tok int

// Token, represent a lexical tokens.
type Token struct {
	// the type of token.
	Tok Tok

	// The literal value of the token(as parsed).
	Value string

	// The rune used for string tokens
	Ending rune

	// Used for numeric tokens.
	Number float64
	Unit   string
}
