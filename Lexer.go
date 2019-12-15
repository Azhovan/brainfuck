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
	read() rune
}

// LexScanner is the interface that adds Unread method to the
// basic LexReader
//
// Unread causes the next call to the Read method return the same
// rune as the same previous call to Read
type LexScanner interface {
	LexReader
	unread() error
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
// err != nil only if there is no more rune to read
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return EOF
	}
	return ch
}

// unread re-buffer the last read data.
func (s *Scanner) unread() error {
	if err := s.r.UnreadRune(); err != nil {
		return err
	}
	return nil
}

// scanWhitespace consumes all subsequent whitespace.
func (s *Scanner) scanWhitespace() Token {
	var buf bytes.Buffer
	for {
		ch := s.read()
		if ch == EOF {
			break
		} else if !isWhitespace(ch) {
			_ = s.unread()
			break
		}
		_, _ = buf.WriteRune(ch)
	}
	return Token{Tok: WhitespaceToken, Value: buf.String()}
}

// scanLetterDigit ignore all subsequent letters or digits.
func (s *Scanner) scanLetterDigit() Token {
	var buff bytes.Buffer
	for {
		ch := s.read()
		if ch == EOF {
			break
		} else if !isLetterDigit(ch) {
			_ = s.unread()
			break
		}
		buff.WriteRune(ch)
	}

	return Token{Tok: IllegalToken, Value: buff.String()}
}

// Scan prepare and returns the next Token.
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
