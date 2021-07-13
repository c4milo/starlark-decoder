package stardecoder

import (
	"testing"

	"github.com/hooklift/assert"
	"go.starlark.net/starlark"
)

func TestDecode(t *testing.T) {
	dict := starlark.NewDict(1)
	dict.SetKey(starlark.String("mykey"), starlark.String("myvalue"))
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

func TestDecodeKWargs(t *testing.T) {
	dict := starlark.NewDict(1)
	dict.SetKey(starlark.String("name"), starlark.String("pablo"))
	dict.SetKey(starlark.String("age"), starlark.MakeInt(36))
	dict.SetKey(starlark.String("hobbies"), starlark.Tuple([]starlark.Value{
		starlark.String("movies"),
		starlark.String("horse riding"),
		starlark.String("car racing"),
	}))
	dict.SetKey(starlark.String("myfn"), &starlark.Function{})

	tests := []struct {
		x  []starlark.Tuple
		y  interface{}
		yy interface{}
	}{
		{
			[]starlark.Tuple{
				[]starlark.Value{
					starlark.String("url"), starlark.String("http://slack.com"),
				},
			},
			struct{ URL string }{},
			struct{ URL string }{URL: "http://slack.com"},
		},
		{
			[]starlark.Tuple{
				[]starlark.Value{
					starlark.String("url"), starlark.String("http://slack.com"),
				},
			},
			struct {
				Blah string `starlark:"url"`
			}{},
			struct {
				Blah string `starlark:"url"`
			}{Blah: "http://slack.com"},
		},
		{
			[]starlark.Tuple{
				[]starlark.Value{
					starlark.String("numbers"),
					starlark.NewList([]starlark.Value{
						starlark.MakeInt(1),
						starlark.MakeInt(2),
						starlark.MakeInt(3),
					}),
				},
			},
			struct {
				Numbers []int
			}{},
			struct {
				Numbers []int
			}{Numbers: []int{1, 2, 3}},
		},
		{
			[]starlark.Tuple{
				[]starlark.Value{
					starlark.String("contact"), dict,
				},
			},
			struct {
				Contact struct {
					Name    string
					Age     int
					Hobbies []string
				}
			}{},
			struct {
				Contact struct {
					Name    string
					Age     int
					Hobbies []string
				}
			}{
				Contact: struct {
					Name    string
					Age     int
					Hobbies []string
				}{
					Name: "pablo", Age: 36,
					Hobbies: []string{"movies", "horse riding", "car racing"},
				},
			},
		},
	}

	for _, tt := range tests {
		err := DecodeKWargs(tt.x, &tt.y)
		assert.Ok(t, err)
		assert.Equals(t, tt.yy, tt.y)
	}
}
