package main

import (
	"io"
)

// interface to write out the execution results
type Writer interface {
	Write() io.Writer
}

// interface for an interpreter
// Run method executes created instructions by Parser
type Interpreter interface {
	Writer
	Run() error
}

// Memory capacity
const MemorySize int = 3000

// BrainFuck is an implementation of the Interpreter
// it has internal parser which builds instructions from the input
// result is written into w
// memory struct keep memory data and cursor to move between memory cells and update their value
// err != nil if any error happen during the print/read operation
type BrainFuck struct {
	p      *Parser
	w      io.Writer
	i      io.Reader
	buf    []byte
	ip     int
	err    error
	memory struct {
		cell [MemorySize]int
		cu   int
	}
}

// NewInterpreter create new Interpreter instance and initialize it's internal Parser.
func NewInterpreter(i io.Reader, w io.Writer, parser *Parser) *BrainFuck {
	return &BrainFuck{
		p:   parser,
		w:   w,
		i:   i,
		buf: make([]byte, 1),
	}
}

// Run method executes the instructions
// err != nil if error happen during read/print operations
// output returns in format of bytes
func (b *BrainFuck) Run() error {

	inst := b.p.Parse()

	for b.ip < len(inst) {
		switch inst[b.ip].t.Value {
		case ">":
			b.seek(inst[b.ip].c)
		case "<":
			b.seek(-inst[b.ip].c)
		case "+":
			b.inc(inst[b.ip].c)
		case "-":
			b.dec(inst[b.ip].c)
		case ".":
			b.write(inst[b.ip].c)
		case ",":
			b.read(inst[b.ip].c)
		case "[":
			if b.val() == 0 {
				b.jump(inst[b.ip].c)
				continue
			}
		case "]":
			if b.val() != 0 {
				b.jump(inst[b.ip].c)
				continue
			}
		}
		b.ip++
	}

	return b.err
}

// Write method writes memory into buffer
//func (b *BrainFuck) Write() io.Writer {
//	b.reset()
//	for {
//		if v, err := b.scan(); err != nil {
//			break
//		} else if v != 0 {
//			b.buf[0] = v
//			_, _ = b.w.Write(b.buf)
//		}
//	}
//	return b.w
//}

// scan method return next valid value in memory
// and move the memory cursor forward
// scan ignores zero filled cell
//func (b *BrainFuck) scan() (byte, error) {
//	if b.cur() >= MemorySize-1 {
//		return byte(0), io.EOF
//	}
//	b.seek(1)
//	if b.val() > 0 {
//		return byte(b.val()), nil
//	}
//	return byte(0), nil
//}

// curr method returns the position of current cursor in the memory
func (b *BrainFuck) cur() int {
	return b.memory.cu
}

// seek method moves the cursor in the memory to given offset
// this move is relative to current cursor position
func (b *BrainFuck) seek(offset int) {
	b.memory.cu += offset
}

// jump method forward the cursor to position p
func (b *BrainFuck) jump(p int) {
	b.ip = p
}

// reset method resets the cursor and writer to point to invalid state
func (b *BrainFuck) reset() {
	b.memory.cu = 0
}

// inc method increment the value of the current cell in memory by v
func (b *BrainFuck) inc(v int) {
	b.memory.cell[b.cur()] = (b.memory.cell[b.cur()] + v) % 255
}

func (b *BrainFuck) dec(v int) {
	if b.memory.cell[b.cur()]-v >= 0 {
		b.memory.cell[b.cur()] -= v
	} else {
		b.memory.cell[b.cur()] = 256 + b.memory.cell[b.cur()] - v
	}
}

// val method returns current value of which cursor is pointing
func (b *BrainFuck) val() int {
	return b.memory.cell[b.cur()]
}

// doPrint method print the value in current cell of the memory
// if any error happen during the Write operation err property will be set
func (b *BrainFuck) write(times int) bool {
	b.buf[0] = byte(b.val())
	for i := 0; i < times; i++ {
		if _, err := b.w.Write(b.buf); err != nil {
			b.err = err
			return false
		}
	}
	return true
}

// doRead read input from io
// if any error happen during the Read operation err property will be set
func (b *BrainFuck) read(times int) bool {
	for i := 0; i < times; i++ {
		if _, err := b.i.Read(b.buf); err != nil {
			b.err = err
			return false
		}
		b.memory.cell[b.cur()] = int(b.buf[0])
	}
	return true
}
