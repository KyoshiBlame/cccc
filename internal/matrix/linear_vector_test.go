package matrix

import (
	"math"
	"testing"
)

func TestSolveGauss(t *testing.T) {
	A := Matrix{{2, 1}, {5, 7}}
	b := []float64{11, 13}
	x, _, err := SolveGauss(A, b)
	if err != nil {
		t.Fatal(err)
	}
	if len(x) != 2 || math.Abs(x[0]-64.0/9.0) > 1e-8 || math.Abs(x[1]+29.0/9.0) > 1e-8 {
		t.Fatalf("unexpected x: %#v", x)
	}
}

func TestVectorCross(t *testing.T) {
	a := []float64{1, 0, 0}
	b := []float64{0, 1, 0}
	v, err := Cross(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if v[0] != 0 || v[1] != 0 || v[2] != 1 {
		t.Fatalf("unexpected cross result: %#v", v)
	}
}
