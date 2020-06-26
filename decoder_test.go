package stardecoder

import (
	"testing"

	"github.com/hooklift/assert"
	"go.starlark.net/starlark"
)

func TestDecode(t *testing.T) {
	dict := starlark.NewDict(1)
	err := dict.SetKey(starlark.String("mykey"), starlark.String("myvalue"))
	assert.Ok(t, err)
	dict.SetKey(starlark.String("list"), starlark.Tuple([]starlark.Value{starlark.MakeInt(99)}))

	tests := []struct {
		x starlark.Value
		y interface{}
	}{
		{starlark.None, nil},
		{starlark.Bool(true), true},
		{starlark.MakeInt(1), 1},
		{starlark.Float(1.23), 1.23},
		{starlark.String("hola"), "hola"},
		{starlark.NewList([]starlark.Value{
			starlark.String("a"),
			starlark.String("b"),
			starlark.String("c"),
		}), []interface{}{"a", "b", "c"}},
		{starlark.NewList([]starlark.Value{
			starlark.Bool(true),
			starlark.None,
			starlark.String("c"),
		}), []interface{}{true, nil, "c"}},
		{starlark.Tuple([]starlark.Value{
			starlark.MakeInt(234567890),
			starlark.Float(1.2),
			starlark.NewList([]starlark.Value{
				starlark.String("nested list"),
			}),
		}), []interface{}{234567890, 1.2, []interface{}{"nested list"}}},
		{
			dict,
			map[interface{}]interface{}{"list": []interface{}{99}, "mykey": "myvalue"},
		},
	}

	for _, tt := range tests {
		value, err := Decode(tt.x)
		assert.Ok(t, err)
		assert.Equals(t, tt.y, value)
	}
}
