package flagvar

type StringsValue struct {
	P *[]string
}

func Strings(p *[]string) *StringsValue {
	return &StringsValue{
		P: p,
	}
}

func (s *StringsValue) Set(v string) error {
	(*s.P) = append((*s.P), v)
	return nil
}
