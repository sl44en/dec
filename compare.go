package dec64

import (
// "fmt"
// "math"
// "math/bits"
// "strconv"
)

/*
// Floating point classification
const (
	decZero   = 1
	decNormal = 2
	decInf    = 4
	decNan    = 8
)
*/

func (a Dec64) Equal(b Dec64) bool {
	return a.asSortable(false) == b.asSortable(true)
}

func (a Dec64) Less(b Dec64) bool {
	return a.asSortable(false) < b.asSortable(true)
}

// these methods assume coefficients are normalised (and no subnormal numbers). With unnormalised
// coefficient, this will not work.

// return a sortable value for a dec64. The coefficient is set to c - lower bound + 1.
// this leaves mantissa range in [1 - 900000...] which can be held in 53 bits
// exponent is kept in same range and place as specified for type 1 decimal 64
// negative values are obtain be negating the above encoding.
// infinity has same encoding as in decimal 64

// The range is sortable, so the next higher number is higher also in integer representation
// and all numbers has only one representation.
// The range is not consecutive (all integer numbers will not be valid reprsentations

// the Go zero value correspondes to zero in this representaion, but it is
// not expected that these values will be created directly.
// the encoding is created to simplify sorting and comparison.

// lownan = true will make nan lowest negative value
func (a Dec64) asSortable(lownan bool) int64 {

	// one below lower bound, for shorting normalised mantissa
	const LOWERBOUND_1 = E_DIGITS_1 - 1

	// what is highests mantissa
	//prn("highest mnatissa", uint64(E_DIGITS_M - LOWERBOUND_1))
	//prn("lowest mnatissa", uint64(E_DIGITS_1 - LOWERBOUND_1))

	t1, s1, e1, m1 := decode(a)

	switch t1 {
	case decNormal:

		// exp bias!
		// exp is an int: int(ufe) - EXP_BIAS

		shortmant := m1 - LOWERBOUND_1
		// mantissa can be kept in 53 bits
		debug_assert(shortmant > 0)
		debug_assert(shortmant < (1 << EXP_SHIFT_T1))

		tempr := (uint64(e1+EXP_BIAS) << EXP_SHIFT_T1) | shortmant

		return transferSign(s1, tempr)

	case decZero:
		return 0

	case decInf:
		// sign!
		return transferSign(s1, uint64(INF_PATTERN))

	case decNan:
		return ternaryIf(lownan, -int64(NAN_PATTERN), int64(NAN_PATTERN))
	}

	panic("we dont get here")
}

func transferSign(s, tempr uint64) int64 {

	// a complete transfer of sign should take abs()
	// but we know input is unsigned

	/*
		benchmark measures same time ...
		This smart no branch version is not faster..
	*/
	t := int64(s) >> 63
	return (int64(tempr) ^ t) - t

	/*
	   	if s == 0 {
	   		return int64(tempr)
	   	}

	   return int64(-tempr)
	*/
}

// return -1, 0, 1 for a < b, a == b, a > b
// as per cmp.Compare function:
// a NaN is considered less than any non-NaN, a NaN is considered equal to a NaN, and -0.0 is equal to 0.0

func Compare(a, b Dec64) int {

	sa := a.asSortable(true)
	sb := b.asSortable(true)
	if sa < sb {
		return -1
	}
	if sa > sb {
		return 1
	}
	return 0
}

// TODO: the library also has the below interpretation of Less
// it is clearly more operational than the NaN googlydook that is ieee 754

// maybe Less() should have this interpretation, and
// PuristLess() could be the ieee complying version

// Less reports whether x is less than y.
// For floating-point types, a NaN is considered less than any non-NaN,
// and -0.0 is not less than (is equal to) 0.0.
