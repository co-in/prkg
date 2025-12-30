package prkg

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var regexpPath = regexp.MustCompile(`^x/(\d+)/(\d+)/(\d+)/(\d+)$`)

type Path [4]uint32

func ParsePath(path string) (Path, error) {
	parts := regexpPath.FindAllStringSubmatch(path, -1)
	if len(parts) != 1 {
		return Path{}, errors.New("invalid path format")
	}

	m := Path{}

	for i := 1; i <= 4; i++ {
		v := parts[0][i]
		u, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return m, fmt.Errorf("invalid path part %s: %w", v, err)
		}

		if u == 0 {
			return m, errors.New("path part cannot be zero")
		}

		m[i-1] = uint32(u)
	}

	return m, nil
}

func NewPath(coin, wallet, kind, index uint32) Path {
	return Path{
		coin,
		wallet,
		kind,
		index,
	}
}

func (m *Path) String() string {
	return fmt.Sprintf("x/%d/%d/%d/%d", m[0], m[1], m[2], m[3])
}

func (m *Path) SetIndex(value uint32) {
	m[3] = value
}

func (m *Path) SetKind(value uint32) {
	m[2] = value
}
