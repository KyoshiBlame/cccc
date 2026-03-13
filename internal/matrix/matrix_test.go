package matrix

import "testing"

func TestTrace(t *testing.T) {
	m := Matrix{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	v, err := Trace(m)
	if err != nil {
		t.Fatal(err)
	}
	if v != 15 {
		t.Fatalf("want 15 got %v", v)
	}
}

func TestDeterminant(t *testing.T) {
	m := Matrix{{1, 2}, {3, 4}}
	v, err := Determinant(m)
	if err != nil {
		t.Fatal(err)
	}
	if v != -2 {
		t.Fatalf("want -2 got %v", v)
	}
}

func TestExpressionOps(t *testing.T) {
	a := Matrix{{1, 0}, {0, 1}}
	b := Matrix{{2, 0}, {0, 2}}
	c, err := Add(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if c[0][0] != 3 || c[1][1] != 3 {
		t.Fatalf("unexpected result: %#v", c)
	}
}
