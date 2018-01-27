package matching

import (
	"testing"

	"github.com/test-go/testify/assert"
	"github.com/trustelem/zxcvbn/match"
)

func Test_sequenceMatch_Matches(t *testing.T) {
	s := sequenceMatch{}

	// doesn't match 1- and 2-character spatial patterns
	assert.Empty(t, s.Matches(""))
	assert.Empty(t, s.Matches("a"))
	assert.Empty(t, s.Matches("1"))

	// matches overlapping patterns
	assert.Equal(t, []*match.Match{
		{
			Pattern:       "sequence",
			Token:         "abc",
			I:             0,
			J:             2,
			Ascending:     true,
			SequenceName:  "lower",
			SequenceSpace: 26,
		},
		{
			Pattern:       "sequence",
			Token:         "cba",
			I:             2,
			J:             4,
			Ascending:     false,
			SequenceName:  "lower",
			SequenceSpace: 26,
		},
		{
			Pattern:       "sequence",
			Token:         "abc",
			I:             4,
			J:             6,
			Ascending:     true,
			SequenceName:  "lower",
			SequenceSpace: 26,
		},
	}, s.Matches("abcbabc"))

	// matches embedded sequence patterns
	word := "jihg"
	for _, pv := range genpws(word, []string{"!", "22"}, []string{"!", "22"}) {
		assert.Equal(t, []*match.Match{
			{
				Pattern:       "sequence",
				Token:         word,
				I:             pv.i,
				J:             pv.j,
				Ascending:     false,
				SequenceName:  "lower",
				SequenceSpace: 26,
			}}, s.Matches(pv.password))

	}

	// matches pattern with the right sequence type
	tests := []struct {
		pattern   string
		name      string
		ascending bool
		space     int
	}{
		{"ABC", "upper", true, 26},
		{"CBA", "upper", false, 26},
		{"PQR", "upper", true, 26},
		{"RQP", "upper", false, 26},
		{"XYZ", "upper", true, 26},
		{"ZYX", "upper", false, 26},
		{"abcd", "lower", true, 26},
		{"dcba", "lower", false, 26},
		{"jihg", "lower", false, 26},
		{"wxyz", "lower", true, 26},
		{"zxvt", "lower", false, 26},
		{"0369", "digits", true, 10},
		{"97531", "digits", false, 10},
	}
	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			matches := s.Matches(tt.pattern)
			assert.Equal(t, []*match.Match{{
				Pattern:       "sequence",
				Token:         tt.pattern,
				I:             0,
				J:             len(tt.pattern) - 1,
				Ascending:     tt.ascending,
				SequenceName:  tt.name,
				SequenceSpace: tt.space,
			}}, matches)
		})
	}
}
