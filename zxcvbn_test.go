package zxcvbn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/trustelem/zxcvbn/match"
	"github.com/trustelem/zxcvbn/scoring"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordStrength(t *testing.T) {
	var testdata struct {
		TimeStamp time.Time `json:"timestamp"`
		Tests     []struct {
			Password string         `json:"password"`
			Guesses  float64        `json:"guesses"`
			Score    int            `json:"score"`
			Sequence []*match.Match `json:"sequence"`
		} `json:"tests"`
	}

	b, err := ioutil.ReadFile(filepath.Join("testdata", "output.json"))
	require.NoError(t, err)

	err = json.Unmarshal(b, &testdata)
	require.NoError(t, err)

	refYear := scoring.ReferenceYear
	defer func() {
		scoring.ReferenceYear = refYear
	}()
	scoring.ReferenceYear = testdata.TimeStamp.Year()
	// maximum epsilon for guesses comparison
	const maxEpsilonGuesses = 1e-15
	for _, td := range testdata.Tests {
		t.Run(td.Password, func(t *testing.T) {
			// map character positions to rune position
			runeMap := make(map[int]int, len(td.Password))
			c := 0
			for i := range td.Password {
				runeMap[i] = c
				c++
			}
			runeMap[len(td.Password)] = c
			s := PasswordStrength(td.Password, nil)
			if len(s.Sequence) == len(td.Sequence) {
				for j := range td.Sequence {
					expect, _ := json.Marshal(td.Sequence[j])
					got, _ := json.Marshal(s.Sequence[j])
					msg := func(f string) string {
						return fmt.Sprintf("Password %+q, field %s: expect=%s got=%s",
							td.Password,
							f,
							string(expect),
							string(got))
					}
					if !assert.Equal(t, td.Sequence[j].I, runeMap[s.Sequence[j].I], msg("i")) {
						return
					}
					if !assert.Equal(t, td.Sequence[j].J, runeMap[s.Sequence[j].J+1]-1, msg("j")) {
						t.Logf("runeMap %v\n", runeMap)
						return
					}
					if !assert.Equal(t, td.Sequence[j].Pattern, s.Sequence[j].Pattern, msg("pattern")) {
						return
					}
					if !assert.Equal(t, td.Sequence[j].Token, s.Sequence[j].Token, msg("token")) {
						return
					}
					if !assert.InEpsilon(t, td.Sequence[j].Guesses, s.Sequence[j].Guesses, maxEpsilonGuesses, msg("guesses")) {
						return
					}
				}
			} else {
				b, _ := json.Marshal(td.Sequence)
				t.Errorf("Expected sequence:\n%s\nGot:\n%s\n",
					string(b),
					match.ToString(s.Sequence))
				return
			}
			assert.InEpsilon(t, td.Guesses, s.Guesses, maxEpsilonGuesses)
			assert.Equal(t, td.Score, s.Score, "Wrong score")
		})
	}

}

func TestCornerCases(t *testing.T) {
	testdata := []string{
		"",
		"wen\x8e\xc6",
	}

	for _, td := range testdata {
		_ = PasswordStrength(td, nil)
	}
}
