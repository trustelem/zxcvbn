package matching

import (
	"strings"
	"testing"

	"github.com/test-go/testify/assert"
	"github.com/trustelem/zxcvbn/match"
)

// removeRepeatBaseData removes extra data not needed for unit tests
func removeRepeatBaseData(matches []*match.Match) []*match.Match {
	for _, m := range matches {
		m.BaseGuesses = 0
		m.BaseMatches = nil
	}
	return matches
}

func TestRepeatMatching(t *testing.T) {
	r := repeatMatch{}

	// doesn't match 0- and 1-character repeat patterns
	assert.Empty(t, r.Matches(""))
	assert.Empty(t, r.Matches("#"))

	// test single-character repeats
	word := "&&&&&"
	for _, pv := range genpws(word, []string{"@", "y4@"}, []string{"u", "u%7"}) {
		matches := removeRepeatBaseData(r.Matches(pv.password))
		assert.Equal(t, []*match.Match{
			{
				Pattern:     "repeat",
				Token:       word,
				I:           pv.i,
				J:           pv.j,
				BaseToken:   "&",
				RepeatCount: 5,
			}}, matches)
	}

	// matches repeats with base character
	for length := 3; length <= 12; length++ {
		for _, chr := range []string{"a", "Z", "4", "&"} {
			word := strings.Repeat(chr, length)
			matches := removeRepeatBaseData(r.Matches(word))

			assert.Equal(t, []*match.Match{
				{
					Pattern:     "repeat",
					Token:       word,
					I:           0,
					J:           length - 1,
					BaseToken:   chr,
					RepeatCount: length,
				}}, matches)
		}
	}

	// matches multiple adjacent repeats
	matches := removeRepeatBaseData(r.Matches("BBB1111aaaaa@@@@@@"))
	assert.Equal(t, []*match.Match{
		{
			Pattern:     "repeat",
			Token:       "BBB",
			I:           0,
			J:           2,
			BaseToken:   "B",
			RepeatCount: 3,
		},
		{
			Pattern:     "repeat",
			Token:       "1111",
			I:           3,
			J:           6,
			BaseToken:   "1",
			RepeatCount: 4,
		},
		{
			Pattern:     "repeat",
			Token:       "aaaaa",
			I:           7,
			J:           11,
			BaseToken:   "a",
			RepeatCount: 5,
		},
		{
			Pattern:     "repeat",
			Token:       "@@@@@@",
			I:           12,
			J:           17,
			BaseToken:   "@",
			RepeatCount: 6,
		},
	}, matches)

	// matches multiple repeats with non-repeats in-between
	matches = removeRepeatBaseData(r.Matches("2818BBBbzsdf1111@*&@!aaaaaEUDA@@@@@@1729"))
	assert.Equal(t, []*match.Match{
		{
			Pattern:     "repeat",
			Token:       "BBB",
			I:           4,
			J:           6,
			BaseToken:   "B",
			RepeatCount: 3,
		},
		{
			Pattern:     "repeat",
			Token:       "1111",
			I:           12,
			J:           15,
			BaseToken:   "1",
			RepeatCount: 4,
		},
		{
			Pattern:     "repeat",
			Token:       "aaaaa",
			I:           21,
			J:           25,
			BaseToken:   "a",
			RepeatCount: 5,
		},
		{
			Pattern:     "repeat",
			Token:       "@@@@@@",
			I:           30,
			J:           35,
			BaseToken:   "@",
			RepeatCount: 6,
		},
	}, matches)

	// test multi-character repeats
	matches = removeRepeatBaseData(r.Matches("abab"))
	assert.Equal(t, []*match.Match{
		{
			Pattern:     "repeat",
			Token:       "abab",
			I:           0,
			J:           3,
			BaseToken:   "ab",
			RepeatCount: 2,
		},
	}, matches)

	// matches aabaab as a repeat instead of the aa prefix
	matches = removeRepeatBaseData(r.Matches("aabaab"))
	assert.Equal(t, []*match.Match{
		{
			Pattern:     "repeat",
			Token:       "aabaab",
			I:           0,
			J:           5,
			BaseToken:   "aab",
			RepeatCount: 2,
		},
	}, matches)

	// identifies ab as repeat string, even though abab is also repeated
	matches = removeRepeatBaseData(r.Matches("abababab"))
	assert.Equal(t, []*match.Match{
		{
			Pattern:     "repeat",
			Token:       "abababab",
			I:           0,
			J:           7,
			BaseToken:   "ab",
			RepeatCount: 4,
		},
	}, matches)
}

func TestCornerCases(t *testing.T) {
	// cases found in fuzzing
	testCases := []string{
		"�\u007f\x00\x00Q",
	}

	r := repeatMatch{}

	for _, password := range testCases {
		_ = r.Matches(password)
	}
}
