package dec64

import (
	// "fmt"
	"math"
	//	"math/bits"
	"testing"
)

var roundData = []Dec64{
	One(1),
	FromParts(true, 0, 0),
	FromParts(true, 0, 1),
	FromParts(false, -2, 125),
	FromParts(true, -1, 125),
	FromParts(true, -2, 150),
	FromExpInt(-2, 9999),
	FromInt(math.MinInt64),
	Inf(-1),
	FromExpInt(-1, 5),
	FromExpInt(-2, 51),
	FromExpInt(-3, 51),
	FromExpInt(-2, 49),
}

func TestRound(t *testing.T) {

	for _, a := range roundData {

		// show data
		prn("round away:", a, a.Round())
		prn("round even:", a, a.RoundEven())
		prn("round trunc:", a, a.Trunc())
		prn("round ceil:", a, a.Ceil())
		prn("round floor:", a, a.Floor())
	}

}

func TestRoundingPrimitives(t *testing.T) {

	prn("test carry even")
	prn("0, 0, 10 ->", roundEven(0, 0, 10))
	prn("0, 4, 10 ->", roundEven(0, 4, 10))
	prn("0, 5, 10 ->", roundEven(0, 5, 10))
	prn("1, 0, 10 ->", roundEven(1, 0, 10))
	prn("1, 4, 10 ->", roundEven(1, 4, 10))
	prn("1, 5, 10 ->", roundEven(1, 5, 10))
	prn("1, 6, 10 ->", roundEven(1, 6, 10))
	prn("2, 1, 10 ->", roundEven(2, 1, 10))
	prn("2, 4, 10 ->", roundEven(2, 4, 10))
	prn("2, 5, 10 ->", roundEven(2, 5, 10))
	prn("2, 6, 10 ->", roundEven(2, 6, 10))
	prn("3, 5, 10 ->", roundEven(3, 5, 10))

	// we use off factors in Div

	var a, b uint64
	a, b = 5, 3
	prn(a, "/", b, "->", roundEven(a/b, a%b, b))
	a, b = 4, 3
	prn(a, "/", b, "->", roundEven(a/b, a%b, b))
	a, b = 45, 30
	prn(a, "/", b, "->", a/b, a%b, "rounded", roundEven(a/b, a%b, b))
	a, b = 6, 3
	prn(a, "/", b, "->", roundEven(a/b, a%b, b))

	prn("test carry away")
	prn("0,0, 10", roundAway(0, 0, 10))
	prn("0,4, 10", roundAway(0, 4, 10))
	prn("1, 5, 10", roundAway(1, 5, 10))
	prn("2,5, 10", roundAway(2, 5, 10))
	prn("2,6, 10", roundAway(2, 6, 10))

	prn("test carry even")
	prn("0, 0, 10 ->", roundEven(0, 0, 10))
	prn("0, 4, 10 ->", roundEven(0, 4, 10))
	prn("0, 5, 10 ->", roundEven(0, 5, 10))
	prn("1, 0, 10 ->", roundEven(1, 0, 10))
	prn("1, 4, 10 ->", roundEven(1, 4, 10))
	prn("1, 5, 10 ->", roundEven(1, 5, 10))
	prn("1, 6, 10 ->", roundEven(1, 6, 10))
	prn("2, 1, 10 ->", roundEven(2, 1, 10))
	prn("2, 4, 10 ->", roundEven(2, 4, 10))
	prn("2, 5, 10 ->", roundEven(2, 5, 10))
	prn("2, 6, 10 ->", roundEven(2, 6, 10))
	prn("3, 5, 10 ->", roundEven(3, 5, 10))

	prn("test carry away")
	prn("0,0, 10", roundAway(0, 0, 10))
	prn("0,4, 10", roundAway(0, 4, 10))
	prn("1, 5, 10", roundAway(1, 5, 10))
	prn("2,5, 10", roundAway(2, 5, 10))
	prn("2,6, 10", roundAway(2, 6, 10))

	if roundEven(0, 0, 10) != 0 {
		t.Error("Rounding primitive error")
	}
	if roundEven(0, 4, 10) != 0 {
		t.Error("Rounding primitive error")
	}
	if roundEven(0, 5, 10) != 0 {
		t.Error("Rounding primitive error")
	}
	if roundEven(0, 6, 10) != 1 {
		t.Error("Rounding primitive error")
	}
	if roundEven(1, 0, 10) != 1 {
		t.Error("Rounding primitive error")
	}
	if roundEven(1, 4, 10) != 1 {
		t.Error("Rounding primitive error")
	}
	if roundEven(1, 5, 10) != 2 {
		t.Error("Rounding primitive error")
	}
	if roundEven(1, 6, 10) != 2 {
		t.Error("Rounding primitive error")
	}
	if roundEven(2, 1, 10) != 2 {
		t.Error("Rounding primitive error")
	}
	if roundEven(2, 4, 10) != 2 {
		t.Error("Rounding primitive error")
	}
	if roundEven(2, 5, 10) != 2 {
		t.Error("Rounding primitive error")
	}
	if roundEven(2, 6, 10) != 3 {
		t.Error("Rounding primitive error")
	}
	if roundEven(3, 5, 10) != 4 {
		t.Error("Rounding primitive error")
	}
	if roundEven(717, 5, 10) != 718 {
		t.Error("Rounding primitive error")
	}
	if roundEven(718, 5, 10) != 718 {
		t.Error("Rounding primitive error")
	}
	if roundEven(717, 49, 100) != 717 {
		t.Error("Rounding primitive error")
	}
	if roundEven(717, 50, 100) != 718 {
		t.Error("Rounding primitive error")
	}
	if roundEven(718, 50, 100) != 718 {
		t.Error("Rounding primitive error")
	}

}

func TestRoundingPrimitives2(t *testing.T) {

	for v := range 16 {
		for _, d := range testRoundData {
			c := roundEven(divPow10(d, v))

			prn("round", d, v, c, c)
			c2 := roundEvenX(divPow10(d, v))

			if c != c2 {
			
				t.Error("Rounding primitive error", d, v, c, c2)
				
			}
		}
	}
}


func TestRoundFuncs(t *testing.T) {

	testRound := func(rf func(Dec64) Dec64, x, r Dec64) {
	y := rf(x)
	if !y.Equal(r) {
	
	t.Error("Rounding failed", x, y, r)
	
	}
	}

	testRound(Dec64.Round, FromInt(9999_9999_9999_9999),FromInt(9999_9999_9999_9999))
	testRound(Dec64.Round, FromInt(9999_9999_9999_999),FromInt(9999_9999_9999_999))
	testRound(Dec64.Round, FromExpInt(-1, 9999_9999_9999_9999),FromInt(1000_0000_0000_0000))
	testRound(Dec64.Round, FromExpInt(-8, 9999_9999_9999_9999),FromInt(1000_0000_0))
	testRound(Dec64.Round, FromExpInt(-15, 9999_9999_9999_9999),FromInt(10))
	testRound(Dec64.Round, FromExpInt(-16, 9999_9999_9999_9999),FromInt(1))
	testRound(Dec64.Round, FromExpInt(-17, 9999_9999_9999_9999),FromInt(0))
	testRound(Dec64.Round, FromExpInt(-2, 1234_5678_9012_3456),FromInt(1234_5678_9012_35))
	testRound(Dec64.Round, FromExpInt(-3, 1234_5678_9012_3456),FromInt(1234_5678_9012_3))
	testRound(Dec64.Round, FromExpInt(-3, 1234_5678_0000_0000),FromInt(1234_5678_0000_0))
	
}


//
// --------------------------------------------------------------------------
// benchmark the rounding functions
// -------------------------------------------------------------------------

var testRoundData = []uint64{
	9999_9999_9999_9999,
	1234_5678_9012_3456,
	1000_0000_0000_0000,
	1234_5000_0000_0000,
	1235_0000_0000_0000,
}

var round_suma uint64

func BenchmarkRoundEven(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for v := range 16 {
			for _, d := range testRoundData {
				c := roundEven(divPow10(d, v))

				suma += c
			}

		}
	}
}

func BenchmarkRoundEvenX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for v := range 16 {
			for _, d := range testRoundData {
				c := roundEvenX(divPow10(d, v))

				suma += c
			}

		}
	}
}

// despite branch free, this measures slower
// returns rounded coefficient, tie to even
func roundEvenX(c DecBase, r, base DecBase) uint64 {
	// rounding: tie-even
	// unclear how to dothis efeficiently and properly

	if r == 0 {
		return c
	}

	r2 := r + r

	if base > r2 {
		return c
	}

	if base < r2 {
		return c + 1
	}
	return c + (c & 1)

	/*

		x := int64(base - r - r)

		if x > 0 {
			return c
		}

		if x < 0 {
			return c + 1
		}
		return c + (c & 1)
	*/

	// seems we could do it with bit operations
	/*
		// x is unsigned
		x := base - (r << 1)
		// if x > 0 -> 0  (signbit of x)
		// if x == 0 -> (1&c)
		// if x < 0 -> 1  (sign bit of x)

		// x  == 0 : abs(x) - 1
		// abs(x) : y = x s>> 63, (x ^ Y) - y
		//y := uint64(int64(x) >> 63)
		// abs := (x ^ y) - y
		//	return c + ((((abs - 1) >> 63) & (c&1)) | (x >> 63))
		//		return c + (((abs - 1)  & (c << 63)) | x ) >>63
		// alternative x == 0: ^(x | -x) (simpler?)

		// havecarry := (x == 0 && c&1 != 0) || x < 0
		// THOUGHT: the no carry predicate may be cheaper to check
		// havenocarry := x > 0 || (x == 0 && c&1 == 0)

		// bitwise calculation of havecarry
		// from hackers delight we have
		// x == 0: ^(x | -x)
		// calculation is in first bit, and then shifted down
		carry := ((^(x | -x) & (c << 63)) | x) >> 63
		return c + carry
	*/

	//	x := int64(r + r - base)

	/*
		if x > 0  || (x == 0 && c&1 != 0) {
		return c+1
		}
		return c
	*/
	/*
	 if x < 0 || (x == 0 && c&1 == 0) {
	 return c
	 }
	 return c + 1
	*/
	/*
		if x > 0 || (x == 0 && c&1 == 0)  { return c }
		return c+1
	*/

	/*
			r2 := r + r
		if base > r2 || (base == r2 && c&1 == 0)  { return c }
		return c+1
	*/

}
