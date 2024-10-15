package dec64

import (
	//	"fmt"
	//	"math/bits"
	"testing"
)

// need to benchmark some functions

var testDivPowData = []uint64{
	9999_9999_9999_9999,
	1234_5678_9012_3456,
	1000_0000_0000_0000,
}

func TestDivPow10(t *testing.T) {
	d := testDivPowData[1]
	for v := range 16 {
		ay, by, fy := divPow10(d, v)
		ax, bx, fx := trivialDivPow10(d, v)

		if ay != ax || by != bx || fy != fx {
			t.Error("Divpow10 fail", v, ay, by, "!=", ax, bx)
		}
	}
}

// TtrivialLen64 returns the number of decimal digits required to represent x;
// the result is 0 for x == 0.
// trivial implementation with no consideration for performance
func TrivialLen64(x uint64) int {
	c := 0
	for x != 0 {
		c++
		x = x / 10
	}
	return c
}

func trivialDivPow10(c uint64, e int) (uint64, uint64, uint64) {

	f := LBASE_POWERS[e]
	return c / f, c % f, f
}

// --------------------------------------------------------------------------
// benchmark
// -------------------------------------------------------------------------

var suma uint64

// it seems to make sense with this function
func BenchmarkDivPow10(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for v := range 16 {
			for _, d := range testDivPowData {
				a, b, _ := divPow10(d, v)

				suma += a
				suma += b

			}
		}
	}
}

// this was aboyt 4% fatser -thats nothing
// could it be faster by skipping nexted func call
func BenchmarkDivPow10X(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for v := range 16 {
			for _, d := range testDivPowData {
				a, b, _ := divPow10X(d, v)
				// a, b, _ := FFUNCS[v](d)

				suma += a
				suma += b

			}
		}
	}
}

// this was aboyt 4% fatser -thats nothing
// could it be faster by skipping nexted func call
func BenchmarkTrivialDivPow10(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for v := range 16 {
			for _, d := range testDivPowData {
				a, b, _ := trivialDivPow10(d, v)

				suma += a
				suma += b

			}
		}
	}
}

// same as trivial divpow
func BenchmarkDivPow10_2(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for v := range 16 {
			for _, d := range testDivPowData {

				f := Pow10(v)
				a, b, _ := d/f, d%f, f

				suma += a
				suma += b

			}
		}
	}
}

// predefined constant - it is fatser than calling divpow10
func BenchmarkDivPow10_3(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for v := range 16 {
			for _, d := range testDivPowData {
				_ = v
				const f = e11
				a, b, _ := d/f, d%f, f

				suma += a
				suma += b

			}
		}
	}
}

func BenchmarkMulPow10(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for v := range 16 {
			for _, d := range testDivPowData {

				f := MulPow10(d, v)

				suma += f

			}
		}
	}
}
