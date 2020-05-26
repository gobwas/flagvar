package flagvar

import (
	"flag"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestStruct(t *testing.T) {
	type Nested struct {
		String string
	}
	type Structure struct {
		Bool     bool
		Int      int
		Int8     int8
		Int16    int16
		Int32    int32
		Int64    int64
		Uint     uint
		Uint8    uint8
		Uint16   uint16
		Uint32   uint32
		Uint64   uint64
		Float32  float32
		Float64  float64
		Duration time.Duration

		//Array   [8]string
		//Map map[string]string
		//Slice   []string
		String string

		Nested Nested
	}
	for _, test := range []struct {
		name string
		args []string
		opts []SetupOption
		in   Structure
		exp  Structure
	}{
		{
			opts: []SetupOption{
				WithRecursion(true),
			},
			args: []string{
				"-string", "string",
				"-nested.string", "nested",
				"-bool", "true",
				"-int", "42",
				"-uint", "42",
				"-int8", strconv.Itoa(math.MinInt8),
				"-int16", strconv.Itoa(math.MinInt16),
				"-int32", strconv.Itoa(math.MinInt32),
				"-int64", strconv.Itoa(math.MinInt64),
				"-uint8", strconv.FormatUint(math.MaxUint8, 10),
				"-uint16", strconv.FormatUint(math.MaxUint16, 10),
				"-uint32", strconv.FormatUint(math.MaxUint32, 10),
				"-uint64", strconv.FormatUint(math.MaxUint64, 10),
				"-float32", strconv.FormatFloat(math.MaxFloat32, 'e', -1, 32),
				"-float64", strconv.FormatFloat(math.MaxFloat32, 'e', -1, 64),
				"-duration", "2562047h47m16.854775807s", // MaxInt64.
			},
			exp: Structure{
				String:   "string",
				Bool:     true,
				Int:      42,
				Uint:     42,
				Int8:     math.MinInt8,
				Int16:    math.MinInt16,
				Int32:    math.MinInt32,
				Int64:    math.MinInt64,
				Uint8:    math.MaxUint8,
				Uint16:   math.MaxUint16,
				Uint32:   math.MaxUint32,
				Uint64:   math.MaxUint64,
				Float32:  math.MaxFloat32,
				Float64:  math.MaxFloat32,
				Duration: time.Duration(math.MaxInt64),
				Nested: Nested{
					String: "nested",
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			fs := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
			if err := Struct(fs, &test.in, test.opts...); err != nil {
				t.Fatal(err)
			}
			fs.PrintDefaults()
			if err := fs.Parse(test.args); err != nil {
				t.Fatal(err)
			}
			if act, exp := test.in, test.exp; !cmp.Equal(act, exp) {
				t.Fatalf("%s", cmp.Diff(act, exp))
			}
			t.Logf("%+v", test.in)
		})
	}
}
