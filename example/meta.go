package example

import (
	"fmt"
	"strings"
)

type Meta struct {
	params []interface{}
	format string
}

func (m *Meta) GetKey() (s string) {

	stringParams := make([]string, len(m.params))
	for i, param := range m.params {
		switch param := param.(type) {
		case string:
			stringParams[i] = param
		default:
			stringParams[i] = fmt.Sprint(param)
		}
	}

	s = strings.Join(stringParams, ":")
	s = fmt.Sprintf(m.format, s)

	return
}

type ProtoWrapper interface {
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
	GetKey() (s string)
}

var _ ProtoWrapper = &Person{}
