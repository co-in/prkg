package prkg

import (
	"crypto/rand"
	"errors"
	"math/big"
)

var (
	Mnemonic12 = entropyConfig{
		bytes: 16,
		mask:  big.NewInt(15),
		shift: big.NewInt(16),
	}
	Mnemonic15 = entropyConfig{
		bytes: 20,
		mask:  big.NewInt(31),
		shift: big.NewInt(8),
	}
	Mnemonic18 = entropyConfig{
		bytes: 24,
		mask:  big.NewInt(63),
		shift: big.NewInt(4),
	}
	Mnemonic21 = entropyConfig{
		bytes: 28,
		mask:  big.NewInt(127),
		shift: big.NewInt(2),
	}
	Mnemonic24 = entropyConfig{
		bytes: 32,
		mask:  big.NewInt(255),
		shift: big.NewInt(0),
	}
)

var wordLengthMapping = map[int]entropyConfig{
	12: Mnemonic12,
	15: Mnemonic15,
	18: Mnemonic18,
	21: Mnemonic21,
	24: Mnemonic24,
}

type entropyConfig struct {
	mask  *big.Int
	shift *big.Int
	bytes int
}

func (m entropyConfig) Size() int {
	return (m.bytes * 264) / (352)
}

func EntropyFromSize(wordLength int) ([]byte, error) {
	cfg, ok := wordLengthMapping[wordLength]
	if !ok {
		return nil, errors.New("invalid word length")
	}

	return NewEntropy(cfg)
}

func NewEntropy(e entropyConfig) ([]byte, error) {
	entropy := make([]byte, e.bytes)
	_, err := rand.Read(entropy)

	return entropy, err
}
