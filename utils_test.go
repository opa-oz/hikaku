package hikaku

import "testing"

func TestMin(t *testing.T) {
	got := min(4, 6)
	want := 4

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestMax(t *testing.T) {
	got := max(4, 6)
	want := 6

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestCheckNegative(t *testing.T) {
	got := check(4, 6, 3, 10)

	if got != false {
		t.Errorf("got %t, wanted %t", got, false)
	}
}

func TestCheck(t *testing.T) {
	got := check(4, 6, 10, 10)

	if got != true {
		t.Errorf("got %t, wanted %t", got, true)
	}
}
