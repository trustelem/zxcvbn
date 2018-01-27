package scoring_test

import (
	"testing"

	"github.com/test-go/testify/assert"
	"github.com/trustelem/zxcvbn/match"
	"github.com/trustelem/zxcvbn/scoring"
)

func TestMostGuessableMatchSequence(t *testing.T) {
	const password = "0123456789"
	tests := []struct {
		name         string
		matches      []*match.Match
		checkGuesses bool
		wantguesses  float64
		wantsequence []*match.Match
	}{
		{
			name:    "returns one bruteforce match given an empty match sequence",
			matches: []*match.Match{},
			wantsequence: []*match.Match{
				{
					Pattern: "bruteforce",
					I:       0,
					J:       9,
					Token:   password,
					Guesses: 10000000000,
				},
			},
		},
		{
			name: "returns match + bruteforce when match covers a prefix of password",
			matches: []*match.Match{
				{
					I:       0,
					J:       5,
					Guesses: 1,
				},
			},
			wantsequence: []*match.Match{
				{
					I:       0,
					J:       5,
					Guesses: 1,
				},
				{
					Pattern: "bruteforce",
					I:       6,
					J:       9,
					Token:   password[6:10],
					Guesses: 10000,
				},
			},
		},
		{
			name: "returns bruteforce + match when match covers a suffix",
			matches: []*match.Match{
				{
					I:       3,
					J:       9,
					Guesses: 1,
				},
			},
			wantsequence: []*match.Match{
				{
					Pattern: "bruteforce",
					I:       0,
					J:       2,
					Token:   password[0:3],
					Guesses: 1000,
				},
				{
					I:       3,
					J:       9,
					Guesses: 1,
				},
			},
		},
		{
			name: "returns bruteforce + match + bruteforce when match covers an infix",
			matches: []*match.Match{
				{
					I:       1,
					J:       8,
					Guesses: 1,
				},
			},
			wantsequence: []*match.Match{
				{
					Pattern: "bruteforce",
					I:       0,
					J:       0,
					Token:   password[0:1],
					Guesses: 11,
				},
				{
					I:       1,
					J:       8,
					Guesses: 1,
				},
				{
					Pattern: "bruteforce",
					I:       9,
					J:       9,
					Token:   password[9:10],
					Guesses: 11,
				},
			},
		},
		{
			name: "chooses lower-guesses match given two matches of the same span",
			matches: []*match.Match{
				{
					I:       0,
					J:       9,
					Guesses: 1,
				},
				{
					I:       0,
					J:       9,
					Guesses: 2,
				},
			},
			wantsequence: []*match.Match{
				{
					I:       0,
					J:       9,
					Guesses: 1,
				},
			},
		},
		{
			name: "chooses lower-guesses match given two matches of the same span (check ordering)",
			matches: []*match.Match{
				{
					I:       0,
					J:       9,
					Guesses: 3,
				},
				{
					I:       0,
					J:       9,
					Guesses: 2,
				},
			},
			wantsequence: []*match.Match{
				{
					I:       0,
					J:       9,
					Guesses: 2,
				},
			},
		},
		{
			name: "when m0 covers m1 and m2, choose [m0] when m0 < m1 * m2 * fact(2)",
			matches: []*match.Match{
				{
					I:       0,
					J:       9,
					Guesses: 3,
				},
				{
					I:       0,
					J:       3,
					Guesses: 2,
				},
				{
					I:       4,
					J:       9,
					Guesses: 1,
				},
			},
			wantsequence: []*match.Match{
				{
					I:       0,
					J:       9,
					Guesses: 3,
				},
			},
			checkGuesses: true,
			wantguesses:  3,
		},
		{
			name: "when m0 covers m1 and m2, choose [m1, m2] when m0 > m1 * m2 * fact(2)",
			matches: []*match.Match{
				{
					I:       0,
					J:       9,
					Guesses: 5,
				},
				{
					I:       0,
					J:       3,
					Guesses: 2,
				},
				{
					I:       4,
					J:       9,
					Guesses: 1,
				},
			},
			wantsequence: []*match.Match{
				{
					I:       0,
					J:       3,
					Guesses: 2,
				},
				{
					I:       4,
					J:       9,
					Guesses: 1,
				},
			},
			checkGuesses: true,
			wantguesses:  4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scoring.MostGuessableMatchSequence(password, tt.matches, true)
			assert.Equal(t, tt.wantsequence, result.Sequence)
			if tt.checkGuesses {
				assert.Equal(t, tt.wantguesses, result.Guesses)
			}
		})
	}

}

func TestCalcGuesses(t *testing.T) {
	// estimate_guesses returns cached guesses when available
	assert.Equal(t, float64(1), scoring.EstimateGuesses(&match.Match{Guesses: 1}, ""))
	m := &match.Match{
		Pattern: "date",
		Token:   "1977",
		Year:    1977,
		Month:   7,
		Day:     14,
	}
	// estimate_guesses delegates based on pattern
	assert.Equal(t, scoring.EstimateGuesses(m, "1977"), scoring.DateGuesses(m))
}

func TestMostGuessableMatchSequenceCoffeeScriptCompat(t *testing.T) {
	password := "eheuczkqyq"
	seq := []*match.Match{
		{
			Pattern:        "dictionary",
			I:              0,
			J:              1,
			Token:          "eh",
			MatchedWord:    "he",
			Rank:           12,
			DictionaryName: "english_wikipedia",
			Reversed:       true,
			L33t:           false},
		{
			Pattern:        "dictionary",
			I:              1,
			J:              2,
			Token:          "he",
			MatchedWord:    "he",
			Rank:           12,
			DictionaryName: "english_wikipedia",
			Reversed:       false,
			L33t:           false},
	}
	result := scoring.MostGuessableMatchSequence(password, seq, false)

	if !assert.Equal(t, []*match.Match{
		{
			Pattern: "bruteforce",
			I:       0,
			J:       9,
			Token:   "eheuczkqyq",
			Guesses: 10000000000},
	}, result.Sequence) {
		t.Logf("Got wrong sequence %s", match.ToString(result.Sequence))
	}

	password = "qwER43@!"
	seq = []*match.Match{
		{
			Pattern:      "spatial",
			I:            0,
			J:            7,
			Token:        "qwER43@!",
			Graph:        "qwerty",
			Turns:        3,
			ShiftedCount: 4},
		{
			Pattern:        "dictionary",
			I:              1,
			J:              2,
			Token:          "wE",
			MatchedWord:    "we",
			Rank:           20,
			DictionaryName: "us_tv_and_film",
			Reversed:       false,
			L33t:           false},
		{
			Pattern:        "dictionary",
			I:              2,
			J:              4,
			Token:          "ER4",
			MatchedWord:    "era",
			Rank:           744,
			DictionaryName: "english_wikipedia",
			Reversed:       false,
			L33t:           true,
			Sub:            map[string]string{"4": "a"}},
		{
			Pattern:        "dictionary",
			I:              3,
			J:              5,
			Token:          "R43",
			MatchedWord:    "rae",
			Rank:           712,
			DictionaryName: "female_names",
			Reversed:       false,
			L33t:           true,
			Sub:            map[string]string{"3": "e", "4": "a"}},
		{
			Pattern:       "sequence",
			I:             4,
			J:             5,
			Token:         "43",
			SequenceName:  "digits",
			SequenceSpace: 10,
			Ascending:     false},
		{
			Pattern:      "spatial",
			I:            4,
			J:            7,
			Token:        "43@!",
			Graph:        "dvorak",
			Turns:        1,
			ShiftedCount: 2},
	}

	expectedSeq := []*match.Match{
		{
			Pattern:      "spatial",
			I:            0,
			J:            7,
			Token:        "qwER43@!",
			Graph:        "qwerty",
			Turns:        3,
			ShiftedCount: 4,
			Guesses:      90470620.03078316},
		{
			Pattern:             "dictionary",
			I:                   1,
			J:                   2,
			Token:               "wE",
			MatchedWord:         "we",
			Rank:                20,
			DictionaryName:      "us_tv_and_film",
			Reversed:            false,
			L33t:                false,
			BaseGuesses:         20,
			UppercaseVariations: 2,
			L33tVariations:      1,
			Guesses:             50},
		{
			Pattern:             "dictionary",
			I:                   2,
			J:                   4,
			Token:               "ER4",
			MatchedWord:         "era",
			Rank:                744,
			DictionaryName:      "english_wikipedia",
			Reversed:            false,
			L33t:                true,
			Sub:                 map[string]string{"4": "a"},
			BaseGuesses:         744,
			UppercaseVariations: 2,
			L33tVariations:      2,
			Guesses:             2976},
		{
			Pattern:             "dictionary",
			I:                   3,
			J:                   5,
			Token:               "R43",
			MatchedWord:         "rae",
			Rank:                712,
			DictionaryName:      "female_names",
			Reversed:            false,
			L33t:                true,
			Sub:                 map[string]string{"3": "e", "4": "a"},
			BaseGuesses:         712,
			UppercaseVariations: 2,
			L33tVariations:      4,
			Guesses:             5696},
		{
			Pattern:       "sequence",
			I:             4,
			J:             5,
			Token:         "43",
			SequenceName:  "digits",
			SequenceSpace: 10,
			Ascending:     false,
			Guesses:       50},
		{
			Pattern:      "spatial",
			I:            4,
			J:            7,
			Token:        "43@!",
			Graph:        "dvorak",
			Turns:        1,
			ShiftedCount: 2,
			Guesses:      12960.000000000002},
	}

	result = scoring.MostGuessableMatchSequence(password, seq, true)
	for i := range seq {
		assert.Equal(t, expectedSeq[i], seq[i])
	}

	if !assert.Equal(t, []*match.Match{
		{
			Pattern: "spatial",
			I:       0,
			J:       7,
			Token:   "qwER43@!", Graph: "qwerty",
			Turns:        3,
			ShiftedCount: 4,
			Guesses:      90470620.03078316},
	}, result.Sequence) {
		t.Logf("Got wrong most guessable sequence %s", match.ToString(result.Sequence))
	}

}
