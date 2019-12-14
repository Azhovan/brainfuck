package main

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

// LexReader is an interface that wraps Read method
// Read method reads and return the next rune from the input
type LexReader interface {
	Read() rune
}

// LexScanner is the interface that adds Unread method to the
// basic LexReader
//
// Unread causes the next call to the Read method return the same
// rune as the same previous call to Read
type LexScanner interface {
	LexReader
	Unread() error
}

// Scanner implements a tokenizer.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		bufio.NewReader(r),
	}
}

// Read method read the next rune from r.
// err != nil only there is no more rune to read
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return EOF
	}
	return ch
}

// Unread unreads the last data from r.
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

// scanWhitespace consumes all subsequent whitespace
func (s *Scanner) scanWhitespace() Token {
	var buf bytes.Buffer
	for {
		ch := s.read()
		if ch == EOF {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		}
		_, _ = buf.WriteRune(ch)
	}
	return Token{Tok: WhitespaceToken, Value: buf.String()}
}

// scanLetterDigit ignore all subsequent letters or digits
func (s *Scanner) scanLetterDigit() Token {
	var buff bytes.Buffer
	for {
		ch := s.read()
		if ch == EOF {
			break
		} else if !isLetterDigit(ch) {
			s.unread()
			break
		}
		buff.WriteRune(ch)
	}

	return Token{Tok: IllegalToken, Value: buff.String()}
}

// Scan returns the next Token and literal value from the reader.
func (s *Scanner) Scan() Token {

	// read next rune
	ch := s.read()

	// If whitespace code point found, then consume all contiguous whitespaces.
	if isWhitespace(ch) {
		return s.scanWhitespace()
	}

	// If letter, digit code point found, then consume all letters, digits
	if isLetterDigit(ch) {
		return s.scanLetterDigit()
	}

	return s.next(ch)
}

func (s *Scanner) next(ch rune) Token {
	// Check against individual code points next.
	switch ch {
	case '>':
		return Token{Tok: RightToken, Value: string(ch)}
	case '<':
		return Token{Tok: LeftToken, Value: string(ch)}
	case '+':
		return Token{Tok: PlusToken, Value: string(ch)}
	case '-':
		return Token{Tok: MinusToken, Value: string(ch)}
	case '[':
		return Token{Tok: LeftBracketToken, Value: string(ch)}
	case ']':
		return Token{Tok: RightBracketToken, Value: string(ch)}
	case '.':
		return Token{Tok: PrintToken, Value: string(ch)}
	case ',':
		return Token{Tok: ReadToken, Value: string(ch)}
	default:
		return Token{Tok: IllegalToken, Value: "<nil>"}
	}
}

// isWhitespace returns True if ch is space, tab, new-line
func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

// isLetterDigit return True if ch is letter or digit
func isLetterDigit(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch)
}