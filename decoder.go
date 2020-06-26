package stardecoder

import (
	"errors"

	"go.starlark.net/starlark"
)

var (
	errUnsupportedType = errors.New("decoder: unsupported starlark type")
)

// Decode recursively decodes the starlark value v into an equivalent shaped Go value g, made up
// of only primitive types.
//
// The following are the conversions made from Starlark types to Go types:
//
// None   → nil
// Bool   → bool
// Int    → int
// Float  → float64
// String → string
// List   → []interface{}
// Tuple  → []interface{}
// Dict   → map[interface{}]interface{}{}
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
	// case *starlark.Set, *starlark.Function, *starlark.Builtin:
	default:
		err = errUnsupportedType
	}

	return
}
