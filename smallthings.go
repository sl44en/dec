package dec64

// three way if....
func ternaryIf[T any](p bool, a, b T) T {
	if p {
		return a
	}
	return b
}

// --------------------------------------------------------
// assertions
// --------------------------------------------------------

func debug_assert(p bool) {
	if !p {
		panic("Assertion failed")
	}
}

/*
// non intrusive version of assert
// what is the overhead from this
func debug_assert(_ bool) {
}
*/

// --------------------------------------------------------
// abs() uabs()
// --------------------------------------------------------

// returns absolute value as an unsigned
// it has the nice property that it does not overflow for MinInt64
func uabs(x int64) uint64 {
	y := x >> 63
	return uint64((x ^ y) - y)
}

// return sign bit of a signed integer
func usign(x int64) uint64 {
	const signmask = 1 << 63
	return signmask & uint64(x)
}
