package BrainFuck

// Token, represent a lexical tokens
type Token int

const (

	// Language Keywords
	LEFT   Token = iota // <
	RIGHT               // >
	PLUS                // +
	MINUS               // -
	PRINT               // .
	READ                // ,
	LBRACK              // [
	RBRACK              // ]

	WS      // white space
	ILLEGAL // not defined token
)

// special var to indicate end of stream
var EOF = rune(0)

