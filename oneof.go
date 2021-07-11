package flagvar

import (
	"fmt"
	"reflect"
	"sort"
)

type Mapping map[string]interface{}

type OneOfValue struct {
	P       interface{}
	Mapping Mapping
}

func OneOf(p interface{}, m Mapping) *OneOfValue {
	val := reflect.ValueOf(p)
	if val.Kind() != reflect.Ptr {
		panic("argument p must be a pointer")
	}
	typ := val.Elem().Type()
	for k, v := range m {
		t := reflect.TypeOf(v)
		if t != typ {
			panic(fmt.Sprintf(
				"value for mapping key %q is %s; but receiver is %s",
				k, t, typ,
			))
		}
	}
	return &OneOfValue{
		P:       p,
		Mapping: m,
	}
}

func (v *OneOfValue) Set(str string) error {
	var keys []string
	for key, val := range v.Mapping {
		if str != key {
			keys = append(keys, key)
			continue
		}
		reflect.ValueOf(v.P).Elem().Set(reflect.ValueOf(val))
		return nil
	}
	sort.Strings(keys)
	return fmt.Errorf("value %q is not one of %q", str, keys)
}

func (v *OneOfValue) String() string {
	return fmt.Sprintf("%s", v.P)
}
