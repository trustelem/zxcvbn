package matching

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trustelem/zxcvbn/adjacency"
	"github.com/trustelem/zxcvbn/match"
)

func Test_spatialMatch(t *testing.T) {
	s := spatialMatch{
		graphs: defaultGraphs,
	}
	// doesn't match 1- and 2-character spatial patterns
	assert.Empty(t, s.Matches(""))
	assert.Empty(t, s.Matches("/"))
	assert.Empty(t, s.Matches("qw"))
	assert.Empty(t, s.Matches("*/"))

	// for testing, make a subgraph that contains a single keyboard
	s = spatialMatch{
		graphs: []*adjacency.Graph{adjacency.Graphs["qwerty"]},
	}

	pattern := "6tfGHJ"
	assert.Equal(t, []*match.Match{
		{
			Pattern:      "spatial",
			Token:        pattern,
			I:            03,
			J:            3 + len(pattern) - 1,
			Graph:        "qwerty",
			Turns:        2,
			ShiftedCount: 3,
		},
	}, s.Matches("rz!"+pattern+"%z"))

	tests := []struct {
		pattern  string
		keyboard string
		turns    int
		shifts   int
	}{
		{"12345", "qwerty", 1, 0},
		{"@WSX", "qwerty", 1, 4},
		{"6tfGHJ", "qwerty", 2, 3},
		{"hGFd", "qwerty", 1, 2},
		{"/;p09876yhn", "qwerty", 3, 0},
		{"Xdr%", "qwerty", 1, 2},
		{"159-", "keypad", 1, 0},
		{"*84", "keypad", 1, 0},
		{"/8520", "keypad", 1, 0},
		{"369", "keypad", 1, 0},
		{"/963.", "mac_keypad", 1, 0},
		{"*-632.0214", "mac_keypad", 9, 0},
		{"aoEP%yIxkjq:", "dvorak", 4, 5},
		{";qoaOQ:Aoq;a", "dvorak", 11, 4},
	}
	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			s := spatialMatch{
				graphs: []*adjacency.Graph{adjacency.Graphs[tt.keyboard]},
			}
			matches := s.Matches(tt.pattern)
			assert.Equal(t,
				[]*match.Match{
					{
						Pattern:      "spatial",
						Token:        tt.pattern,
						I:            0,
						J:            len(tt.pattern) - 1,
						Graph:        tt.keyboard,
						Turns:        tt.turns,
						ShiftedCount: tt.shifts,
					},
				},
				matches,
			)
		})
	}
}
