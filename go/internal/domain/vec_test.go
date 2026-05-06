package domain

import "testing"

func TestVecOps(t *testing.T) {
	a := Vec3{1, 2, 3}
	b := Vec3{4, 5, 6}
	if got := a.Add(b); got != (Vec3{5, 7, 9}) {
		t.Fatalf("Add: %#v", got)
	}
	if got := b.Sub(a); got != (Vec3{3, 3, 3}) {
		t.Fatalf("Sub: %#v", got)
	}
	if got := a.Scale(2); got != (Vec3{2, 4, 6}) {
		t.Fatalf("Scale: %#v", got)
	}
}
