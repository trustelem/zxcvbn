package zxcvbn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/trustelem/zxcvbn/match"

	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

func TestPasswordStrength(t *testing.T) {
	var testdata []struct {
		Password string         `json:"password"`
		Guesses  float64        `json:"guesses"`
		Score    int            `json:"score"`
		Sequence []*match.Match `json:"sequence"`
	}

	b, err := ioutil.ReadFile(filepath.Join("testdata", "output.json"))
	require.NoError(t, err)

	err = json.Unmarshal(b, &testdata)
	require.NoError(t, err)

	for _, td := range testdata {
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
					if !assert.Equal(t, td.Sequence[j].Guesses, s.Sequence[j].Guesses, msg("guesses")) {
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
			assert.Equal(t, td.Guesses, s.Guesses)
			assert.Equal(t, td.Score, s.Score, "Wrong score")
		})
	}

}

func TestCornerCases(t *testing.T) {
	testdata := []string{
		"",
		"wen\x8e\xc6",
		"İҦİ",
		"\xcd|",
		"0\xefi",
	}

	for _, td := range testdata {
		_ = PasswordStrength(td, nil)
	}
}
