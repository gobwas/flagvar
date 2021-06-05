package flagvar

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// SetupOption is the interface to configure flags definitions.
type SetupOption interface {
	configureExporter(*exporter)
}

type exporter struct {
	recursion    bool
	lispCase     bool
	setSeparator string
}

// WithLispCase specifies whether flag name should be in lisp-case form.
type WithLispCase bool

func (x WithLispCase) configureExporter(e *exporter) {
	e.lispCase = bool(x)
}

// WithRecursion specifies whether flag definition should be recursive.
type WithRecursion bool

func (x WithRecursion) configureExporter(e *exporter) {
	e.recursion = bool(x)
}

// WithSetSeparator specifies how nested flag sets should be separated.
type WithSetSeparator string

func (x WithSetSeparator) configureExporter(e *exporter) {
	e.setSeparator = string(x)
}

// Struct iterates over given struct x fields and defines its fields as
// appropriate flag variables.
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
			name = parent + e.setSeparator + name
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
func (v durationValue) Get() interface{} {
	return v.x.Interface()
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
func (v floatValue) Get() interface{} {
	f := v.x.Float()
	switch v.x.Kind() {
	case reflect.Float32:
		return float32(f)
	case reflect.Float64:
		return float64(f)
	}
	panic(fmt.Sprintf("unexpected kind: %s", v.x.Kind()))
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
func (v stringValue) Get() interface{} {
	return v.x.Interface()
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
func (v boolValue) Get() interface{} {
	return v.x.Interface()
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
func (v intValue) Get() interface{} {
	return v.x.Interface()
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
func (v uintValue) Get() interface{} {
	return v.x.Interface()
}
