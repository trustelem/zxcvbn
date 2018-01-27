package matching

import (
	"github.com/test-go/testify/assert"
	"github.com/trustelem/zxcvbn/match"
	"testing"
)

func TestRegexpMatching(t *testing.T) {
	rm := regexpMatch{regexes: defaultRegexpMatch}
	assert.Equal(t, []*match.Match{
		{
			Pattern:   "regex",
			Token:     "1922",
			I:         0,
			J:         3,
			RegexName: "recent_year",
		},
	},
		rm.Matches("1922"),
	)

	assert.Equal(t, []*match.Match{
		{
			Pattern:   "regex",
			Token:     "2017",
			I:         0,
			J:         3,
			RegexName: "recent_year",
		},
	},
		rm.Matches("2017"),
	)
}
