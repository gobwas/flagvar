package flagvar

import (
	"time"
)

type TimeValue struct {
	P      *time.Time
	Layout string
}

var zeroTime time.Time

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
	if t.P == nil {
		return zeroTime.Format(t.Layout)
	}
	return t.P.Format(t.Layout)
}

func (t *TimeValue) Get() interface{} {
	if t.P == nil {
		return nil
	}
	return *t.P
}
