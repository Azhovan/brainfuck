package BrainFuck

import (
	"strings"
	"testing"
)

func TestBrainFuckMachine_SimpleCalculation(t *testing.T) {
	input := strings.NewReader("+++---+++")
	int := NewInterpreter(input)
	_, err := int.Run()
	if err != nil {
		t.Fatal(err)
	}
	if int.memory.cell[0] != 3 {
		t.Fatalf("incorrect value, ecpected 3, given %d", int.memory.cell[0])
	}
}
