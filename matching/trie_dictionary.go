// +build use_trie

package matching

import (
	"github.com/okzk/marisa"
	"github.com/trustelem/zxcvbn/match"
	"strings"
)

type trieRankedDictionary struct {
	trie      *marisa.Trie
	rankTable []int
}

func newTrieRankedDictionary(unrankedList []string) *trieRankedDictionary {
	keyset := marisa.NewKeyset()
	defer keyset.Dispose()
	for _, v := range unrankedList {
		err := keyset.PushBack(strings.ToLower(v))
		if err != nil {
			panic(err)
		}
	}
	trie, err := keyset.Build()
	if err != nil {
		panic(err)
	}

	rankTable := make([]int, len(unrankedList))
	for i, v := range unrankedList {
		id, err := trie.Lookup(strings.ToLower(v))
		if err != nil {
			panic(err)
		}
		rankTable[id] = i + 1
	}
	return &trieRankedDictionary{
		trie:      trie,
		rankTable: rankTable,
	}
}

func (dict *trieRankedDictionary) matches(dictionaryName, password string) []*match.Match {
	var results []*match.Match
	lowerPassword := strings.ToLower(password)
	for i := range password {
		dict.trie.CommonPrefixSearch(lowerPassword[i:], func(id uint64, key string) error {
			j := i + len(key) - 1
			matchDic := &match.Match{
				Pattern:        "dictionary",
				I:              i,
				J:              j,
				Token:          password[i : j+1],
				MatchedWord:    key,
				Rank:           dict.rankTable[id],
				DictionaryName: dictionaryName,
			}
			results = append(results, matchDic)
			return nil
		})
	}
	return results
}
