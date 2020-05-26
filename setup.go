package flagvar

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type SetupOption interface {
	configureExporter(*exporter)
}

type exporter struct {
	recursion bool
	lispCase  bool
}

type WithLispCase bool

func (x WithLispCase) configureExporter(e *exporter) {
	e.lispCase = bool(x)
}

type WithRecursion bool

func (x WithRecursion) configureExporter(e *exporter) {
	e.recursion = bool(x)
}

func Struct(flag *flag.FlagSet, x interface{}, opts ...SetupOption) error {
	e := exporter{
		lispCase: true,
	}
	for _, opt := range opts {
		opt.configureExporter(&e)
	}

	v := reflect.ValueOf(x)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("can't setup non-addressable value: %s", v.Type())
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("can't setup non-struct value: %s", v.Type())
	}
	t := v.Type()
	return e.setupStruct(flag, "", t, v)
}

func (e exporter) setupStruct(flag *flag.FlagSet, parent string, t reflect.Type, v reflect.Value) error {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		v := v.Field(i)
		if !v.CanSet() {
			continue
		}
		name := f.Name
		if e.lispCase {
			name = lispCase(name)
		}
		if parent != "" {
			name = parent + "." + name
		}
		if e.recursion && v.Kind() == reflect.Struct {
			err := e.setupStruct(flag, name, v.Type(), v)
			if err != nil {
				return err
			}
			continue
		}
		if err := e.setup(flag, name, "", v); err != nil {
			return fmt.Errorf(
				"set up %s's field %s error: %v",
				v.Type(), f.Name, err,
			)
		}
	}
	return nil
}

const bitSize = 32 << (^uint(0) >> 63)

var timeDuration = reflect.TypeOf(time.Duration(0))

func (e exporter) setup(fs *flag.FlagSet, name, usage string, v reflect.Value) error {
	val := typeValue(v)
	if val == nil {
		val = kindValue(v)
	}
	if val != nil {
		if f := fs.Lookup(name); f != nil {
			return fmt.Errorf("can't define flag %q: already exists", name)
		}
		fs.Var(val, name, usage)
	}
	return nil
}

func typeValue(v reflect.Value) flag.Value {
	switch v.Type() {
	case timeDuration:
		return durationValue{x: v}
	}
	return nil
}

func kindValue(v reflect.Value) flag.Value {
	switch v.Kind() {
	case reflect.Bool:
		return boolValue{x: v}
	case reflect.Int:
		return intValue{x: v, bits: bitSize, base: 10}
	case reflect.Int8:
		return intValue{x: v, bits: 8, base: 10}
	case reflect.Int16:
		return intValue{x: v, bits: 16, base: 10}
	case reflect.Int32:
		return intValue{x: v, bits: 32, base: 10}
	case reflect.Int64:
		return intValue{x: v, bits: 64, base: 10}
	case reflect.Uint:
		return uintValue{x: v, bits: bitSize, base: 10}
	case reflect.Uint8:
		return uintValue{x: v, bits: 8, base: 10}
	case reflect.Uint16:
		return uintValue{x: v, bits: 16, base: 10}
	case reflect.Uint32:
		return uintValue{x: v, bits: 32, base: 10}
	case reflect.Uint64:
		return uintValue{x: v, bits: 64, base: 10}
	case reflect.Float32:
		return floatValue{x: v, bits: 32}
	case reflect.Float64:
		return floatValue{x: v, bits: 64}
	//case reflect.Array:
	//case reflect.Map:
	//case reflect.Slice:
	case reflect.String:
		return stringValue{x: v}
	default:
	}
	return nil
}

type durationValue struct {
	x reflect.Value
}

func (v durationValue) Set(s string) error {
	x, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	v.x.SetInt(int64(x))
	return nil
}
func (v durationValue) String() string {
	if !v.x.IsValid() {
		return "<zero>"
	}
	return v.x.Interface().(time.Duration).String()
}

type floatValue struct {
	x    reflect.Value
	bits int
}

func (v floatValue) Set(s string) error {
	x, err := strconv.ParseFloat(s, v.bits)
	if err != nil {
		return err
	}
	v.x.SetFloat(x)
	return nil
}

func (v floatValue) String() string {
	if !v.x.IsValid() {
		return "<zero>"
	}
	return strconv.FormatFloat(v.x.Float(), 'f', -1, v.bits)
}

type stringValue struct {
	x reflect.Value
}

func (v stringValue) Set(s string) error {
	v.x.Set(reflect.ValueOf(s))
	return nil
}

func (v stringValue) String() string {
	if !v.x.IsValid() {
		return "<zero>"
	}
	return v.x.String()
}

type boolValue struct {
	x reflect.Value
}

func (v boolValue) Set(s string) error {
	x, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	v.x.Set(reflect.ValueOf(x))
	return nil
}

func (v boolValue) String() string {
	if !v.x.IsValid() {
		return "<zero>"
	}
	return strconv.FormatBool(v.x.Bool())
}

type intValue struct {
	x    reflect.Value
	bits int
	base int
}

func (v intValue) Set(s string) error {
	x, err := strconv.ParseInt(s, v.base, v.bits)
	if err != nil {
		return err
	}
	v.x.SetInt(x)
	return nil
}

func (v intValue) String() string {
	if !v.x.IsValid() {
		return "<zero>"
	}
	return strconv.FormatInt(v.x.Int(), v.base)
}

type uintValue struct {
	x    reflect.Value
	bits int
	base int
}

func (v uintValue) Set(s string) error {
	x, err := strconv.ParseUint(s, v.base, v.bits)
	if err != nil {
		return err
	}
	v.x.SetUint(x)
	return nil
}

func (v uintValue) String() string {
	if !v.x.IsValid() {
		return "<zero>"
	}
	return strconv.FormatUint(v.x.Uint(), v.base)
}
