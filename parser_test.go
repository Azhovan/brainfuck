package main

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
		t.Errorf("incorrect instruction's length, expected 5 got %+v", len(instructions))
	}
	expected := []*inst{
		&inst{c: 5, t: Token{Tok:PlusToken, Value:"+"}},
		&inst{c: 2, t: Token{Tok:MinusToken, Value:"-"}},
		&inst{c: 4, t: Token{Tok:LeftBracketToken, Value:"["}},
		&inst{c: 1, t: Token{Tok:MinusToken, Value:"-"}},
		&inst{c: 2, t: Token{Tok:RightBracketToken, Value:"]"}},
	}
	for i, v := range expected {
		if *v != *instructions[i] {
			t.Errorf("incorrect instruction. expected %+v got %+v", *v, *instructions[i])
		}
	}
}

func TestInnerLoops(t *testing.T) {
	input := strings.NewReader("-[--[+]--]")
	p := NewParser(input)
	instructions := p.Parse()
	expected := []*inst{
		{t:Token{Tok:MinusToken, Value:"-"}, c:1},
		{t:Token{Tok:LeftBracketToken, Value:"["}, c:7},
		{t:Token{Tok:MinusToken, Value:"-"}, c:2},
		{t:Token{Tok:LeftBracketToken, Value:"["}, c:5},
		{t:Token{Tok:PlusToken, Value:"+"}, c:1},
		{t:Token{Tok:RightBracketToken, Value:"]"}, c:3},
		{t:Token{Tok:MinusToken, Value:"-"}, c:2},
		{t:Token{Tok:RightBracketToken, Value:"]"}, c:1},
	}

	for i, v := range expected {
		if *v != *instructions[i] {
			t.Errorf("incorrect instruction. expected %+v got %+v", *v, *instructions[i])
		}
	}
}

func Test_MoveBetweenCells(t *testing.T) {
	input := strings.NewReader("+>>>+++++++>>+++ --<<")
	p := NewParser(input)
	instructions := p.Parse()
	expected := []*inst{
		{t:Token{Tok:PlusToken, Value:"+"}, c:1},
		{t:Token{Tok:RightToken, Value:">"}, c:3},
		{t:Token{Tok:PlusToken, Value:"+"}, c:7},
		{t:Token{Tok:RightToken, Value:">"}, c:2},
		{t:Token{Tok:PlusToken, Value:"+"}, c:3},
		{t:Token{Tok:MinusToken, Value:"-"}, c:2},
		{t:Token{Tok:LeftToken, Value:"<"}, c:2},
	}

	for i, v := range expected {
		if *v != *instructions[i] {
			t.Errorf("incorrect instruction. expected %+v got %+v", *v, *instructions[i])
		}
	}

}