package BrainFuck

import (
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	input := strings.NewReader("+++++")
	p := NewParser(input)
	instructions := p.Parse()

	// since we are folding instructions
	// there is one instruction +, but 4 times
	if len(instructions) != 1 {
		t.Errorf("wrong instruction length, expected 1 given %d", len(instructions))
	}

	expected := []*inst{
		&inst{c: 5, i: "+"},
	}

	for i, v := range expected {
		if *v != *instructions[i] {
			t.Errorf("wrong instruction literal. expected %+v given %+v", *v, *instructions[i])
		}
	}

}
