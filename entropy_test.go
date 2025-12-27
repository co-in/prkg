package prkg_test

import (
	"encoding/hex"
	"testing"

	"github.com/co-in/prkg"
	"github.com/stretchr/testify/assert"
)

func TestEntropy(t *testing.T) {
	expectedEntropy, _ := hex.DecodeString("67776fe4e70fe34fbf9f1432546ff0f3d5e4b9d8")
	bytes, err := prkg.DictionaryEnglish.Entropy([]string{
		"guess", "rocket", "weird", "sock", "wreck", "pond",
		"wrist", "tip", "crane", "pet", "wire", "tray",
		"furnace", "friend", "genuine",
	})
	assert.NoError(t, err)
	assert.Equal(t, expectedEntropy, bytes)

	entropy, err := prkg.NewEntropy(prkg.Mnemonic15)
	assert.NoError(t, err)
	assert.Len(t, entropy, 20)

	entropy, err = prkg.NewEntropy(prkg.Mnemonic24)
	assert.NoError(t, err)
	assert.Len(t, entropy, 32)
}
