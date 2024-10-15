package dec64

import (
	"fmt"
	"math"
	"math/bits"
	"strconv"
)

// Floating point classification
const (
	fpZero = iota
	fpNormal
	fpInf
	fpNaN
)

const (
	signmask    = 1 << 63
	expbits     = 11
	mantbits    = 52
	expmask     = 0x7FF << mantbits
	mantimplbit = 1 << mantbits
	mantmask    = mantimplbit - 1
	bias        = 1023
	emax        = bias
	emin        = 1 - emax
)

// returns class, sign for floating point value
func fpClassify(f float64) (int, uint64) {

	bf := math.Float64bits(f)
	sign := bf & signmask
	exp := bf & expmask
	if exp == 0 {
		if bf&mantmask == 0 {
			return fpZero, sign
		}
		// this is subnormal, but they are handled the same as normal
		return fpNormal, sign
	}
	if exp == expmask {
		// this is inf or nan
		if bf&mantmask == 0 {
			return fpInf, sign
		}
		return fpNaN, sign
	}
	return fpNormal, sign
}

// returns class, sign, exponent, mantissa
func fpDecode(f float64) (int, uint64, int, uint64) {

	bf := math.Float64bits(f)
	sign := bf & signmask
	exp := bf & expmask
	mant := bf & mantmask
	if exp == 0 {
		if mant == 0 {
			return fpZero, sign, 0, 0
		}
		// subnormal value
		// could normalise mantissa and reduce exp
		return fpNormal, sign, emin - mantbits, mant
	}
	if exp == expmask {

		// this is inf or nan
		if mant == 0 {
			return fpInf, sign, 0, 0
		}
		return fpNaN, sign, 0, mant
	}
	return fpNormal, sign, int(exp>>mantbits) - bias - mantbits, mantimplbit | mant
}

// it is intended to convert strings in teh form
// strconv.FormatFloat(math.Abs(d), 'g', -1, 64), only!
// And only finite numbers (Inf and nan not handled).
// only for abs values
// this is not intended for parsing user input!

func fpConvertString(s string) (int, uint64) {
	/*
	   leading sign will not be provided in input string
	   ddd.dddde[+|-]eddd

	   this processing does not check syntax is correct,
	   as string is output from strconv it is assumed to be correct
	*/

	var negexp bool
	var m uint64
	var ix, tenexp, exp int
	var d, c byte

	for ; ix < len(s); ix++ {
		c = s[ix]
		d = c - '0'
		if d <= 9 {
			// digit
			m = m*10 + uint64(d)

		} else if c == '.' {
			// start counting decimals
			ix++
			goto DECIMALS
		} else if c == 'e' {
			// end of coefficient
			ix++
			goto EXPONENT
		}
	}

	// there was no decimals, exponent
	return 0, m

DECIMALS:
	for ; ix < len(s); ix++ {
		c = s[ix]
		d = c - '0'
		if d <= 9 {
			// digit
			m = m*10 + uint64(d)
			tenexp--
		} else if c == 'e' {
			// end of coefficient
			ix++
			goto EXPONENT
		}
	}

	// there was no exponent
	return tenexp, m

EXPONENT:
	// we have had the 'e'
	// exponent has an optional sign, and one or more digits.
	// it seems the exponent always has sign in outbput from FormatFloat e/g
	// [+|-] dddd

	// sign -optional, but it's always there

	if ix < len(s) {
		c = s[ix]
		if c == '-' {
			ix++
			negexp = true
		} else if c == '+' {
			ix++
		}
	}

	for ; ix < len(s); ix++ {
		c = s[ix]
		d = c - '0'
		if d <= 9 {
			// digit
			exp = exp*10 + int(d)
		}
	}

	if negexp {
		return tenexp - exp, m
	}
	return tenexp + exp, m
}

// frankly, not many values can be converted simply
func try_convert() {

	for _, d := range td {
		cls, sgn, exp, mant := fpDecode(d)

		fmt.Println("float: ", d)
		fmt.Println("float classify", sgn, exp, mant, cls)

		// what can we do to transform number - only for normals
		if cls == decNormal {
			tz := bits.TrailingZeros64(mant)
			imant := mant >> tz
			// leading zeros: original 11 + tz
			lz := 11 + tz
			// adjust exp with tz
			exp = exp + tz
			fmt.Println("float adjusted", exp, imant)

			if exp >= 0 {
				// number is int, but can it be contained in uint64
				if exp <= lz {
					number := imant << exp
					fmt.Println("number was an int: ", number)
				} else {
					// how get a ten exp out of this
					fmt.Println("number was an int, but too big: mant, exp", imant, exp)
				}
			} else {
				// number is fraction, exp is smaller than zero
				fmt.Println("number is a fraction", imant, exp)
				// if there is room to adjust
				tenexp := 0
				for exp < 0 {
					imant = imant * 5
					tenexp--
					exp++
					fmt.Println("   adjusting", imant, exp, tenexp)
				}
			}
		}
		fmt.Println()

	}
}

func FromFloat(f float64) Dec64 {

	cls, sgn := fpClassify(f)
	switch cls {

	case fpZero:
		return Dec64{sgn}

	case fpNormal:
		s := strconv.FormatFloat(math.Abs(f), 'g', -1, 64)
		exp, m := fpConvertString(s)
		return encodeAny(sgn, exp, m)

	case fpInf:
		return Dec64{sgn | INF_PATTERN}

	default:
		return Dec64{NAN_PATTERN}
	}
}
