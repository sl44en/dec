package dec64

import (
// "fmt"
// "math"
// "math/bits"
)

// -----------------------------------------------------------------------------
// Arithmetic  Add and Sub
// -----------------------------------------------------------------------------
/*
    * There is a complex system for the type of the result depending on
    * input types for calculations
    *
    * The table below list the result types for each combination.
    * Note that this implemenattion does not implement signed zeroes,
    * so zero results will be unsigned. Below table is based on how f64 type
    * interpretes input types and generate results.
    * -------------------------------------------------------------------------------

    x + y     y =            0          2     -56.89        inf       -inf        NaN
   x =          0            0          2     -56.89        inf       -inf        NaN
   x =          2            2          4     -54.89        inf       -inf        NaN
   x =     -56.89       -56.89     -54.89    -113.78        inf       -inf        NaN
   x =        inf          inf        inf        inf        inf        NaN        NaN
   x =       -inf         -inf       -inf       -inf        NaN       -inf        NaN
   x =        NaN          NaN        NaN        NaN        NaN        NaN        NaN

    x - y     y =            0          2     -56.89        inf       -inf        NaN
   x =          0            0         -2      56.89       -inf        inf        NaN
   x =          2            2          0      58.89       -inf        inf        NaN
   x =     -56.89       -56.89     -58.89          0       -inf        inf        NaN
   x =        inf          inf        inf        inf        NaN        inf        NaN
   x =       -inf         -inf       -inf       -inf       -inf        NaN        NaN
   x =        NaN          NaN        NaN        NaN        NaN        NaN        NaN
    * -------------------------------------------------------------------------------
*/

func (a Dec64) Sub(b Dec64) Dec64 {

	// this will negate NaN as well as other values, but that will not afect the result
	return a.Add(Dec64{b.d ^ SIGN_MASK})
}

func (a Dec64) Add(b Dec64) Dec64 {

	t1, s1, e1, c1 := decode(a)
	t2, s2, e2, c2 := decode(b)

	// the classification returns int 1,2,4,8
	// so can be combined as (t1 << 4) | t2

	// tcomb := (t1 << 4) | t2

	// todo consider having these logics in separate function
	// and do intial split between subm, add on sign only
	// keeping the normals operations together with this logic

	if s1 == s2 {
		// if same sign, add

		if t1|t2 == decNormal {

			if e1 >= e2 {
				// first parameter must have highest exponent
				//return a + b
				return addNormals(s1, e1, c1, e2, c2)
			} else {
				// return b + a
				return addNormals(s2, e2, c2, e1, c1)
			}
		}

		if max(t1, t2) <= decNormal {
			if t1 == decZero {
				// 0 + 0
				// 0 + b
				return b
			}
			// a + 0
			return a
		}

		if max(t1, t2) == decInf {
			// a and b have same sign
			// one or more of them is inf
			return Dec64{s1 | INF_PATTERN}
		}

		return NaN()

	} else {

		// different sign, subtract

		switch t1 | t2 {

		case decNormal: // both normal
			if e1 >= e2 {
				// first parameter must have highest exponent
				// return a - b, result with sign of a
				return subtractNormals(s1, e1, c1, e2, c2)
			} else {
				// return b - a, result with sign of b
				return subtractNormals(s2, e2, c2, e1, c1)
			}

		case decZero: // 0 - 0 is a special case
			// result is 0 regardless of input signs
			return Dec64{}

		case decZero | decNormal:
			// 0 - b, a - 0
			if t1 == decZero {
				return b
			}
			return a

		case decInf:
			// both are inf: inf - inf -> nan
			return NaN()
		}

		if max(t1, t2) == decInf {
			// inf vs finite values
			if t1 == decInf {
				// inf - finite -> inf
				return a
			}
			// finite - inf -> -inf
			return b
		}

		return NaN()

		panic("should not get here")
	}
}

/*
For add, there is two types of overrun checking
- coeeficient can overflow becuase of rounding. Needs to be scaled down; This can in turn trigger
  exponent overrun
- exponent overrun because of upwards exponent adjustment
*/

// encode value after rounding, exponent will be well behaved, from original value
// coefficient is normalised, except it can have overflowed
func encodeChkCoefficient(s uint64, e int, c DecBase) Dec64 {
	// assumptions: only if there is round up can exponent be too high

	debug_assert(c >= E_DIGITS_1 && c <= E_DIGITS)
	if c == E_DIGITS {
		if e >= EXP_MAX_NORM {
			return Dec64{s | INF_PATTERN}
		}
		return encodeFinal(s, e+1, E_DIGITS_1)
	}
	return encodeFinal(s, e, c)
}

// this assumes coefficient is in normalised in range
// it checks exponent value for overlow and return Inf if needed
func encodeChkExponent(s uint64, e int, c DecBase) Dec64 {

	debug_assert(c >= E_DIGITS_1 && c < E_DIGITS)
	if e > EXP_MAX_NORM {
		return Dec64{s | INF_PATTERN}
	}
	return encodeFinal(s, e, c)
}

// add two values of same sign
//func addNormals(a, b Dec64) Dec64 {

//	_, s1, e1, c1 := decode(a)
//	_, _, e2, c2 := decode(b)

func addNormals(s1 DecBase, e1 int, c1 DecBase, e2 int, c2 DecBase) Dec64 {

	// they need to have same sign, for this to be an add

	var c, sf, hi_c2, lo_c2, lo_c uint64

	// they need be sorted accoring to scale
	debug_assert(e1 >= e2)
	sdiff := uint64(e1 - e2)
	if sdiff <= 3 {
		switch sdiff {
		case 0:
			// same scale
			c = c1 + c2
			if c < E_DIGITS {
				// no carry from addition
				return encodeFinal(s1, e1, c)
			} else {
				// addition has caused coefficient overrun, adjust and round
				// This situation, should never trigger next range coefficient overrun
				// check exponent overrun
				return encodeChkExponent(s1, e1+1, roundEven(c/10, c%10, 10))
			}

		case 1:
			c = c1*10 + c2
			if c < E_DIGITS*10 {
				// no carry from addition
				// divisions are not parameterised, to use const division optimise
				return encodeChkCoefficient(s1, e1, roundEven(c/10, c%10, 10))
			} else {
				// carry from addition
				return encodeChkExponent(s1, e1+1, roundEven(c/100, c%100, 100))
			}
		case 2:
			c = c1*100 + c2
			if c < E_DIGITS*100 {
				// no carry from addition
				// divisions are not parameterised, to use const divisor
				return encodeChkCoefficient(s1, e1, roundEven(c/100, c%100, 100))
			} else {
				// carry from addition
				return encodeChkExponent(s1, e1+1, roundEven(c/1000, c%1000, 1000))
			}
		case 3:
			c = c1*1000 + c2
			if c < E_DIGITS*1000 {
				// no carry from addition
				// divisions are not parameterised, to use const divisor
				return encodeChkCoefficient(s1, e1, roundEven(c/1000, c%1000, 1000))
			} else {
				// carry from addition
				return encodeChkExponent(s1, e1+1, roundEven(c/10000, c%10000, 10000))
			}

			// these shorter versions are slower!
			// should be possible to handle sdiff 1- 3 this way 0 - 3, actually)
			// AlT: We could call encodeAny - there could be a dedicated optimiseeed encode high?
			// although.. maybe better to utilise division by constants
			/*
				} else if sdiff <= 3 {
					c = c1*Pow10(sdiff) + c2
					if c < E_DIGITS*Pow10(sdiff) { // could this be constant
						// no carry from addition
						return encodeChkCoefficient(s1, e1, roundEven(divPow10(c, sdiff)))
					} else {
						// carry from addition
						return encodeChkExponent(s1, e1+1, roundEven(divPow10(c, sdiff+1)))
					}
			*/

			/*
				} else if sdiff <= 3 {
					c = c1*Pow10(sdiff) + c2
					f := Len10(c) - DIGITS
					e := e1 - sdiff + f // e might overflow range alreday at this point
					c = roundEven(divPow10(c, f))

					// rounding can in turn overflow the coefficient,
					if c == E_DIGITS {
						c = E_DIGITS_1
						e++
					}

					if e > EXP_MAX_NORM {
						return Dec64{s1 | INF_PATTERN}
					}

					return encodeFinal(s1, e, c)

			*/

		}
	}

	if sdiff < DIGITS {

		hi_c2, lo_c2, sf = divPow10(c2, int(sdiff))
		c = c1 + hi_c2
		if c < E_DIGITS {
			// rounding is based on the low part truncatd from c2
			// check coefficient overrun
			return encodeChkCoefficient(s1, e1, roundEven(c, lo_c2, sf))
		} else {
			// coefficient overflowed; adjust c and round
			// for rounding to be correct for tie-even, the discarded digits must be
			// included in rounding
			// it seems impossible to avoid this extra division operation,
			// but the cases where this is required should be rare
			// check exponent overrun
			// return encodeChkExponent(s1, e1+1, roundEvenTail(c/10, c%10, 10, lo_c2))
			// the solution woth tail seems simpler..
			c, lo_c = c/10, c%10
			lo_c = lo_c*sf + lo_c2
			return encodeChkExponent(s1, e1+1, roundEven(c, lo_c, Pow10(int(sdiff+1))))

		}
	} else if sdiff == DIGITS {
		// second operand only affects result through possible rounding
		// check coefficient overrun
		return encodeChkCoefficient(s1, e1, roundEven(c1, c2, E_DIGITS))
	} else {
		// second operand is so small scale, it will not influence result
		// somehow we should return original value unchanged...
		return encodeFinal(s1, e1, c1)
	}
}

/*
normalsie functions for sub
*/

// following subtraction of equal scale or one off, coefficient can need much scaling
// only scales up
func encodeScaleUp(s uint64, e int, c DecBase) Dec64 {

	debug_assert(c > 0 && c < E_DIGITS)
	if c >= E_DIGITS_1 {
		// if c is in normalised range
		return encodeFinal(s, e, c)
	}
	// need to scale up
	f := DIGITS - Len10(c)
	return encodeChkUnderflow(s, e-f, MulPow10(c, f))
}

func encodeChkUnderflow(s uint64, e int, c DecBase) Dec64 {
	if e < EXP_MIN_NORM {
		return Dec64{s}
	}
	return encodeFinal(s, e, c)
}

// func subtractNormals(a, b Dec64) Dec64 {

// 	_, s1, e1, c1 := decode(a)
// 	_, _, e2, c2 := decode(b)

func subtractNormals(s1 DecBase, e1 int, c1 DecBase, e2 int, c2 DecBase) Dec64 {

	// thy need to have different signs for this to be sub

	var c, sf, hi_c2, lo_c2, lo_c uint64

	// they need even be sorted accoring to scale
	debug_assert(e1 >= e2)
	sdiff := e1 - e2
	if sdiff == 0 {
		// same scale, only significands should be subtracted.
		// Much precision can be lost if c1 and c2 are almost equal.
		// 1234 - 1233 -> 1 (normalised)-> 1.000
		// zero is possible, exp can underflow

		if c1 > c2 {
			return encodeScaleUp(s1, e1, c1-c2)
		}
		if c1 < c2 {
			return encodeScaleUp(s1^SIGN_MASK, e1, c2-c1)
		}

		// zero from Diff has no sign
		return Dec64{}

		// and turn sign if c1 < c2
		return encodeScaleUp(s1^SIGN_MASK, e1, c2-c1)

		/*
			// branch fre...
			c = c1 - c2
			return encodeScaleUp(s1^(c&SIGN_MASK), e1, uabs(int64(c)))
		*/

	} else if sdiff == 1 {

		c = c1*10 - c2
		if c >= E_DIGITS {
			// high range
			// div with 10 and round, exp unchanged.
			// in this case, carry cannot cause overflow - rounding cannot bring c higher than c1
			// eg: 2200 - 999.9 -> 1200.1 (round)-> 1200
			// eg: 9999 - 100.1 -> 9898.9 (round)-> 9899
			return encodeFinal(s1, e1, roundEven(c/10, c%10, 10))
		} else {
			// lower scale
			// eg: 1200 - 999.9 -> 200,1 scale up, no rounding
			// much precision can be lost, worst case being
			// 1000 - 999.9 -> 0.1 (for 4 digit coefficient)
			// scale up if needed
			// zero is not  possible, exp underflow is
			// start exponent (e1-1) does not underflow as e1-1 == e2
			return encodeScaleUp(s1, e1-1, c)
		}
		/*
		   lo_c2 = c2 % 10
		   if lo_c2 == 0 {
		       return encode(s1, e1, c1 - c2 / 10)
		   } else {
		       hi_c := c1 - c2 / 10 - 1
		       lo_c := 10 - lo_c2
		       if hi_c >= E_DIGITS_1 {
		           return encode(s1, e1, hi_c + r_carry(lo_c, 10))
		       } else {
		           // shift digit left
		           return encode(s1, e1 - 1, hi_c * 10 + lo_c)
		       }
		   }
		*/
	} else if sdiff == 1 {

		c = c1*10 - c2
		if c >= E_DIGITS {
			// high range
			// div with 10 and round, exp unchanged.
			// rounding cannot cause overflow here - rounding cannot bring c higher than c1
			// eg: 2200 - 999.9 -> 1200.1 (round)-> 1200
			// eg: 9999 - 100.1 -> 9898.9 (round)-> 9899
			return encodeFinal(s1, e1, roundEven(c/10, c%10, 10))
		} else {
			// lower scale
			// eg: 1200 - 999.9 -> 200,1 scale up, no rounding
			// much precision can be lost, worst case being
			// 1000 - 999.9 -> 0.1 (for 4 digit coefficient)
			// scale up if needed
			// zero is not  possible, exp underflow is
			// start exponent (e1-1) does not underflow as e1-1 == e2
			return encodeScaleUp(s1, e1-1, c)
		}

		/*
			} else if sdiff == 2 {

				c = c1*100 - c2
				if c >= E_DIGITS * 10 {
					// high range
					return encodeFinal(s1, e1, roundEven(c/100, c%100, 100))
				} else {
					// low range
					// coefficient can overflow but not exp
					return encodeChkCoefficient(s1, e1-1, roundEven(c/10, c%10, 10))
				}
		*/
	} else if sdiff == 3 {

		c = c1*1000 - c2
		if c >= E_DIGITS*100 {
			// high range
			return encodeFinal(s1, e1, roundEven(c/1000, c%1000, 1000))
		} else {
			// lower scale
			// coefficient can overflow but not exp
			return encodeChkCoefficient(s1, e1-1, roundEven(c/100, c%100, 100))
		}

	} else if sdiff < DIGITS {
		/*
			hi_c2, lo_c2, sf = divPow10(c2, sdiff-1)

			   // let sf = Pow10(sdiff - 1)
			   // let lo_c2 = c2 % sf

			   c = c1 * 10 - hi_c2

			   // - (lo_c2 != 0) as DecBase;

			   if c >= E_DIGITS {
			   // this rounding ignores after first decimal - not sufficient for tie-even
			       // reduce with 10 and round, exp unchanged
			       // there  will be no coefficient or exponent overflow
			       encodeNormalised(s1, e1, c / 10 + r_carry(c % 10, 10))
			   } else {
			       if lo_c2 == 0 {
			           construct_repr(s1, e1 - 1, c)
			       } else {
			           // this can not exp underflow,
			           // but it may need to round up
			           ext_construct(s1, e1 - 1, c + r_carry(sf - lo_c2, sf))
			       }
			   }

		*/

		// adjust c1 down with sf and c2 up with same, making c2 "positive"
		// convert to addition - this can pull both coeffs below normalised range
		c1 = c1 - Pow10(DIGITS-sdiff)
		c2 = E_DIGITS - c2

		hi_c2, lo_c2, sf = divPow10(c2, sdiff)
		c = c1 + hi_c2
		if c >= E_DIGITS_1 {
			// still within range; no coeff overflow
			return encodeFinal(s1, e1, roundEven(c, lo_c2, sf))
		}

		// use one more digit from c2
		if lo_c2 == 0 {
			// there are no more digits
			return encodeFinal(s1, e1-1, c*10)
		}

		// shift a digit - hopefully rarely
		hi_c2, lo_c2, sf = divPow10(lo_c2, sdiff-1)
		// this can overflow coefficient! but not exponent
		return encodeChkCoefficient(s1, e1-1,
			roundEven(c*10+hi_c2, lo_c2, sf))

		/*
			hi_c2, lo_c2, sf = divPow10(c2, sdiff)
			c = c1 - hi_c2

			if lo_c2 == 0 {
				// no tail
				if c >= E_DIGITS_1 {
					// c in normalised range - likely the common case
					return encodeFinal(s1, e1, c)
				} else {
					// scale up, c cannot be lower than next lower range
					// exponent cannot underflow because e2 is lower than e1-1
					return encodeFinal(s1, e1-1, c*10)
				}
			} else {
				// consider tail
				c--
				lo_c = sf - lo_c2

				if c > E_DIGITS_1 {
					// c will be in range, adjusted for tail
					return encodeFinal(s1, e1, roundEven(c, lo_c, sf))
				} else {
					// shift a digit
					hi_c2, lo_c2, sf = divPow10(lo_c, sdiff-1)
					c = c1*10+hi_c2
					// this can overflow coefficient! but not exponent
					return encodeChkCoefficient(s1, e1-1,
						roundEven(c, lo_c2, sf))


						// shift digit left - expensive?
					//	t := 10 * lo_c
					//	c = c*10 + t/sf
					//	return encode(s1, e1-1, roundEven(c, t%sf, sf))


				}
			}
		*/

	} else if sdiff == DIGITS {
		// subtraction like: 1001 - 0.9999 -> 1000.0001 (round)-> 1000
		// subtraction like: 1001 - 0.1000 -> 1000.9000 (round)-> 1001
		// only if c is low bound can subtrcation go below range
		// subtract: 1000 - 0.9999 -> 999.0001 (round)-> 999.0
		// subtract: 1000 - 0.1000 -> 999.9000 (round)-> 999.9
		c = c1 - 1
		lo_c = E_DIGITS - c2

		if c >= E_DIGITS_1 {
			return encodeFinal(s1, e1, roundEven(c, lo_c, E_DIGITS))
		}
		// shift first digit of lo_c into c
		return encodeFinal(s1, e1-1, c*10+
			roundEven(lo_c/E_DIGITS_1, lo_c%E_DIGITS_1, E_DIGITS_1))

		/*
		   * alt if c1 > minimum, the effect will only be a possible sub of 1

		   if c1 > E_DIGITS_1 { c = c1 - 1 + r_carry(E_DIGITS-c2, E_DIGITS) }
		   else (c1 == E_DIGITS_1)

		*/
	} else if sdiff == DIGITS+1 && c1 == E_DIGITS_1 {
		// no effect on result if c1 higher than low range:
		// 1001 - 0.0xxxx -> 1000.9yyyy (round4)-> 1001
		// where yyyy = 10000 - xxxx
		// 1000 - 0.0xxxx -> 999.9yyyy (round4)-> result depends on yyyy
		// rounding at 4 digits, tie to even:
		// 999.95000 (round4)-> 1000
		// 999.94999 (round4)-> 999.9
		// so, yyyy >= 5000 causes roundup to unmodified c1 value
		// yyyy >= 5000 <=> 10000 - xxxx >= 5000 <=> 5000 >= xxxx
		if c2 <= E_DIGITS_H {
			// return unmodified
			return encodeFinal(s1, e1, c1)
		}
		return encodeFinal(s1, e1-1, E_DIGITS-1)

	} else {
		// second operand doesnt affect result
		// we should return self, if possble
		return encode(s1, e1, c1)
		// return a
	}

}
