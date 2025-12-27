package prkg_test

import (
	"encoding/hex"
	"testing"

	"github.com/co-in/prkg"
	"github.com/stretchr/testify/assert"
)

func TestDeterministic(t *testing.T) {
	seed, _ := hex.DecodeString("f02d1f7b01c9142a1226d87f52dd15010aadcbe96257eded4396e64c7d86d7ecc07820e635fab3e17c3f3afdf975b36e8acfee8d0ebfe1b4f6ee28f0da39575b")
	key, _ := hex.DecodeString("4bd1067c24f96e5344447e71708a21a6a2487f123fb7a6da3e93b09a0c96c91e")

	hd1, err := prkg.NewDK([64]byte(seed))
	assert.NoError(t, err)
	assert.NotNil(t, hd1)
	key1, err := hd1.Jump(44, 0, 0, 1)
	assert.NoError(t, err)
	assert.Equal(t, key, key1)

	hd2, err := prkg.NewDK([64]byte(seed))
	assert.NoError(t, err)
	assert.NotNil(t, hd2)
	key2, err := hd2.Jump(44, 0, 0, 1)
	assert.NoError(t, err)
	assert.Equal(t, key, key2)
}
