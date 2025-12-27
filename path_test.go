package prkg_test

import (
	"testing"

	"github.com/co-in/prkg"
	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {
	pathStr := "p/1/11/2/3333"
	path := prkg.Path([]uint32{1, 11, 2, 3333})

	assert.Equal(t, path.String(), pathStr)

	p, err := prkg.ParsePath(pathStr)
	assert.NoError(t, err)
	assert.Equal(t, p, path)
}
