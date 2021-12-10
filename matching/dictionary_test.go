package matching

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/trustelem/zxcvbn/match"
)

func Test_dictionaryMatch(t *testing.T) {
	dm := dictionaryMatch{
		rankedDictionaries: map[string]rankedDictionnary{
			"d1": rankedDictionnary{
				"motherboard": 1,
				"mother":      2,
				"board":       3,
				"abcd":        4,
				"cdef":        5,
				"touché":      6,
				"先生":          7,
				"猫咪":          8,
			},
			"d2": rankedDictionnary{
				"z":          1,
				"8":          2,
				"99":         3,
				"$":          4,
				"asdf1234&*": 5,
			},
		},
	}
	tests := []struct {
		name     string
		password string
		want     []*match.Match
	}{
		{
			name:     "matches words that contain other words",
			password: "motherboard",
			want: []*match.Match{
				{
					Pattern:        "dictionary",
					Token:          "mother",
					MatchedWord:    "mother",
					Rank:           2,
					DictionaryName: "d1",
					I:              0,
					J:              5,
				},
				{
					Pattern:        "dictionary",
					Token:          "motherboard",
					MatchedWord:    "motherboard",
					Rank:           1,
					DictionaryName: "d1",
					I:              0,
					J:              10,
				},
				{
					Pattern:        "dictionary",
					Token:          "board",
					MatchedWord:    "board",
					Rank:           3,
					DictionaryName: "d1",
					I:              6,
					J:              10,
				},
			},
		},
		{
			name:     "matches multiple words when they overlap",
			password: "abcdef",
			want: []*match.Match{
				{
					Pattern:        "dictionary",
					Token:          "abcd",
					MatchedWord:    "abcd",
					Rank:           4,
					DictionaryName: "d1",
					I:              0,
					J:              3,
				},
				{
					Pattern:        "dictionary",
					Token:          "cdef",
					MatchedWord:    "cdef",
					Rank:           5,
					DictionaryName: "d1",
					I:              2,
					J:              5,
				},
			},
		},
		{
			name:     "ignores uppercasing",
			password: "BoaRdZ",
			want: []*match.Match{
				{
					Pattern:        "dictionary",
					Token:          "BoaRd",
					MatchedWord:    "board",
					Rank:           3,
					DictionaryName: "d1",
					I:              0,
					J:              4,
				},
				{
					Pattern:        "dictionary",
					Token:          "Z",
					MatchedWord:    "z",
					Rank:           1,
					DictionaryName: "d2",
					I:              5,
					J:              5,
				},
			},
		},
		{
			name:     "matches words that contain composite characters",
			password: "touchÉ先生猫咪",
			want: []*match.Match{
				{
					Pattern:        "dictionary",
					Token:          "touchÉ",
					MatchedWord:    "touché",
					Rank:           6,
					DictionaryName: "d1",
					I:              0,
					J:              6,
				},
				{
					Pattern:        "dictionary",
					Token:          "先生",
					MatchedWord:    "先生",
					Rank:           7,
					DictionaryName: "d1",
					I:              7,
					J:              12,
				},
				{
					Pattern:        "dictionary",
					Token:          "猫咪",
					MatchedWord:    "猫咪",
					Rank:           8,
					DictionaryName: "d1",
					I:              13,
					J:              18,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, dm.Matches(tt.password))
		})
	}

	// identifies words surrounded by non-words
	word := "asdf1234&*"
	for _, pv := range genpws(word, []string{"q", "%%"}, []string{"%", "qq"}) {
		assert.Equal(t, []*match.Match{
			{
				Pattern:        "dictionary",
				Token:          word,
				MatchedWord:    word,
				Rank:           5,
				DictionaryName: "d2",
				I:              pv.i,
				J:              pv.j,
			}}, dm.Matches(pv.password))

	}

	// matches against all words in provided dictionaries
	for name, dict := range dm.rankedDictionaries {
		for word, rank := range dict {
			if word == "motherboard" {
				// skip words that contain others
				continue
			}
			assert.Equal(t, []*match.Match{
				{
					Pattern:        "dictionary",
					Token:          word,
					MatchedWord:    word,
					Rank:           rank,
					DictionaryName: name,
					I:              0,
					J:              len(word) - 1,
				}}, dm.Matches(word))
		}
	}
}

func Test_defaultdictionary(t *testing.T) {
	got := defaultRankedDictionnaries.Matches("wow")
	assert.Equal(t, []*match.Match{
		{
			Pattern:        "dictionary",
			Token:          "wow",
			MatchedWord:    "wow",
			Rank:           322,
			DictionaryName: "us_tv_and_film",
			I:              0,
			J:              2,
		}}, got)

	d := defaultRankedDictionnaries.withDict(
		"user_inputs",
		buildRankedDict([]string{"foo", "bar"}),
	)
	var filtered []*match.Match
	for _, m := range d.Matches("foobar") {
		if m.DictionaryName == "user_inputs" {
			filtered = append(filtered, m)
		}
	}
	assert.Equal(t, []*match.Match{
		{
			Pattern:        "dictionary",
			Token:          "foo",
			MatchedWord:    "foo",
			Rank:           1,
			DictionaryName: "user_inputs",
			I:              0,
			J:              2,
		},
		{
			Pattern:        "dictionary",
			Token:          "bar",
			MatchedWord:    "bar",
			Rank:           2,
			DictionaryName: "user_inputs",
			I:              3,
			J:              5,
		},
	}, filtered)
}
