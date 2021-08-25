package dfcache

type ByteView struct {
	b []byte
}

// Get the length of data
func (bv ByteView) Len() int {
	return len(bv.b)
}

// Get a copy of the data as a byte slice
func (bv ByteView) ByteSlice() []byte {
	return cloneBytes(bv.b)
}

// Get the data as a string
func (bv ByteView) String() string {
	return string(bv.b)
}

// Copy bytes
func cloneBytes(b []byte) []byte {
	var dst []byte
	return append(dst, b...)
}
