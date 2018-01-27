package matching

import (
	"fmt"
	"strings"
	"testing"

	"github.com/test-go/testify/assert"
	"github.com/trustelem/zxcvbn/match"
)

func Test_dateMatch(t *testing.T) {
	// matches dates that use '#{sep}' as a separator"
	for _, sep := range []string{"", " ", "-", "/", "\\", "_", "."} {
		password := "13" + sep + "2" + sep + "1921"
		matches := dateMatch{}.Matches(password)

		assert.Equal(t, []*match.Match{
			{
				Pattern:   "date",
				Token:     password,
				I:         0,
				J:         len(password) - 1,
				Separator: sep,
				Year:      1921,
				Month:     2,
				Day:       13,
			},
		},
			matches,
		)
	}

	// matches dates with "#{order}" format
	for _, order := range []string{"mdy", "dmy", "ymd", "ydm"} {
		password := strings.Replace(order, "y", "88", 1)
		password = strings.Replace(password, "m", "8", 1)
		password = strings.Replace(password, "d", "8", 1)
		matches := dateMatch{}.Matches(password)

		assert.Equal(t, []*match.Match{
			{
				Pattern:   "date",
				Token:     password,
				I:         0,
				J:         len(password) - 1,
				Separator: "",
				Year:      1988,
				Month:     8,
				Day:       8,
			},
		},
			matches,
		)
	}

	// matches the date with year closest to REFERENCE_YEAR when ambiguous
	password := "111504"
	assert.Equal(t, []*match.Match{
		{
			Pattern:   "date",
			Token:     password,
			I:         0,
			J:         len(password) - 1,
			Separator: "",
			Year:      2004, // picks '04' -> 2004 as year, not '1504'
			Month:     11,
			Day:       15,
		},
	},
		dateMatch{}.Matches(password),
	)

	// matches various dates
	for _, tt := range []struct {
		day   int
		month int
		year  int
	}{
		{1, 1, 1999},
		{11, 8, 2000},
		{9, 12, 2005},
		{22, 11, 1551},
	} {
		password := fmt.Sprintf("%d%d%d", tt.year, tt.month, tt.day)
		matches := dateMatch{}.Matches(password)
		month, day := matches[0].Month, matches[0].Day
		assert.Equal(t, []*match.Match{
			{
				Pattern:   "date",
				Token:     password,
				I:         0,
				J:         len(password) - 1,
				Separator: "",
				Year:      tt.year,
				Month:     month,
				Day:       day,
			},
		},
			matches,
			"matching %s", password,
		)

		password = fmt.Sprintf("%d.%d.%d", tt.year, tt.month, tt.day)
		matches = dateMatch{}.Matches(password)
		month, day = matches[0].Month, matches[0].Day
		assert.Equal(t, []*match.Match{
			{
				Pattern:   "date",
				Token:     password,
				I:         0,
				J:         len(password) - 1,
				Separator: ".",
				Year:      tt.year,
				Month:     month,
				Day:       day,
			},
		},
			matches,
			"matching %s", password,
		)

	}

	// matches zero-padded dates
	password = "02/02/02"
	assert.Equal(t, []*match.Match{
		{
			Pattern:   "date",
			Token:     password,
			I:         0,
			J:         len(password) - 1,
			Separator: "/",
			Year:      2002,
			Month:     2,
			Day:       2,
		},
	},
		dateMatch{}.Matches(password),
	)

	// matches embedded dates
	word := "1/1/91"
	for _, pv := range genpws(word, []string{"a", "ab"}, []string{"!"}) {
		assert.Equal(t, []*match.Match{
			{
				Pattern:   "date",
				Token:     word,
				I:         pv.i,
				J:         pv.j,
				Separator: "/",
				Year:      1991,
				Month:     1,
				Day:       1,
			}}, dateMatch{}.Matches(pv.password))

	}

	// matches overlapping dates
	password = "12/20/1991.12.20"
	assert.Equal(t, []*match.Match{
		{
			Pattern:   "date",
			Token:     "12/20/1991",
			I:         0,
			J:         9,
			Separator: "/",
			Year:      1991,
			Month:     12,
			Day:       20,
		},
		{
			Pattern:   "date",
			Token:     "1991.12.20",
			I:         6,
			J:         15,
			Separator: ".",
			Year:      1991,
			Month:     12,
			Day:       20,
		},
	},
		dateMatch{}.Matches(password),
	)

	// matches dates padded by non-ambiguous digits
	password = "912/20/919"
	assert.Equal(t, []*match.Match{
		{
			Pattern:   "date",
			Token:     "12/20/91",
			I:         1,
			J:         8,
			Separator: "/",
			Year:      1991,
			Month:     12,
			Day:       20,
		}}, dateMatch{}.Matches(password))
}

func Test_twoToFourDigitYear(t *testing.T) {
	tests := []struct {
		year int
		want int
	}{
		{60, 1960},
		{960, 960},
		{20, 2020},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, twoToFourDigitYear(tt.year))
	}
}
