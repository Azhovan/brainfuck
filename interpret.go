package BrainFuck

import (
	"bytes"
	"io"
)

type Interpreter interface {
	Run()
}

// BrainFuckMachine is contains w Writer, p Parser and memory
type BrainFuckMachine struct {
	p      *Parser
	w      bytes.Buffer
	i      io.Reader
	memory struct {
		cell [3000]int
		cu   int // memory cursor
	}
	err error
}

func NewInterpreter(r io.Reader) *BrainFuckMachine {
	return &BrainFuckMachine{
		p: NewParser(r),
	}
}

func (b *BrainFuckMachine) Run() (bytes.Buffer, error) {
	for _, inst := range b.p.Parse() {
		switch inst.i {
		case ">":
			b.move(inst.c)
		case "<":
			b.move(-1 * inst.c)
		case "+":
			b.inc(inst.c)
		case "-":
			b.inc(-1 *inst.c)
		case ".":
			b.doPrint(inst.c)
		case ",":
			b.doRead(inst.c)
		case "[":
			if b.val() == 0 {
				b.jump(inst.c)
				continue
			}
		case "]":
			if b.val() != 0 {
				b.jump(inst.c)
				continue
			}
		}

	}
	return b.w, b.err
}

// curr method returns current cursor position in the memory
func (b *BrainFuckMachine) cur() int {
	return b.memory.cu
}

// move method moves the cursor in the memory by given offset
func (b *BrainFuckMachine) move(offset int) {
	b.memory.cu += offset
}

// jump method forward the cursor to position p without any processing
func (b *BrainFuckMachine) jump(p int) {
	b.memory.cu = p
}

// inc method increment the value of the current cell in memory by v
// v can be positive or negative
func (b *BrainFuckMachine) inc(v int) {
	b.memory.cell[b.cur()] += v
}

// val method returns current value of which cursor is pointing
func (b *BrainFuckMachine) val() int {
	return b.memory.cell[b.cur()]
}

// doPrint method print the value in current cell of the memory
// it sets err and return false if any happened.
func (b *BrainFuckMachine) doPrint(times int) bool {
	v := byte(b.val())
	for i := 0; i < times; i++ {
		if err := b.w.WriteByte(v); err != nil {
			b.err = err
			return false
		}
	}
	return true
}

// doRead read input from io
// set error and return false if any happened
func (b *BrainFuckMachine) doRead(times int) bool {
	buf := make([]byte, 1)
	for i := 0; i < times; i++ {
		if _, err := b.i.Read(buf); err != nil {
			b.err = err
			return false
		}
		b.memory.cell[b.cur()] = int(buf[0])
	}
	return true
}
