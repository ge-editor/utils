package utils

import (
	"reflect"
	"testing"
)

func TestMoveToFront(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		i    int
		want []string
	}{
		{"middle", []string{"a", "b", "c", "d"}, 2, []string{"c", "a", "b", "d"}},
		{"last", []string{"a", "b", "c", "d"}, 3, []string{"d", "a", "b", "c"}},
		{"second", []string{"a", "b", "c"}, 1, []string{"b", "a", "c"}},
		{"first", []string{"a", "b", "c"}, 0, []string{"a", "b", "c"}},
		{"negative", []string{"a", "b", "c"}, -1, []string{"a", "b", "c"}},
		{"out of range", []string{"a", "b", "c"}, 3, []string{"a", "b", "c"}},
		{"single", []string{"a"}, 0, []string{"a"}},
		{"empty", []string{}, 0, []string{}},
		{"nil", nil, 1, nil},
	}
	for _, c := range cases {
		var s []string
		if c.in != nil {
			s = make([]string, len(c.in))
			copy(s, c.in)
		}
		MoveToFront(s, c.i)
		if !reflect.DeepEqual(s, c.want) {
			t.Errorf("%s: MoveToFront(%v, %d) = %v, want %v", c.name, c.in, c.i, s, c.want)
		}
	}
}

// The order of use: a buffer that is shown again goes to the front and the
// others keep their order.
func TestMoveToFrontKeepsOrderOfUse(t *testing.T) {
	order := []string{"A", "B"} // A is shown, B was opened before and not used
	use := func(name string) {
		for i, v := range order {
			if v == name {
				MoveToFront(order, i)
				return
			}
		}
		order = append(order, name) // a new one is added at the end ...
		MoveToFront(order, len(order)-1) // ... and moves to the front when shown
	}

	use("C") // open C while A is shown
	if !reflect.DeepEqual(order, []string{"C", "A", "B"}) {
		t.Fatalf("after opening C: %v", order)
	}
	order = order[1:] // close C: the buffer to go back to is the first one
	if order[0] != "A" {
		t.Fatalf("closing C must go back to A, got %v", order)
	}

	use("B")
	use("C")
	order = order[1:] // close C again
	if order[0] != "B" {
		t.Fatalf("closing C must go back to B (used before C), got %v", order)
	}
}
