package matching

import (
	"strings"

	"github.com/trustelem/zxcvbn/match"
)

type dictionaryMatch struct {
	rankedDictionaries map[string]rankedDictionary
}

func (dm dictionaryMatch) Matches(password string) []*match.Match {
	var results []*match.Match

	for dictionaryName, rankedDict := range dm.rankedDictionaries {
		results = append(results, rankedDict.matches(dictionaryName, password)...)
	}

	match.Sort(results)
	return results
}

func (dm dictionaryMatch) withDict(name string, d rankedDictionary) dictionaryMatch {
	rd2 := make(map[string]rankedDictionary, len(dm.rankedDictionaries)+1)
	for k, v := range dm.rankedDictionaries {
		rd2[k] = v
	}
	rd2[name] = d
	return dictionaryMatch{rankedDictionaries: rd2}
}

type rankedDictionary interface {
	matches(dictionaryName, password string) []*match.Match
}

type mapRankedDictionary map[string]int

func newMapRankedDictionary(unrankedList []string) mapRankedDictionary {
	result := make(mapRankedDictionary)

	for i, v := range unrankedList {
		result[strings.ToLower(v)] = i + 1
	}
	return result
}

func (dict mapRankedDictionary) matches(dictionaryName, password string) []*match.Match {
	var results []*match.Match
	lowerPassword := strings.ToLower(password)
	for i := range password {
		for delta := range password[i:] {
			j := i + delta
			word := lowerPassword[i : j+1]
			if val, ok := dict[word]; ok {
				matchDic := &match.Match{
					Pattern:        "dictionary",
					I:              i,
					J:              j,
					Token:          password[i : j+1],
					MatchedWord:    word,
					Rank:           val,
					DictionaryName: dictionaryName,
				}
				results = append(results, matchDic)
			}
		}
	}
	return results
}

type linearSearchRankedDictionary []string

func newLinearSearchRankedDictionary(unrankedList []string) linearSearchRankedDictionary {
	result := make(linearSearchRankedDictionary, len(unrankedList))

	for i, v := range unrankedList {
		result[i] = strings.ToLower(v)
	}
	return result
}

func (dict linearSearchRankedDictionary) matches(dictionaryName, password string) []*match.Match {
	var results []*match.Match
	lowerPassword := strings.ToLower(password)
	for idx, word := range dict {
		pos := 0
		for {
			n := strings.Index(lowerPassword[pos:], word)
			if n == -1 {
				break
			}
			i := pos + n
			j := i + len(word) - 1
			matchDic := &match.Match{
				Pattern:        "dictionary",
				I:              i,
				J:              j,
				Token:          password[i : j+1],
				MatchedWord:    word,
				Rank:           idx + 1,
				DictionaryName: dictionaryName,
			}
			results = append(results, matchDic)
			pos = i + 1
		}
	}
	return results
}

func buildRankedDict(unrankedList []string) rankedDictionary {
	// linearSearchRankedDictionary is faster than mapRankedDictionary when dictionary size is small enough.
	if len(unrankedList) < 64 {
		return newLinearSearchRankedDictionary(unrankedList)
	}
	return newMapRankedDictionary(unrankedList)
}
