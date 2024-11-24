package dec64

import (
	"fmt"
	"math/bits"
)

// for unused
var Printfn = fmt.Println

// decimal integer math

const (
	e0  = 1
	e1  = 10
	e2  = 100
	e3  = 1000
	e4  = 10000
	e5  = 100000
	e6  = 1000000
	e7  = 10000000
	e8  = 100000000
	e9  = 1000000000
	e10 = 10000000000
	e11 = 100000000000
	e12 = 1000000000000
	e13 = 10000000000000
	e14 = 100000000000000
	e15 = 1000000000000000
	e16 = 10000000000000000
	e17 = 100000000000000000
	e18 = 1000000000000000000
	e19 = 10000000000000000000
)

var BASE_POWERS = [...]uint32{
	1, 10, 100, 1000, e4, e5, e6, e7, e8, e9,
}

var LBASE_POWERS = [...]uint64{
	e0, e1, e2, e3, e4, e5, e6, e7, e8, e9, e10,
	e11, e12, e13, e14, e15, e16, e17, e18, e19,
}

// are they needed
func IntLog2_32(x uint32) int {
	return 31 - bits.LeadingZeros32(x|1)
}

func IntLog2_64(x uint64) int {
	return 63 - bits.LeadingZeros64(x|1)
}

// Len10 returns the number of decimal digits required to represent x;
// the result is 0 for x == 0.
func Len10(x uint64) int {
	dc := bits.Len64(x) * 3 / 10
	if x >= LBASE_POWERS[dc] {
		return dc + 1
	}
	return dc
}

// the assumption is: len10(2 ^ p) == len2(2 ^ p) * 3 / 10
// len2 * 3 / 10 is an aproximation of log10. (1 / log2810) ¨= 0.301)
// where ^ denotes raise to power
// in the interval [2^p : 2^(p+1)[ ( excluding hi end) all numbers will have same binary length, len2
// len10 in lowest end of interval can be determined as shown above: len2 * 3 / 10

// *****************************************

// return 10 raised to power i for non negative i
func Pow10(i int) uint64 {
	return LBASE_POWERS[i]
}

// multiply v by nonnegative power of 10
// not defined for negative i, or i causing overflow
// multiply by 10 power modulo 16. We dont need 10 powers higher than 15 here
func MulPow10(v uint64, i int) uint64 {
	return v * LBASE_POWERS[i]
}

// NOTE: only handling powers [0 - 15]

// can we dop this fats, useful ?
// return x/p , c%p, p  where p = 10^e
// divide by 10 power modulo 16. We dont need 10 powers higher than 15 here
func divPow10X(c uint64, e int) (uint64, uint64, uint64) {

	// utilise division by constants
	// TODO benchmark
	switch uint(e) % 16 {
	case 0:
		return c, 0, 1
	case 1:
		return c / e1, c % e1, e1
	case 2:
		return c / e2, c % e2, e2
	case 3:
		return c / e3, c % e3, e3
	case 4:
		return c / e4, c % e4, e4
	case 5:
		return c / e5, c % e5, e5
	case 6:
		return c / e6, c % e6, e6
	case 7:
		return c / e7, c % e7, e7
	case 8:
		return c / e8, c % e8, e8
	case 9:
		return c / e9, c % e9, e9
	case 10:
		return c / e10, c % e10, e10
	case 11:
		return c / e11, c % e11, e11
	case 12:
		return c / e12, c % e12, e12
	case 13:
		return c / e13, c % e13, e13
	case 14:
		return c / e14, c % e14, e14
	case 15:
		return c / e15, c % e15, e15
	}

	panic("no not here")
	// return 0,0,0
}

// utilise division by constants
// benchmark (amd64) shows 113 ns for this, 130 ns for the above simpler version

// divide  10 power modulo 16. We dont need 10 powers higher than 15 here
func divPow10(c uint64, e int) (uint64, uint64, uint64) {
	return FFUNCS[uint(e)](c)
}

type ffunc func(uint64) (uint64, uint64, uint64)

var FFUNCS = [...]ffunc{
	func(c uint64) (uint64, uint64, uint64) {
		return c, 0, 1
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e1, c % e1, e1
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e2, c % e2, e2
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e3, c % e3, e3
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e4, c % e4, e4
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e5, c % e5, e5
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e6, c % e6, e6
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e7, c % e7, e7
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e8, c % e8, e8
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e9, c % e9, e9
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e10, c % e10, e10
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e11, c % e11, e11
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e12, c % e12, e12
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e13, c % e13, e13
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e14, c % e14, e14
	},
	func(c uint64) (uint64, uint64, uint64) {
		return c / e15, c % e15, e15
	},
}

// this is fastest by current benchmarking
func TrailingZeros64(n uint64) int {

	if n == 0 {
		return 20
	}

	// if n is not even - but is this not just adding time for a rare case

	if (n & 1) != 0 {
		return 0
	}

	//var c uint64
	//	var zeros int

	if (n % e12) == 0 {

		if (n % e16) == 0 {

			if (n % e18) == 0 {
				if (n % e19) == 0 {
					// if n == 0 { return 20 }
					return 19
				}
				return 18

			}
			if (n % e17) == 0 {
				return 17
			}
			return 16
		}

		if (n % e14) == 0 {
			if (n % e15) == 0 {
				return 15
			}
			return 14
		}
		if (n % e13) == 0 {
			return 13
		}
		return 12
	}

	if (n % 1e8) == 0 {

		if (n % e10) == 0 {
			if (n % e11) == 0 {
				return 11
			}
			return 10
		}
		if (n % e9) == 0 {
			return 9
		}
		return 8
	}

	if (n % e4) == 0 {
		if (n % e6) == 0 {
			if (n % e7) == 0 {
				return 7
			}
			return 6
		}
		if (n % e5) == 0 {
			return 5
		}
		return 4
	}

	if (n % e2) == 0 {
		if (n % e3) == 0 {
			return 3
		}
		return 2
	}
	if (n % 10) == 0 {
		return 1
	}
	return 0
}

// if we know it is less 16 and n not zero...
// how to name it...
func TrailingZeros64_16(n uint64) int {

	// if n is not even - but is this not just adding time for a rare case

	if (n & 1) != 0 {
		return 0
	}

	if (n % 1e8) == 0 {
		if (n % e12) == 0 {

			if (n % e14) == 0 {
				if (n % e15) == 0 {
					return 15
				}
				return 14
			}
			if (n % e13) == 0 {
				return 13
			}
			return 12
		}

		if (n % e10) == 0 {
			if (n % e11) == 0 {
				return 11
			}
			return 10
		}
		if (n % e9) == 0 {
			return 9
		}
		return 8
	}

	if (n % e4) == 0 {
		if (n % e6) == 0 {
			if (n % e7) == 0 {
				return 7
			}
			return 6
		}
		if (n % e5) == 0 {
			return 5
		}
		return 4
	}

	if (n % e2) == 0 {
		if (n % e3) == 0 {
			return 3
		}
		return 2
	}
	if (n % 10) == 0 {
		return 1
	}
	return 0
}

// trying to be clever with using number of trailing 0 bits
// this can for some values stop iteration earlier, result in slightly faster average times
func AltTrailingZeros64(v uint64) int {

	if v == 0 {
		return 20
	}

	// devisible by 10 means devisible by 5 and 2
	// divisible by 2 can be determined by trailing zero bits
	// so result is <= trailning zero bits
	tz := bits.TrailingZeros64(v)
	v = v >> tz
	for tzd := 0; tzd < tz; tzd++ {
		if v%5 != 0 {
			return tzd
		}
		v /= 5

		//		fmt.Println("how many loops v maxres tz", v, tz, tzd)

	}
	return tz
}
