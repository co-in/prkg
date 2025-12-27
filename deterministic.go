package prkg

import (
	"crypto/hkdf"
	"crypto/sha512"
	"fmt"
)

const (
	masterKeyHKDFSalt = "key derivation"
	masterKeyHKDFInfo = "hardened HDKF"
)

type DK struct {
	masterKey      []byte
	keyEntropySize uint16
	maxLevel       uint8
}

// DKOption
/*
	- WithDKKeyEntropySize
	- WithDKMaxLevel
*/
type DKOption func(*DK)

func WithDKKeyEntropySize(value uint16) DKOption {
	return func(m *DK) {
		m.keyEntropySize = value
	}
}

func WithDKMaxLevel(value uint8) DKOption {
	return func(m *DK) {
		m.maxLevel = value
	}
}

func NewDK(seed [64]byte, options ...DKOption) (*DK, error) {
	m := &DK{
		keyEntropySize: 32,
		maxLevel:       4,
	}

	for _, opt := range options {
		opt(m)
	}

	if seed == [64]byte{} {
		return nil, fmt.Errorf("empty seed")
	}

	//Vulnerable entropy pruning
	if m.keyEntropySize < 32 {
		return nil, fmt.Errorf("invalid key entropy size: %d", m.keyEntropySize)
	}

	//prevent return part of MasterKey
	if m.maxLevel < 1 {
		return nil, fmt.Errorf("invalid max level: %d", m.maxLevel)
	}

	var err error

	if m.masterKey, err = hkdf.Key(sha512.New, seed[:],
		[]byte(masterKeyHKDFSalt), masterKeyHKDFInfo, int(m.keyEntropySize*2),
	); err != nil {
		return nil, fmt.Errorf("derive master key: %w", err)
	}

	return m, nil
}

func (m *DK) Jump(path ...uint32) (key []byte, err error) {
	if m.masterKey == nil {
		return key, fmt.Errorf("master key is nil")
	}

	pl := len(path)
	if pl == 0 || pl > int(m.maxLevel) {
		return key, fmt.Errorf("invalid path length: %d", pl)
	}

	keyI := m.masterKey

	for _, idx := range path {
		if keyI, err = hkdf.Key(sha512.New,
			keyI[m.keyEntropySize:],
			keyI[:m.keyEntropySize],
			fmt.Sprintf("%d", idx),
			int(m.keyEntropySize*2),
		); err != nil {
			return key, err
		}
	}

	return keyI[:m.keyEntropySize], nil
}
