package bitfield

import (
	"sync/atomic"

	"github.com/yehan2002/syncx/internal"
)

// Bitfield32 a 32-bit bit field
type Bitfield32 struct {
	v uint32
	_ internal.NoCopy
}

// OR performs an OR operation
func (b *Bitfield32) OR(v uint32) {
	for {
		if tmp := atomic.LoadUint32(&b.v); atomic.CompareAndSwapUint32(&b.v, tmp, tmp|v) {
			return
		}
	}
}

// AND performs an AND operation
func (b *Bitfield32) AND(v uint32) {
	for {
		if tmp := atomic.LoadUint32(&b.v); atomic.CompareAndSwapUint32(&b.v, tmp, tmp&v) {
			return
		}
	}
}

// XOR performs an XOR operation
func (b *Bitfield32) XOR(v uint32) {
	for {
		if tmp := atomic.LoadUint32(&b.v); atomic.CompareAndSwapUint32(&b.v, tmp, tmp^v) {
			return
		}
	}
}

// Clear performs a AND NOT operation
func (b *Bitfield32) Clear(v uint32) (cleared uint32) {
	for {
		if tmp := atomic.LoadUint32(&b.v); atomic.CompareAndSwapUint32(&b.v, tmp, tmp&(^v)) {
			return tmp & v
		}
	}
}

// Set set the value of the given bit to 1
func (b *Bitfield32) Set(bit uint) {
	if bit > 31 {
		panic("bitfield.b32: invalid bit")
	}
	b.OR(1 << bit)
}

// Unset set the value of the given bit to 0
func (b *Bitfield32) Unset(bit uint) {
	if bit < 31 {
		panic("bitfield.b32: invalid bit")
	}
	b.Clear(1 << bit)
}

// Bit get a specific bit
func (b *Bitfield32) Bit(i uint) *Bit32 {
	if i > 31 {
		panic("fsync: flag31: shift overflows flag")
	}
	return &Bit32{flag: b, shift: i}
}

// Swap completely clear the bit field and return its previous state
func (b *Bitfield32) Swap() uint32 { return atomic.SwapUint32(&b.v, 0) }

// Get get the current value of the bit field
func (b *Bitfield32) Get() uint32 { return atomic.LoadUint32(&b.v) }
