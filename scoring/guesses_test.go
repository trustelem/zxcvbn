package scoring_test

import (
	"math"
	"testing"

	"github.com/test-go/testify/assert"
	"github.com/trustelem/zxcvbn/adjacency"
	"github.com/trustelem/zxcvbn/internal/mathutils"
	"github.com/trustelem/zxcvbn/match"
	"github.com/trustelem/zxcvbn/matching"
	"github.com/trustelem/zxcvbn/scoring"
)

func TestRepeatGuesses(t *testing.T) {
	tests := []struct {
		Token       string
		BaseToken   string
		RepeatCount int
	}{
		{"aa", "a", 2},
		{"999", "9", 3},
		{"$$$$", "$", 4},
		{"abab", "ab", 2},
		{"batterystaplebatterystaplebatterystaple", "batterystaple", 3},
	}
	for _, tt := range tests {
		baseGuesses := scoring.MostGuessableMatchSequence(
			tt.BaseToken,
			matching.Omnimatch(tt.BaseToken, nil),
			false,
		).Guesses
		match := &match.Match{
			Token:       tt.Token,
			BaseToken:   tt.BaseToken,
			BaseGuesses: baseGuesses,
			RepeatCount: tt.RepeatCount,
		}
		expectedGuesses := baseGuesses * float64(tt.RepeatCount)
		// the repeat pattern '#{token}' has guesses of #{expected_guesses}
		assert.Equal(t, scoring.RepeatGuesses(match), expectedGuesses)
	}
}

func TestSequenceGuesses(t *testing.T) {
	tests := []struct {
		Token     string
		Ascending bool
		Guesses   float64
	}{
		{"ab", true, 4 * 2},         // obvious start * len-2
		{"XYZ", true, 26 * 3},       // base26 * len-3
		{"4567", true, 10 * 4},      // base10 * len-4
		{"7654", false, 10 * 4 * 2}, // base10 * len 4 * descending
		{"ZYX", false, 4 * 3 * 2},   // obvious start * len-3 * descending
	}
	for _, tt := range tests {
		guesses := scoring.SequenceGuesses(&match.Match{
			Token:     tt.Token,
			Ascending: tt.Ascending,
		})
		// the repeat pattern '#{token}' has guesses of #{expected_guesses}
		assert.Equal(t, tt.Guesses, guesses)
	}
}

func TestRegexGuesses(t *testing.T) {
	// guesses of 26^7 for 7-char lowercase regex
	assert.Equal(t, math.Pow(26, 7), scoring.RegexGuesses(&match.Match{
		Token:     "aizocdk",
		RegexName: "alpha_lower",
	}))

	// guesses of 62^5 for 5-char alphanumeric regex
	assert.Equal(t, math.Pow(2*26+10, 5), scoring.RegexGuesses(&match.Match{
		Token:     "ag7C8",
		RegexName: "alphanumeric",
	}))

	// "guesses of |year - REFERENCE_YEAR| for distant year matches"
	assert.EqualValues(t, mathutils.Abs(scoring.ReferenceYear-1972), scoring.RegexGuesses(&match.Match{
		Token:     "1972",
		RegexName: "recent_year",
	}))

	assert.EqualValues(t, mathutils.Abs(scoring.MinYearSpace), scoring.RegexGuesses(&match.Match{
		Token:     "2005",
		RegexName: "recent_year",
	}))
}

func TestDateGuesses(t *testing.T) {
	// guesses for #{match.token} is 365 * distance_from_ref_year
	m := &match.Match{
		Token: "1923",
		Year:  1923,
		Month: 1,
		Day:   1,
	}
	assert.EqualValues(t, 365*mathutils.Abs(scoring.ReferenceYear-m.Year), scoring.DateGuesses(m))
	// recent years assume MIN_YEAR_SPACE
	// extra guesses are added for separators.
	m = &match.Match{
		Token:     "1/1/2010",
		Year:      2010,
		Month:     1,
		Day:       1,
		Separator: "/",
	}
	assert.EqualValues(t, 365*scoring.MinYearSpace*4, scoring.DateGuesses(m))

}

func TestSpatialGuesses(t *testing.T) {
	keyboardStartingPositions := float64(len(adjacency.Graphs["qwerty"].Graph))

	// with no turns or shifts, guesses is starts * degree * (len-1)
	m := &match.Match{
		Token:        "zxcvbn",
		Graph:        "qwerty",
		Turns:        1,
		ShiftedCount: 0,
	}
	baseGuesses := keyboardStartingPositions *
		adjacency.Graphs["qwerty"].AverageDegree *
		//     # - 1 term because: not counting spatial patterns of length 1
		//     # eg for length==6, multiplier is 5 for needing to try len2,len3,..,len6
		float64(len(m.Token)-1)

	assert.Equal(t, baseGuesses, scoring.SpatialGuesses(m))

	// guesses is added for shifted keys, similar to capitals in dictionary matching
	m = &match.Match{
		Token:        "ZxCvbn",
		Graph:        "qwerty",
		Turns:        1,
		ShiftedCount: 2,
	}
	shiftedGuesses := baseGuesses * (mathutils.NCk(6, 2) + mathutils.NCk(6, 1))
	assert.Equal(t, shiftedGuesses, scoring.SpatialGuesses(m))

	// when everything is shifted, guesses are doubled
	m = &match.Match{
		Token:        "ZXCVBN",
		Graph:        "qwerty",
		Turns:        1,
		ShiftedCount: 6,
	}
	shiftedGuesses = baseGuesses * 2
	assert.Equal(t, shiftedGuesses, scoring.SpatialGuesses(m))

	// spatial guesses accounts for turn positions, directions and starting keys
	m = &match.Match{
		Token:        "zxcft6yh",
		Graph:        "qwerty",
		Turns:        2,
		ShiftedCount: 0,
	}
	guesses := float64(0)
	l := len(m.Token)
	s := keyboardStartingPositions
	d := adjacency.Graphs["qwerty"].AverageDegree
	for i := 2; i <= l; i++ {
		for j := 1; j <= m.Turns && j <= i-1; j++ {
			guesses += mathutils.NCk(i-1, j-1) * s * math.Pow(d, float64(j))
		}
	}
	assert.Equal(t, guesses, scoring.SpatialGuesses(m))
}

func TestDictionaryGuess(t *testing.T) {
	// base guesses == the rank
	assert.EqualValues(t, 32, scoring.DictionaryGuesses(&match.Match{
		Token: "aaaaa",
		Rank:  32,
	}))

	// extra guesses are added for capitalization
	assert.EqualValues(t, 32*scoring.UppercaseVariations("AAAaaa"), scoring.DictionaryGuesses(&match.Match{
		Token: "AAAaaa",
		Rank:  32,
	}))

	// guesses are doubled when word is reversed
	assert.EqualValues(t, 32*2, scoring.DictionaryGuesses(&match.Match{
		Token:    "aaa",
		Rank:     32,
		Reversed: true,
	}))

	// extra guesses are added for common l33t substitutions
	m := &match.Match{
		Token: "aaa@@@",
		Rank:  32,
		L33t:  true,
		Sub:   map[string]string{"@": "a"},
	}
	assert.EqualValues(t, 32*scoring.L33tVariations(m), scoring.DictionaryGuesses(m))

	// extra guesses are added for both capitalization and common l33t substitutions
	m = &match.Match{
		Token: "AaA@@@",
		Rank:  32,
		L33t:  true,
		Sub:   map[string]string{"@": "a"},
	}
	assert.EqualValues(t, 32*scoring.L33tVariations(m)*scoring.UppercaseVariations(m.Token), scoring.DictionaryGuesses(m))
}

func TestUppercaseVariants(t *testing.T) {
	tests := []struct {
		Word     string
		Variants float64
	}{
		{"", 1},
		{"a", 1},
		{"A", 2},
		{"abcdef", 1},
		{"Abcdef", 2},
		{"abcdeF", 2},
		{"ABCDEF", 2},
		{"aBcdef", mathutils.NCk(6, 1)},
		{"aBcDef", mathutils.NCk(6, 1) + mathutils.NCk(6, 2)},
		{"ABCDEf", mathutils.NCk(6, 1)},
		{"aBCDEf", mathutils.NCk(6, 1) + mathutils.NCk(6, 2)},
		{"ABCdef", mathutils.NCk(6, 1) + mathutils.NCk(6, 2) + mathutils.NCk(6, 3)},
	}
	for _, tt := range tests {
		// check guess multiplier of word
		assert.Equal(t, tt.Variants, scoring.UppercaseVariations(tt.Word))
	}
}

func TestL33tVariants(t *testing.T) {
	// 1 variant for non-l33t matches
	assert.Equal(t, float64(1), scoring.L33tVariations(&match.Match{L33t: false}))

	// extra l33t guesses of #{word} is #{variants}"
	for _, tt := range []struct {
		Word     string
		Variants float64
		Sub      map[string]string
	}{
		{"", 1, map[string]string{}},
		{"a", 1, map[string]string{}},
		{"4", 2, map[string]string{"4": "a"}},
		{"4pple", 2, map[string]string{"4": "a"}},
		{"abcet", 1, map[string]string{}},
		{"4bcet", 2, map[string]string{"4": "a"}},
		{"a8cet", 2, map[string]string{"8": "b"}},
		{"abce+", 2, map[string]string{"+": "t"}},
		{"48cet", 4, map[string]string{"4": "a", "8": "b"}},
		{"a4a4aa", mathutils.NCk(6, 2) + mathutils.NCk(6, 1), map[string]string{"4": "a"}},
		{"4a4a44", mathutils.NCk(6, 2) + mathutils.NCk(6, 1), map[string]string{"4": "a"}},
		{"a44att+", (mathutils.NCk(4, 2) + mathutils.NCk(4, 1)) * mathutils.NCk(3, 1), map[string]string{"4": "a", "+": "t"}},
	} {
		m := &match.Match{Token: tt.Word, Sub: tt.Sub, L33t: len(tt.Sub) > 0}
		assert.Equal(t, tt.Variants, scoring.L33tVariations(m))
	}

	// capitalization doesn't affect extra l33t guesses calc
	m := &match.Match{Token: "Aa44aA", Sub: map[string]string{"4": "a"}, L33t: true}
	variants := mathutils.NCk(6, 2) + mathutils.NCk(6, 1)
	assert.Equal(t, variants, scoring.L33tVariations(m))
}
