package BrainFuck

import (
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	input := strings.NewReader("+++++ -- [-]")
	p := NewParser(input)
	instructions := p.Parse()
	// since we are folding instructions
	// there is one instruction +, but 4 times
	if len(instructions) != 5 {
		t.Errorf("incorrect instruction's length, expected 5 given %+v", len(instructions))
	}
	expected := []*inst{
		&inst{c: 5, i: "+"},
		&inst{c: 2, i: "-"},
		&inst{c: 4, i: "["},
		&inst{c: 1, i: "-"},
		&inst{c: 2, i: "]"},
	}
	for i, v := range expected {
		if *v != *instructions[i] {
			t.Errorf("incorrect instruction. expected %+v given %+v", *v, *instructions[i])
		}
	}
}

func TestInnerLoops(t *testing.T) {
	input := strings.NewReader("-[--[+]--]")
	p := NewParser(input)
	instructions := p.Parse()
	expected := []*inst{
		{i:"-", c:1},
		{i:"[", c:7},
		{i:"-", c:2},
		{i:"[", c:5},
		{i:"+", c:1},
		{i:"]", c:3},
		{i:"-", c:2},
		{i:"]", c:1},
	}

	for i, v := range expected {
		if *v != *instructions[i] {
			t.Errorf("incorrect instruction. expected %+v given %+v", *v, *instructions[i])
		}
	}
}