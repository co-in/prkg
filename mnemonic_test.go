package prkg_test

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"strings"
	"testing"

	"github.com/co-in/prkg"
	"github.com/stretchr/testify/assert"
)

func TestMnemonicChecksum(t *testing.T) {
	// Ensure a word list is correct https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/english.txt
	checksum := crc32.ChecksumIEEE([]byte(strings.Join(prkg.DictionaryEnglish.Words(), "\n")))
	assert.Equal(t, uint32(0xb5a54d12), checksum)
}

func TestEntropyFromMnemonic(t *testing.T) {
	entropy, _ := hex.DecodeString("67776fe4e70fe34fbf9f1432546ff0f3d5e4b9d8")
	words, err := prkg.DictionaryEnglish.Mnemonic(entropy)
	assert.NoError(t, err)
	assert.Equal(t, words, []string{
		"guess", "rocket", "weird", "sock", "wreck", "pond",
		"wrist", "tip", "crane", "pet", "wire", "tray",
		"furnace", "friend", "genuine",
	})
}

func TestUniqueFirstFour(t *testing.T) {
	words := prkg.DictionaryEnglish.Words()

	var groups3 = make([]string, 0, len(words))
	var groups4 = make(map[string]int, len(words))
	var maxLen int
	var minLen int

	for idx, word := range words {
		if len(word) > 3 {
			groups4[word[:4]]++
		} else {
			groups3 = append(groups3, fmt.Sprintf("%04d: %s", idx+1, word))
		}

		if len(word) > maxLen {
			maxLen = len(word)
		}
		if len(word) < minLen || minLen == 0 {
			minLen = len(word)
		}
	}

	assert.Equal(t, 3, minLen)
	assert.Equal(t, 8, maxLen)

	for key, count := range groups4 {
		if count > 1 {
			assert.Equal(t, 1, count, key)
		}
	}

	fmt.Println(strings.Join(groups3, "\n"))
}

func TestDictionaryUnique(t *testing.T) {
	words := prkg.DictionaryEnglish.Words()
	var uniqueWords = make(map[string]bool, len(words))

	for idx, word := range words {
		assert.False(t, uniqueWords[word], fmt.Sprintf("word %d: %s", idx, word))
		uniqueWords[word] = true
	}
}

func TestSeed(t *testing.T) {
	mnemonic := []string{
		"guess", "rocket", "weird", "sock", "wreck", "pond",
		"wrist", "tip", "crane", "pet", "wire", "tray",
		"furnace", "friend", "genuine",
	}

	expectedSeed, _ := hex.DecodeString("be9ac19cea2f3552e16b1391db6a1b4f057a53cfe6831734" +
		"b96e9c85739697fd5b7abe63f16eafa093af4c26eed99198c722c362b743f24fcbf8dbbd198d903e")
	seed, err := prkg.DictionaryEnglish.Seed(mnemonic, "")
	assert.NoError(t, err)
	assert.Equal(t, [64]byte(expectedSeed), seed)

	expectedSeed, _ = hex.DecodeString("f02d1f7b01c9142a1226d87f52dd15010aadcbe96257eded4" +
		"396e64c7d86d7ecc07820e635fab3e17c3f3afdf975b36e8acfee8d0ebfe1b4f6ee28f0da39575b")
	seed, err = prkg.DictionaryEnglish.Seed(mnemonic, "secret")
	assert.NoError(t, err)
	assert.Equal(t, [64]byte(expectedSeed), seed)

	seed, err = prkg.DictionaryEnglish.Seed(mnemonic[1:], "")
	assert.Error(t, err)
	assert.Equal(t, [64]byte{}, seed)

	fmt.Println("Seed-Length:", len(seed))
}
