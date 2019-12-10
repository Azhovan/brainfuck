package BrainFuck

import (
	"bytes"
	"io"
)

// interface to write out the execution results
type Writer interface {
	Write() bytes.Buffer
}

// interface for an interpreter
// Run method executes created instructions by Parser
type Interpreter interface {
	Writer
	Run() error
}

// Memory capacity
const MemorySize int = 3000

// BrainFuckMachine is an implementation of the Interpreter
// it has internal parser which builds instructions from the input
// result is written into w
// memory struct keep memory data and cursor to move between memory cells and update their value
// err != nil if any error happen during the print/read operation
type BrainFuckMachine struct {
	p      *Parser
	w      bytes.Buffer
	i      io.Reader
	memory struct {
		cell [MemorySize]int
		cu   int // memory cursor
	}
	err error
}

// NewInterpreter create new Interpreter instance and initialize it's internal Parser with Reader r.
// Parser has internal functionality to create instructions based on the input
func NewInterpreter(r io.Reader) *BrainFuckMachine {
	return &BrainFuckMachine{
		p: NewParser(r),
	}
}

// Run method executes the instructions
// err != nil if error happen during read/print operations
// output returns in format of bytes
func (b *BrainFuckMachine) Run() error {
	for _, inst := range b.p.Parse() {
		switch inst.i {
		case ">":
			b.seek(inst.c)
		case "<":
			b.seek(-1 * inst.c)
		case "+":
			b.inc(inst.c)
		case "-":
			b.inc(-1 * inst.c)
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
	return b.err
}

// Write method writes memory into buffer
func (b *BrainFuckMachine) Write() bytes.Buffer {
	b.reset()
	for {
		if v, err := b.scan(); err != nil {
			break
		} else if v != 0 {
			b.w.WriteRune(rune(v))
		}
	}
	return b.w
}

// scan method return next valid value in memory
// and move the memory cursor forward
// scan ignores zero filled cell
func (b *BrainFuckMachine) scan() (int, error) {
	if b.cur() >= MemorySize-1 {
		return 0, io.EOF
	}
	b.seek(1)
	if b.val() > 0 {
		return b.val(), nil
	}
	return 0, nil
}

// curr method returns the position of current cursor in the memory
func (b *BrainFuckMachine) cur() int {
	return b.memory.cu
}

// seek method moves the cursor in the memory to given offset
// this move is relative to current cursor position
func (b *BrainFuckMachine) seek(offset int) {
	b.memory.cu += offset
}

// jump method forward the cursor to position p
func (b *BrainFuckMachine) jump(p int) {
	b.memory.cu = p
}

// reset method resets the cursor to point to nowhere until an move operation happened
func (b *BrainFuckMachine) reset() {
	b.memory.cu = -1
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
// if any error happen during the Write operation err property will be set
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
// if any error happen during the Read operation err property will be set
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
