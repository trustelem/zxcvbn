package matching

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/trustelem/zxcvbn/match"
	"os"
	"testing"
)

func TestOmnimatch(t *testing.T) {
	assert.Empty(t, Omnimatch("", nil))
	password := "r0sebudmaelstrom11/20/91aaaa"
	matches := Omnimatch(password, nil)
	for _, tt := range []struct {
		pattern string
		i       int
		j       int
	}{
		{"dictionary", 0, 6},
		{"dictionary", 7, 15},
		{"date", 16, 23},
		{"repeat", 24, 27},
	} {
		found := false
		for _, m := range matches {
			if m.I == tt.i && m.J == tt.j && m.Pattern == tt.pattern {
				found = true
				break
			}
		}
		assert.True(t, found, "Pattern %s (i=%d, j=%d) not found")
	}

	password = "abcde"
	matches = Omnimatch(password, nil)
	assert.Equal(t, []*match.Match{
		{
			Pattern:       "sequence",
			I:             0,
			J:             4,
			Token:         "abcde",
			SequenceName:  "lower",
			SequenceSpace: 26,
			Ascending:     true},
		{
			Pattern:      "spatial",
			I:            2,
			J:            4,
			Token:        "cde",
			Graph:        "qwerty",
			Turns:        1,
			ShiftedCount: 0},
	}, matches)

	password = "qwER43@!"
	matches = Omnimatch(password, nil)
	json.NewEncoder(os.Stdout).Encode(matches)
	assert.Equal(t, []*match.Match{
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
	}, matches)

	password = "eheuczkqyq"
	matches = Omnimatch(password, nil)
	assert.Equal(t, []*match.Match{
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
	}, matches)
}
