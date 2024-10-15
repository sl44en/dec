package dec64

// import "fmt"

// String representation of Decimal value
func (d Dec64) String() string {
	// TODO: this implemenatation is obviously not optimised

	t, s, e, v := decode(d)

	switch t {
	case decZero:
		return ternaryIf(s != 0, "-0", "0")
	case decNormal:
		goto NORMAL
	case decInf:
		return ternaryIf(s != 0, "-Inf", "+Inf")
	case decNan:
		return "NaN"
	}

NORMAL:

	ix := 0
	buf := [P]byte{}

	// split in int and fraction

	sr := ternaryIf(s != 0, "-", "")

	if e <= -P {
		//  value is below 1
		sr += "0."
		for i := P + e; i < 0; i++ {
			sr += "0"
		}

		fraction := v
		// fraction to digits
		const divisor = E_DIGITS_1 // 1e15
		for ix = 0; fraction != 0; ix++ {
			buf[ix] = byte('0' + fraction/divisor)
			fraction = (fraction % divisor) * 10
		}
		sr += string(buf[:ix])

	} else if e < 0 {
		// it will have both int and fraction parts
		p := Pow10(-e)

		// int part
		intpart := v / p
		ix = len(buf)

		// this loop will print all digits in intpart, incl potential trailing zeros
		for intpart != 0 {
			ix--
			buf[ix] = byte('0' + intpart%10)
			intpart /= 10
		}
		sr += string(buf[ix:])

		fraction := v % p

		if fraction != 0 {
			sr += "."

			ix = 0
			divi := p
			for fraction != 0 {
				fraction *= 10
				buf[ix] = byte('0' + fraction/divi)
				fraction %= divi
				ix++
			}
			sr += string(buf[:ix])
		}
	} else {
		// value >= 1e15, 16 digits or more
		ix = len(buf)

		// this loop will print all 16 digits, incl any trailing zeros
		for v != 0 {
			ix--
			buf[ix] = byte('0' + v%10)
			v /= 10
		}
		sr += string(buf[ix:])

		// is this right??
		// add trailing zeros
		for i := e; i > 0; i-- {
			sr += "0"
		}
	}
	return sr
}
