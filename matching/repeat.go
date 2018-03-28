package matching

import (
	"github.com/dlclark/regexp2"
	"github.com/trustelem/zxcvbn/match"
	"github.com/trustelem/zxcvbn/scoring"
)

type repeatMatch struct{}

var greedy = regexp2.MustCompile(`(.+)\1+`, 0)
var lazy = regexp2.MustCompile(`(.+?)\1+`, 0)
var lazyAnchored = regexp2.MustCompile(`^(.+?)\1+$`, 0)

func (repeatMatch) Matches(password string) []*match.Match {
	var matches []*match.Match

	lastIndex := 0
	indices := map[int]bool{
		lastIndex: true,
	}

	for lastIndex < len(password) {
		greedyMatch, err := greedy.FindStringMatchStartingAt(password, lastIndex)
		if err != nil || greedyMatch == nil {
			break
		}
		lazyMatch, _ := lazy.FindStringMatchStartingAt(password, lastIndex)

		var rmatch *regexp2.Match
		var baseToken string
		if greedyMatch.Captures[0].Length > lazyMatch.Captures[0].Length {
			// greedy beats lazy for 'aabaab'
			//   greedy: [aabaab, aab]
			//   lazy:   [aa,     a]
			rmatch = greedyMatch
			// greedy's repeated string might itself be repeated, eg.
			// aabaab in aabaabaabaab.
			// run an anchored lazy match on greedy's repeated string
			// to find the shortest repeated string
			if m, err := lazyAnchored.FindStringMatch(rmatch.Captures[0].String()); err == nil {
				baseToken = m.GroupByNumber(1).String()
			}
		} else {
			// lazy beats greedy for 'aaaaa'
			//   greedy: [aaaa,  aa]
			//   lazy:   [aaaaa, a]
			rmatch = lazyMatch
			baseToken = rmatch.GroupByNumber(1).String()
		}
		i := rmatch.Index
		j := rmatch.Index + rmatch.Captures[0].Length - 1
		// recursively match and score the base string
		baseAnalysis := scoring.MostGuessableMatchSequence(
			baseToken,
			Omnimatch(baseToken, nil),
			false,
		)
		matches = append(matches, &match.Match{
			Pattern:     "repeat",
			I:           i,
			J:           j,
			Token:       rmatch.Captures[0].String(),
			BaseToken:   baseToken,
			BaseGuesses: baseAnalysis.Guesses,
			BaseMatches: baseAnalysis.Sequence,
			RepeatCount: rmatch.Captures[0].Length / len(baseToken),
		})
		lastIndex = j + 1
		if _, ok := indices[lastIndex]; ok {
			//already seen this index...avoid an infinite loop
			break
		}
		indices[lastIndex] = true
	}
	return matches
}
