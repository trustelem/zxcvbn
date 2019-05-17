package feedback

import (
	"encoding/json"
	"strings"

	"github.com/trustelem/zxcvbn/match"
	"github.com/trustelem/zxcvbn/scoring"
)

type Feedback struct {
	Warning     string   `json:"warning"`
	Suggestions []string `json:"suggestions"`
}

func GetFeedback(score int, sequence []*match.Match) *Feedback {
	// starting feedback
	if len(sequence) == 0 {
		return defaultFeedback()
	}

	// no feedback if score is good or great.
	if score > 2 {
		return emptyFeedback()
	}

	// tie feedback to the longest match for longer sequences
	longestMatch := sequence[0]
	for _, match := range sequence[1:] {
		if len(match.Token) > len(longestMatch.Token) {
			longestMatch = match
		}
	}

	feedback := getMatchFeedback(longestMatch, len(sequence) == 1)

	extraFeedback := "Add another word or two. Uncommon words are better."

	if feedback == nil {
		return &Feedback{
			Warning:     "",
			Suggestions: []string{extraFeedback},
		}
	}

	feedback.Suggestions = append([]string{extraFeedback}, feedback.Suggestions...)

	return feedback
}

func getMatchFeedback(match *match.Match, isSoleMatch bool) *Feedback {
	switch match.Pattern {
	case "dictionary":
		return getDictionaryMatchFeedback(match, isSoleMatch)
	case "spatial":
		var warning string
		if match.Turns == 1 {
			warning = "Straight rows of keys are easy to guess"
		} else {
			warning = "Short keyboard patterns are easy to guess"
		}

		return &Feedback{
			Warning:     warning,
			Suggestions: []string{"Use a longer keyboard pattern with more turns"},
		}
	case "repeat":
		var warning string
		if len(match.BaseToken) == 1 {
			warning = "Repeats like \"aaa\" are easy to guess"
		} else {
			warning = "Repeats like \"abcabcabc\" are only slightly harder to guess than \"abc\""
		}

		return &Feedback{
			Warning:     warning,
			Suggestions: []string{"Avoid repeated words and characters"},
		}
	case "sequence":
		return &Feedback{
			Warning:     "Sequences like abc or 6543 are easy to guess",
			Suggestions: []string{"Avoid sequences"},
		}
	case "regex":
		if match.RegexName == "recent_year" {
			return &Feedback{
				Warning: "Recent years are easy to guess",
				Suggestions: []string{
					"Avoid recent years",
					"Avoid years that are associated with you",
				},
			}
		}

		return emptyFeedback()
	case "date":
		return &Feedback{
			Warning:     "Dates are often easy to guess",
			Suggestions: []string{"Avoid dates and years that are associated with you"},
		}
	default:
		return emptyFeedback()
	}
}

func getDictionaryMatchFeedback(match *match.Match, isSoleMatch bool) *Feedback {
	var warning string

	if match.DictionaryName == "passwords" {
		if isSoleMatch && !match.L33t && !match.Reversed {
			if match.Rank <= 10 {
				warning = "This is a top-10 common password"
			} else if match.Rank <= 100 {
				warning = "This is a top-100 common password"
			} else {
				warning = "This is a very common password"
			}
		}
	} else if match.GuessesLog10 <= 4 {
		warning = "This is similar to a commonly used password"
	} else if match.DictionaryName == "english_wikipedia" {
		if isSoleMatch {
			warning = "A word by itself is easy to guess"
		}
	} else if stringInList(match.DictionaryName, []string{"surnames", "male_names", "female_names"}) {
		if isSoleMatch {
			warning = "Names and surnames by themselves are easy to guess"
		} else {
			warning = "Common names and surnames are easy to guess"
		}
	}

	var suggestions []string
	word := match.Token
	if scoring.ReStartUpper.MatchString(word) {
		suggestions = append(suggestions, "Capitalization doesn't help very much")
	} else if scoring.ReAllUpper.MatchString(word) && strings.ToLower(word) != word {
		suggestions = append(suggestions, "All-uppercase is almost as easy to guess as all-lowercase")
	}

	if match.Reversed && len(match.Token) == 4 {
		suggestions = append(suggestions, "Reversed words aren't much harder to guess")
	}
	if match.L33t {
		suggestions = append(suggestions, "Predictable substitutions like '@' instead of 'a' don't help very much")
	}

	return &Feedback{
		Warning:     warning,
		Suggestions: suggestions,
	}
}

func stringInList(x string, list []string) bool {
	for _, item := range list {
		if item == x {
			return true
		}
	}
	return false
}

// ToString returns a string representation of the feedback
func ToString(feedback *Feedback) string {
	b, _ := json.Marshal(feedback)
	return string(b)
}

func defaultFeedback() *Feedback {
	return &Feedback{
		Warning: "",
		Suggestions: []string{
			"Use a few words, avoid common phrases",
			"No need for symbols, digits, or uppercase letters"},
	}
}

// Why not just return nil? In practice it's much the same, but, there is a
// difference between a nil and a empty slice. The main motivation for doing it
// this way was when deserialising the test data, which would deserialise to
// an empty slice, and not a nil slice.
func emptyFeedback() *Feedback {
	return &Feedback{
		Warning:     "",
		Suggestions: []string{},
	}
}
