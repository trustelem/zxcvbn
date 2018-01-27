package matching

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"strconv"
	"testing"

	"github.com/test-go/testify/assert"
	"github.com/trustelem/zxcvbn/match"
)

var testl33tTable = map[string][]string{
	"a": {"4", "@"},
	"c": {"(", "{", "[", "<"},
	"g": {"6", "9"},
	"o": {"0"},
}

func Test_relevantSubtable(t *testing.T) {
	// reduces l33t table to only the substitutions that a password might be employing
	tests := []struct {
		password string
		want     map[string][]string
	}{
		{
			password: "",
			want:     map[string][]string{},
		},
		{
			password: "abcdefgo123578!#$&*)]}>",
			want:     map[string][]string{},
		},
		{
			password: "a",
			want:     map[string][]string{},
		},
		{
			password: "4",
			want:     map[string][]string{"a": {"4"}},
		},
		{
			password: "4@",
			want:     map[string][]string{"a": {"4", "@"}},
		},
		{
			password: "4({60",
			want:     map[string][]string{"a": {"4"}, "c": {"(", "{"}, "g": {"6"}, "o": {"0"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			if got := relevantSubtable(tt.password, testl33tTable); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("relevantSubtable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_enumerateLeetSubs(t *testing.T) {
	// enumerates the different sets of l33t substitutions a password might be using
	type args struct {
		table map[string][]string
	}
	tests := []struct {
		table map[string][]string
		want  []map[string]string
	}{
		{
			table: map[string][]string{},
			want:  []map[string]string{{}},
		},
		{
			table: map[string][]string{"a": {"@"}},
			want:  []map[string]string{{"@": "a"}},
		},
		{
			table: map[string][]string{"a": {"@", "4"}},
			want:  []map[string]string{{"@": "a"}, {"4": "a"}},
		},
		{
			table: map[string][]string{"a": {"@", "4"}, "c": {"("}},
			want:  []map[string]string{{"@": "a", "(": "c"}, {"4": "a", "(": "c"}},
		},
	}
	for i, tt := range tests {
		t.Run("test_"+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tt.want, enumerateLeetSubs(tt.table))
		})
	}
}

func Test_l33tMatch(t *testing.T) {
	lm := l33tMatch{
		dm: dictionaryMatch{
			rankedDictionaries: map[string]rankedDictionnary{
				"words": rankedDictionnary{
					"aac":       1,
					"password":  3,
					"paassword": 4,
					"asdf0":     5,
				},
				"words2": rankedDictionnary{
					"cgo": 1,
				},
			},
		},
		table: testl33tTable,
	}
	tests := []struct {
		name     string
		password string
		want     []*match.Match
	}{
		{
			name:     "doesn't match ''",
			password: "",
			want:     []*match.Match{},
		},
		{
			name:     "doesn't match pure dictionary words",
			password: "password",
			want:     []*match.Match{},
		},
		{
			name:     "matches against common l33t substitutions",
			password: "p4ssword",
			want: []*match.Match{
				{
					Pattern:        "dictionary",
					Token:          "p4ssword",
					MatchedWord:    "password",
					Rank:           3,
					DictionaryName: "words",
					I:              0,
					J:              7,
					L33t:           true,
					Sub:            map[string]string{"4": "a"},
				},
			},
		},
		{
			name:     "matches against common l33t substitutions",
			password: "p@ssw0rd",
			want: []*match.Match{
				{
					Pattern:        "dictionary",
					Token:          "p@ssw0rd",
					MatchedWord:    "password",
					Rank:           3,
					DictionaryName: "words",
					I:              0,
					J:              7,
					L33t:           true,
					Sub:            map[string]string{"@": "a", "0": "o"},
				},
			},
		},
		{
			name:     "matches against common l33t substitutions",
			password: "aSdfO{G0asDfO",
			want: []*match.Match{
				{
					Pattern:        "dictionary",
					Token:          "{G0",
					MatchedWord:    "cgo",
					Rank:           1,
					DictionaryName: "words2",
					I:              5,
					J:              7,
					L33t:           true,
					Sub:            map[string]string{"{": "c", "0": "o"},
				},
			},
		},
		{
			name:     "matches against overlapping l33t patterns",
			password: "@a(go{G0",
			want: []*match.Match{
				{
					Pattern:        "dictionary",
					Token:          "@a(",
					MatchedWord:    "aac",
					Rank:           1,
					DictionaryName: "words",
					I:              0,
					J:              2,
					L33t:           true,
					Sub:            map[string]string{"@": "a", "(": "c"},
				},
				{
					Pattern:        "dictionary",
					Token:          "(go",
					MatchedWord:    "cgo",
					Rank:           1,
					DictionaryName: "words2",
					I:              2,
					J:              4,
					L33t:           true,
					Sub:            map[string]string{"(": "c"},
				},
				{
					Pattern:        "dictionary",
					Token:          "{G0",
					MatchedWord:    "cgo",
					Rank:           1,
					DictionaryName: "words2",
					I:              5,
					J:              7,
					L33t:           true,
					Sub:            map[string]string{"{": "c", "0": "o"},
				},
			},
		},
		{
			name:     "doesn't match when multiple l33t substitutions are needed for the same letter",
			password: "p4@ssword",
			want:     []*match.Match{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, lm.Matches(tt.password))
		})
	}

	// doesn't match single-character l33ted words
	assert.Len(t, lm.Matches("4 1 @"), 0)

	// known issue: subsets of substitutions aren't tried.
	// for long inputs, trying every subset of every possible substitution could quickly get large,
	// but there might be a performant way to fix.
	// (so in this example: {'4': a, '0': 'o'} is detected as a possible sub,
	// but the subset {'4': 'a'} isn't tried, missing the match for asdf0.)
	// TODO: consider partially fixing by trying all subsets of size 1 and maybe 2
	assert.Len(t, lm.Matches("4sdf0"), 0)
}

func TestDeterministicOutput(t *testing.T) {
	password := "coRrecth0rseba++ery9.23.2007staple$"

	lm := l33tMatch{
		dm:    defaultRankedDictionnaries,
		table: l33tTable,
	}

	var lastMatches []*match.Match
	for i := 0; i < 100; i++ {
		matches := lm.Matches(password)
		if i > 0 {
			if d := cmp.Diff(matches, lastMatches); d != "" {
				t.Fatalf("Got two different values %s %s \n%s", match.ToString(lastMatches), match.ToString(matches), d)
			}
		}
		lastMatches = matches
	}
}
