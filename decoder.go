package stardecoder

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"go.starlark.net/starlark"
)

var (
	errUnsupportedType = errors.New("decoder: unsupported starlark type")
)

// DecodeKWargs decodes Starlark keyword arguments into the value pointed by g.
// If g is a struct, it uses its fields names by default to find the arguments
// and do the mapping. The struct field name can be redefined through struct tags
// prefixed with "starlark". Example:
//
//	type Repo struct {
//		SourceFile string `starlark:"file"`
//	}
//
func DecodeKWargs(kwargs []starlark.Tuple, g interface{}) error {
	structDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "starlark",
		Result:  g,
	})
	if err != nil {
		return err
	}

	input := make(map[string]interface{}, len(kwargs))
	for _, arg := range kwargs {
		name, argValue := arg[0].(starlark.String).GoString(), arg[1]
		value, err := Decode(argValue)
		if err != nil {
			return err
		}
		input[name] = value
	}
	return structDecoder.Decode(input)
}

// Decode recursively decodes the starlark value v into an equivalent shaped Go value g, made up
// of only primitive types.
//
// The following are the conversions made from Starlark types to Go types:
//
//	None   → nil
//	Bool   → bool
//	Int    → int
//	Float  → float64
//	String → string
//	List   → []interface{}
//	Tuple  → []interface{}
//	Dict   → map[interface{}]interface{}{}
//
func Decode(v starlark.Value) (g interface{}, err error) {
	switch x := v.(type) {
	case starlark.NoneType:
		g = nil
	case starlark.Bool:
		g = x.Truth() == starlark.True
	case starlark.Int:
		g, err = starlark.AsInt32(x)
	case starlark.Float:
		g, _ = starlark.AsFloat(x)
	case starlark.String:
		g = x.GoString()
	case *starlark.List, starlark.Tuple:
		var index int
		var item starlark.Value

		slice := make([]interface{}, x.(starlark.Sequence).Len())

		iter := x.(starlark.Sequence).Iterate()
		defer iter.Done()

		for iter.Next(&item) {
			slice[index], err = Decode(item)
			if err != nil {
				return
			}
			index++
		}
		g = slice
	case *starlark.Dict:
		var key, value starlark.Value
		hmap := map[interface{}]interface{}{}

		iter := x.Iterate()
		defer iter.Done()

		for iter.Next(&key) {
			var k interface{}

			k, err = Decode(key)
			if err != nil {
				return
			}

			value, _, err = x.Get(key)
			if err != nil {
				return
			}

			hmap[k], err = Decode(value)
			if err != nil {
				return
			}
			g = hmap
		}
	default:
		err = errUnsupportedType
	}

	return
}
