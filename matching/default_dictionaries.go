// +build !use_trie

package matching

import "github.com/trustelem/zxcvbn/frequency"

func loadDefaultDictionaries() dictionaryMatch {
	rd := make(map[string]rankedDictionary)
	for n, list := range frequency.FrequencyLists {
		rd[n] = buildRankedDict(list)
	}
	return dictionaryMatch{
		rankedDictionaries: rd,
	}
}
