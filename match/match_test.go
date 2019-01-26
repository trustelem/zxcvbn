package match

import "testing"

func TestToString(t *testing.T) {
	matches := []*Match{
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
			Pattern:   "regex",
			Token:     "1922",
			I:         0,
			J:         3,
			RegexName: "recent_year",
		},
	}

	want := `[{"pattern":"dictionary","i":0,"j":3,"token":"abcd","matched_word":"abcd","rank":4,"dictionary_name":"d1"},` +
		`{"pattern":"regex","i":0,"j":3,"token":"1922","regex_name":"recent_year"}]`

	if got := ToString(matches); got != want {
		t.Errorf("ToString() = %v, want %v", got, want)
	}
}
