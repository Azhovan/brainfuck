package BrainFuck

import (
	"strings"
	"testing"
)

func TestScanner_Read(t *testing.T) {
	r := strings.NewReader("<>   this is string [+++]")
	s := NewScanner(r)
	ch := s.Read()

	if ch != '<' {
		t.Errorf("expect < given %q", ch)
	}
}

func TestScanner_Scan(t *testing.T) {
	//below string contains long white space and three runes
	r := strings.NewReader("        [+]")
	s := NewScanner(r)

	// consume all white spaces
	s.Scan()

	// read  [ rune
	ch := s.Read()

	if ch != '[' {
		t.Errorf("expect [ given %q", ch)
	}

	// read  + rune
	ch = s.Read()
	if ch != '+' {
		t.Errorf("expect + given %q", ch)
	}

	// read the last rune
	ch = s.Read()
	if ch != ']' {
		t.Errorf("expect ] given %q", ch)
	}

}
