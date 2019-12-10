package BrainFuck

import (
	"bytes"
	"strings"
	"testing"
)

func TestBrainFuckMachine_SimpleCalculation(t *testing.T) {
	output := prepare("++++++---", t)
	if r, _, _ := output.ReadRune(); r != 3 {
		t.Fatalf("incorrect value, ecpected 3, given %d", r)
	}
}

func TestBrainFuckMachine_MovingForwardAndBackward(t *testing.T) {
	output := prepare("+>>+++++++>+++", t)
	for _, v := range []rune {1,7,3} {
		if r, _, _ := output.ReadRune(); r != v {
			t.Fatalf("incorrect value, ecpected %d, given %d", v, r)
		}
	}
}

func prepare(str string, t *testing.T) bytes.Buffer {
	input := strings.NewReader(str)
	inst := NewInterpreter(input)
	err := inst.Run()

	if err != nil {
		t.Fatalf("executing the instructions failed, got %+v", err)
	}

	return inst.Write()
}
