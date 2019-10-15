// +build use_trie

package matching

import (
	"github.com/trustelem/zxcvbn/frequency"
	"testing"
)

func Test_trieRankedDictionaryMatch(t *testing.T) {
	dm := dictionaryMatch{
		rankedDictionaries: map[string]rankedDictionary{
			"d1": newTrieRankedDictionary([]string{
				"motherboard",
				"mother",
				"board",
				"abcd",
				"cdef",
			}),
			"d2": newTrieRankedDictionary([]string{
				"z",
				"8",
				"99",
				"$",
				"asdf1234&*",
			}),
		},
	}
	test_dictionaryMatch(t, dm)
}

func BenchmarkTrieRankedDictionaryMatch_8_64(b *testing.B) {
	// password length:8, dictionary size: 64
	password := testPassword64[:8]
	dict := newTrieRankedDictionary(frequency.FrequencyLists["passwords"][:64])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.matches("test", password)
	}
}

func BenchmarkTrieRankedDictionaryMatch_64_64(b *testing.B) {
	// password length:64, dictionary size: 64
	dict := newTrieRankedDictionary(frequency.FrequencyLists["passwords"][:64])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.matches("test", testPassword64)
	}
}
