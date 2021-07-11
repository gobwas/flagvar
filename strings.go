package flagvar

import (
	"fmt"
	"strings"
)

type StringsValue struct {
	P         *[]string
	Separator string
}

func Strings(p *[]string, sep string) *StringsValue {
	return &StringsValue{
		P:         p,
		Separator: sep,
	}
}

func (s *StringsValue) Set(v string) error {
	(*s.P) = append((*s.P), strings.Split(v, s.Separator)...)
	return nil
}

func (s *StringsValue) String() string {
	if s.P == nil {
		return ""
	}
	if sep := s.Separator; sep != "" {
		return strings.Join((*s.P), sep)
	}
	return fmt.Sprintf("%s", (*s.P))
}
