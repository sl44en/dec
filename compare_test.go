package dec64

import (
	// "fmt"
	"math"
	//	"math/bits"
	"testing"
)

var compareData = []Dec64{
	One(1),
	FromParts(true, 0, 0),
	FromParts(true, 0, 1),
	FromInt(9999),
	FromInt(math.MinInt64),
	Inf(-1),
	NaN(),
	FromInt(-9999999999999999),
	FromInt(-99999999999999990),
	FromInt(-99999999999999995),
	FromExpUint(0, 125),
	FromExpUint(0, 126),
	FromExpUint(-1, 125),
	FromParts(true, -2, 125),
	FromExpUint(-3, 125),
	FromExpUint(-4, 125),

	FromExpUint(-1, 1),
	FromExpUint(-2, 1),
	FromExpUint(-18, 12345678901234567),
	FromExpUint(-3, 10002),
	FromExpUint(-1, 1234567890123456),
	FromExpUint(-15, 1234567890123456),
	FromParts(false, -1, 9988827272800001),
	FromExpUint(-2, 9988827272800001),

	// large number
	FromExpUint(10, 123456789012),
	FromExpUint(0, 1234567890123456),
	FromParts(true, 0, 12345678901234567),
	FromExpUint(25, 12),
}

func TestCompare(t *testing.T) {

	// this test only works where rounding is not involved.. and maybe even issues with negative zeros

	minf := Inf(-1)
	inf := Inf(1)
	nan := NaN()
	zero := Zero()

	if nan.Equal(nan) {
		t.Error("nan equal nan")
	}

	if nan.Less(nan) {
		t.Error("nan less nan")
	}
	if FromExpInt(-1, 125).Equal(FromExpInt(-2, 125)) {
		t.Error("exp smaller")
	}
	if !FromExpInt(-1, 125).Equal(FromExpInt(-1, 125)) {
		t.Error("exp smaller")
	}
	if !FromExpInt(-1, -125).Equal(FromExpInt(-1, -125)) {
		t.Error("exp smaller")
	}

	if inf.Less(nan) {
		t.Error("inf less nan")
	}

	if nan.Less(inf) {
		t.Error("nan less inf")
	}

	if zero.Less(nan) {
		t.Error("inf less nan")
	}

	if nan.Less(zero) {
		t.Error("nan less inf")
	}

	if !minf.Less(inf) {
		t.Error("!-inf less inf")
	}

	if !minf.Less(MaxValue()) {
		t.Error("!-inf less inf")
	}
	if inf.Less(MaxValue()) {
		t.Error("inf less max")
	}

	if !MaxValue().Less(inf) {
		t.Error("max less inf")
	}

	if FromExpInt(-2, 125).Less(FromExpInt(-2, 125)) {
		t.Error("same is less")
	}

	if FromExpInt(-2, 125).Less(FromExpInt(-2, -125)) {
		t.Error("pos less neg")
	}

	if FromExpInt(-1, 125).Less(FromExpInt(-2, 125)) {
		t.Error("exp smaller")
	}

	if !FromExpInt(-1, -125).Less(FromExpInt(-2, -125)) {
		t.Error("exp smaller")
	}

	if zero.Less(FromExpInt(-2, -125)) {
		t.Error("exp smaller")
	}

	if !FromExpInt(-2, -125).Less(FromExpInt(-2, 125)) {
		t.Error("!neg less pos")
	}

	/*
		for _, a := range compareData {
			for _, b := range compareData {

				d1 := a.Equal(b)
				d2 := a.Less(b)

				prn("a, b", a, b, "=", d1, "<", d2)

				//			if d1 == d2 {
				//				t.Error("add sub mismatch", a, ",", b, d1, d2)
				//			}

			}
		}
	*/

}

// some test vals
var testvals = []uint64{
	1, 2, 0, INF_PATTERN, NAN_PATTERN,
}

// sime signs
var signs = []uint64{
	0, SIGN_MASK, 0,
}

func TestTransferSign(t *testing.T) {

	for _, v := range testvals {

		for _, s := range signs {
			a := transferSign(s, v)
			b := stupidTransferSign(s, v)

			prn("result of sign transfer", s, v, a, b)

			if a != b {
				t.Error("Error in sign transfer", s, v, a, b)
			}

		}
	}

	prn("print sum from bench", sumit)
}

var sumit int64

func updatesum(a int64) {
	sumit += a

}

// it says there is no difference

func BenchmarkTransferSign(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for _, v := range testvals {
			for _, s := range signs {
				a := transferSign(s, v)

				updatesum(a)

			}
		}
	}
}

func BenchmarkStupidTransferSign(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for _, v := range testvals {
			for _, s := range signs {
				a := stupidTransferSign(s, v)

				updatesum(a)

			}
		}
	}
}

func BenchmarkTernarySign(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for _, v := range testvals {
			for _, s := range signs {

				_ = v

				a := ternaryIf(s == 0, int64(INF_PATTERN), -int64(INF_PATTERN))

				updatesum(a)

			}
		}
	}
}

// should i benchmark these
/*
func transferSign(s, tempr uint64) int64 {

	t := int64(s) >> 63
	return (int64(tempr) ^ t) - t
}
*/

func stupidTransferSign(s, tempr uint64) int64 {

	if s == 0 {
		return int64(tempr)
	}
	return int64(-tempr)

}
