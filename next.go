package dec64

import (
// "fmt"
//
//	"math"
//
// "math/bits"
)

// -----------------------------------------------------------------------------
// Arithmetic
// -----------------------------------------------------------------------------

// / next floating point value towards positive infinity
func (self Dec64) NextUp() Dec64 {

	t, s, e, c := decode(self)

	switch t {
	case decNormal:
		if s != 0 {
			// this will go to zero if underflow normal range
			if c == E_DIGITS_1 {
				// construct with underflow check
				// TODO this fn used nowhere else
				return encodeNormalised(s, e-1, E_DIGITS-1)
			} else {
				return encodeFinal(s, e, c-1)
			}
		} else {
			// this will return inf if overflow normal range
			c++
			if c == E_DIGITS {
				// Overflowed range; next lower in range
				return encodeNormalised(0, e+1, E_DIGITS_1)
			} else {
				return encodeFinal(0, e, c)
			}
		}

	case decZero:

		// what is min positive
		return Dec64{MIN_POSITIVE}
	case decInf:

		if s != 0 { // not right!
			return Inf(1)
		} else {
			return Inf(-1)
		}

		/*
		   							      DecClass::Inf => {
		                               if self.is_signed() {
		                                   Self::MIN
		                               } else {
		                                   Self::INFINITY
		                               }
		                           }  */

	default:
		return NaN()
	}
}
