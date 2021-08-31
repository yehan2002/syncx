package bitfield

//Bit32 flag bit
type Bit32 struct {
	flag  *Bitfield32
	shift uint
}

//Set set
func (f *Bit32) Set(v bool) {
	bit := uint32(1 << f.shift)
	if v {
		f.flag.OR(bit)
	} else {
		f.flag.Clear(bit)
	}
}

//Bit64 flag bit
type Bit64 struct {
	flag  *Bitfield64
	shift uint
}

//Set set
func (f *Bit64) Set(v bool) {
	bit := uint64(1 << f.shift)
	if v {
		f.flag.OR(bit)
	} else {
		f.flag.Clear(bit)
	}
}
