package hist

import "testing"

func TestBins(t *testing.T) {
	got := NewBins(0, 10, 5)
	expected := Bins{0, 2, 4, 6, 8, 10}
	t.Logf("sample for 5 bins on range [0, 10)")

	for i := range got {
		if got[i] != expected[i] {
			t.Fatalf("for generated bins: got %+v, expected %+v", got, expected)
		}
	}

	val := 0.0
	exp := 0
	if b := got.Bin(val); b != exp {
		t.Errorf("value %v binned to %v - expected bin %v", val, b, exp)
	}

	val = 3.0
	exp = 1
	if b := got.Bin(val); b != exp {
		t.Errorf("value %v binned to %v - expected bin %v", val, b, exp)
	}

	val = 9.99999999
	exp = 4
	if b := got.Bin(val); b != exp {
		t.Errorf("value %v binned to %v - expected bin %v", val, b, exp)
	}

	val = 42
	exp = 5
	if b := got.Bin(val); b != exp {
		t.Errorf("value %v binned to %v - expected bin %v", val, b, exp)
	}
}
