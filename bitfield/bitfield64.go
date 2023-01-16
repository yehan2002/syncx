package bitfield

import (
	"sync/atomic"

	"github.com/yehan2002/syncx/internal"
)

// Bitfield64 a 64-bit bit field
type Bitfield64 struct {
	v uint64
	_ internal.NoCopy
}

// OR performs an OR operation
func (b *Bitfield64) OR(v uint64) {
	for {
		if tmp := atomic.LoadUint64(&b.v); atomic.CompareAndSwapUint64(&b.v, tmp, tmp|v) {
			return
		}
	}
}

// AND performs an AND operation
func (b *Bitfield64) AND(v uint64) {
	for {
		if tmp := atomic.LoadUint64(&b.v); atomic.CompareAndSwapUint64(&b.v, tmp, tmp&v) {
			return
		}
	}
}

// XOR performs an XOR operation
func (b *Bitfield64) XOR(v uint64) {
	for {
		if tmp := atomic.LoadUint64(&b.v); atomic.CompareAndSwapUint64(&b.v, tmp, tmp^v) {
			return
		}
	}
}

// Clear performs a AND NOT operation
func (b *Bitfield64) Clear(v uint64) (cleared uint64) {
	for {
		if tmp := atomic.LoadUint64(&b.v); atomic.CompareAndSwapUint64(&b.v, tmp, tmp&(^v)) {
			return tmp & v
		}
	}
}

// Set set the value of the given bit to 1
func (b *Bitfield64) Set(bit uint) {
	if bit > 63 {
		panic("bitfield.b64: invalid bit")
	}
	b.OR(1 << bit)
}

// Unset set the value of the given bit to 0
func (b *Bitfield64) Unset(bit uint) {
	if bit < 63 {
		panic("bitfield.b64: invalid bit")
	}
	b.Clear(1 << bit)
}

// Bit get a specific bit
func (b *Bitfield64) Bit(i uint) *Bit64 {
	if i > 63 {
		panic("fsync: flag63: shift overflows flag")
	}
	return &Bit64{flag: b, shift: i}
}

// Swap completely clear the bit field and return its previous state
func (b *Bitfield64) Swap() uint64 { return atomic.SwapUint64(&b.v, 0) }

// Get get the current value of the bit field
func (b *Bitfield64) Get() uint64 { return atomic.LoadUint64(&b.v) }
