package prkg

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

type Dictionary struct {
	wordList []string
	wordMap  map[string]int
}

func NewDictionary(wordListRaw string) Dictionary {
	m := Dictionary{
		wordList: strings.Split(strings.TrimSpace(wordListRaw), "\n"),
	}

	m.wordMap = make(map[string]int, len(m.wordList))
	for i, v := range m.wordList {
		m.wordMap[v] = i
	}

	return m
}

func (m Dictionary) Words() []string {
	return m.wordList
}

func (m Dictionary) Entropy(words []string) (entropy []byte, err error) {
	wLen := len(words)

	if _, ok := wordLengthMapping[wLen]; !ok {
		return nil, fmt.Errorf("invalid word count %d", wLen)
	}

	b := big.NewInt(0)

	for _, v := range words {
		index, ok := m.wordMap[v]
		if !ok {
			return nil, fmt.Errorf("word `%v` not found in reverse map", v)
		}

		var wordBytes [2]byte
		binary.BigEndian.PutUint16(wordBytes[:], uint16(index))

		b = b.Mul(b, shift11BitsMask)
		b = b.Or(b, big.NewInt(0).SetBytes(wordBytes[:]))
	}

	ec := wordLengthMapping[wLen]
	expectedSum := big.NewInt(0).And(b, ec.mask)

	b.Div(b, big.NewInt(0).Add(ec.mask, bigOne))
	entropy = padBytes(b.Bytes(), ec.bytes)

	entropySum := big.NewInt(int64(sha256.Sum256(entropy)[0]))
	if l := wLen; l != 24 {
		checksumShift := wordLengthMapping[l].shift
		entropySum.Div(entropySum, checksumShift)
	}
	if entropySum.Cmp(expectedSum) != 0 {
		return nil, fmt.Errorf("checksum mismatch")
	}

	return entropy, nil
}

func (m Dictionary) Mnemonic(entropy []byte) (words []string, err error) {
	entropyBitLength := len(entropy) * 8
	checksumBitLength := entropyBitLength / 32
	sentenceLength := (entropyBitLength + checksumBitLength) / 11

	ec, ok := wordLengthMapping[sentenceLength]
	if !ok || ec.bytes*8 != entropyBitLength {
		return nil, fmt.Errorf("invalid entropy bit length %d", entropyBitLength)
	}

	entropy, err = addChecksum(entropy)
	if err != nil {
		return nil, err
	}

	entropyInt := new(big.Int).SetBytes(entropy)
	words = make([]string, sentenceLength)
	word := big.NewInt(0)

	for i := sentenceLength - 1; i >= 0; i-- {
		word.And(entropyInt, last11BitsMask)
		entropyInt.Div(entropyInt, shift11BitsMask)
		wordBytes := padBytes(word.Bytes(), 2)
		words[i] = m.wordList[binary.BigEndian.Uint16(wordBytes)]
	}

	return words, nil
}

func (m Dictionary) Seed(words []string, password string) (seed [64]byte, err error) {
	if _, ok := wordLengthMapping[len(words)]; !ok {
		return seed, errors.New("invalid word count")
	}

	for _, v := range words {
		if _, ok := m.wordMap[v]; !ok {
			return seed, fmt.Errorf("word `%v` not found in dictionary", v)
		}
	}

	var seedD []byte
	seedD, err = pbkdf2.Key(sha512.New, strings.Join(words, " "), []byte("mnemonic"+password), 2048, 64)
	if err != nil {
		return seed, err
	}

	return [64]byte(seedD), err
}

var (
	last11BitsMask  = big.NewInt(2047)
	shift11BitsMask = big.NewInt(2048)
	bigOne          = big.NewInt(1)
	bigTwo          = big.NewInt(2)
)

func addChecksum(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	firstChecksumByte := hash[0]
	checksumBitLength := uint(len(data) / 4)
	dataBigInt := new(big.Int).SetBytes(data)

	for i := uint(0); i < checksumBitLength; i++ {
		dataBigInt.Mul(dataBigInt, bigTwo)

		if firstChecksumByte&(1<<(7-i)) > 0 {
			dataBigInt.Or(dataBigInt, bigOne)
		}
	}

	return dataBigInt.Bytes(), nil
}

func padBytes(slice []byte, length int) []byte {
	offset := length - len(slice)
	if offset <= 0 {
		return slice
	}

	newSlice := make([]byte, length)
	copy(newSlice[offset:], slice)

	return newSlice
}
