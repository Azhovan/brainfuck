package BrainFuck

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

// Scanner provides a way to read data
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new scanner to read data from r.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		bufio.NewReader(r),
	}
}

// Read method read the next rune from r.
// err != nil only there is no more rune to read
func (s *Scanner) Read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return EOF
	}
	return ch
}

// Unread unreads the last data from r.
func (s *Scanner) Unread() {
	_ = s.r.UnreadRune()
}

// scanWhitespace ignores all consecutive whitespace
// the methods interrupts whenever reaches to end of data or
// extracts non whitespace rune
func (s *Scanner) scanWhitespace() (Token, string) {
	var buf bytes.Buffer
	for {
		ch := s.Read()
		if ch == EOF {
			break
		} else if !isWhitespace(ch) {
			s.Unread()
			break
		}
		buf.WriteRune(ch)
	}
	return WS, buf.String()
}

// scanIllegal ignore all consecutive illegal runes 
func (s *Scanner) scanIllegal() (Token, string) {
	var buff bytes.Buffer
	for {
		ch := s.Read()
		if ch == EOF {
			break
		} else if !isIllegal(ch) {
			s.Unread()
			break
		}
		buff.WriteRune(ch)
	}

	return ILLEGAL, buff.String()
}

// Scan returns the next Token and literal value
func (s *Scanner) Scan() (Token, string) {

	// read next rune
	ch := s.Read()

	// Ignore all consecutive whitespaces
	if isWhitespace(ch) {
		return s.scanWhitespace()

		// Ignore letters, digits
	} else if isIllegal(ch) {
		return s.scanIllegal()
	}

	// process the legal runes
	switch ch {
	case '>':
		return RIGHT, string(ch)
	case '<':
		return LEFT, string(ch)
	case '+':
		return PLUS, string(ch)
	case '-':
		return MINUS, string(ch)
	case '[':
		return LBRACK, string(ch)
	case ']':
		return RBRACK, string(ch)
	case '.':
		return PRINT, string(ch)
	case ',':
		return READ, string(ch)
	default:
		return ILLEGAL, "<nil>"
	}
}


// isWhitespace returns True if ch is space, tab, new-line
func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

// isIllegal return True if ch is letter or digits
func isIllegal(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch)
}
