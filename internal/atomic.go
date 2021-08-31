package internal

//NoCopy a struct that prevents shallow copying and comparing
type NoCopy struct{ _ func() }

//Lock noop
func (*NoCopy) Lock() {}

//Unlock noop
func (*NoCopy) Unlock() {}
