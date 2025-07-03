package enblock

type FindBlockInput struct {
	WeakHash   uint32
	StrongHash [16]byte
}

// RollingChecksum is the rsync‐style rolling checksum.
type RollingChecksum struct {
	s1 uint32 // sum of bytes
	s2 uint32 // sum of running sums
	n  int    // window size
}

// NewRolling computes the initial checksum for block (len = n).
func NewRolling(block []byte) *RollingChecksum {
	r := &RollingChecksum{n: len(block)}
	for _, bb := range block {
		r.s1 += uint32(bb)
		r.s2 += r.s1
	}
	return r
}

// Sum returns the current 32-bit weak checksum.
func (r *RollingChecksum) Sum() uint32 {
	return (r.s1 & 0xffff) | (r.s2 << 16)
}

// Roll slides the window by dropping oldByte and adding newByte,
// then returns the updated 32-bit checksum.
func (r *RollingChecksum) Roll(oldByte, newByte byte) uint32 {
	r.s1 += uint32(newByte) - uint32(oldByte)
	r.s2 += r.s1 - uint32(r.n)*uint32(oldByte)
	return r.Sum()
}
