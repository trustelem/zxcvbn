package match

import "testing"

func TestSort(t *testing.T) {
	matches := []*Match{
		{
			Pattern: "a",
			I:       1,
			J:       3,
		},
		{
			Pattern: "b",
			I:       1,
			J:       3,
		},
		{
			Pattern: "c",
			I:       1,
			J:       1,
		},
		{
			Pattern: "d",
			I:       0,
			J:       4,
		},
		{
			Pattern: "e",
			I:       0,
			J:       5,
		},
		{
			Pattern: "f",
			I:       0,
			J:       1,
		},
	}

	expected := []string{"f", "d", "e", "c", "a", "b"}

	Sort(matches)

	for i := range matches {
		if matches[i].Pattern != expected[i] {
			t.Errorf("matches[%d].Pattern() = %s, want %s", i, matches[i].Pattern, expected[i])
		}
	}
}
