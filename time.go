package flagvar

import (
	"time"
)

type TimeValue struct {
	P      *time.Time
	Layout string
}

func Time(p *time.Time, layout string) *TimeValue {
	return &TimeValue{
		P:      p,
		Layout: layout,
	}
}

func (t *TimeValue) Set(s string) error {
	v, err := time.Parse(t.Layout, s)
	if err != nil {
		return err
	}
	(*t.P) = v
	return nil
}

func (t *TimeValue) String() string {
	return t.P.Format(t.Layout)
}
