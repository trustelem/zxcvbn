package matching

import (
	"testing"

	"github.com/test-go/testify/assert"
	"github.com/trustelem/zxcvbn/match"
)

func Test_reverseDictionnaryMatch(t *testing.T) {
	rdm := reverseDictionnaryMatch{
		dm: dictionaryMatch{
			rankedDictionaries: map[string]rankedDictionnary{
				"d1": rankedDictionnary{
					"123": 1,
					"321": 2,
					"456": 3,
					"654": 4,
				},
			},
		},
	}

	password := "0123456789"
	assert.Equal(t, []*match.Match{
		{
			Pattern:        "dictionary",
			Token:          "123",
			MatchedWord:    "321",
			Rank:           2,
			DictionaryName: "d1",
			I:              1,
			J:              3,
			Reversed:       true,
		},
		{
			Pattern:        "dictionary",
			Token:          "456",
			MatchedWord:    "654",
			Rank:           4,
			DictionaryName: "d1",
			I:              4,
			J:              6,
			Reversed:       true,
		},
	}, rdm.Matches(password))
}
