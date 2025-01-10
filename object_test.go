package libjson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectAtom(t *testing.T) {
	input := []struct {
		inp      string
		path     string
		expected any
	}{
		{"12", ".", 12},
		{`"str"`, ".", "str"},
		{"true", ".", true},
		{"false", ".", false},
		{"null", ".", nil},
		{`{"key": "value"}`, ".key", "value"},
		{`{ "hello": {"world": ["hi"] } }`, ".hello.world.0", "hi"},
	}
	for _, i := range input {
		t.Run(i.inp+i.path, func(t *testing.T) {
			obj, err := New([]byte(i.inp))
			assert.NoError(t, err)
			assert.NotNil(t, obj)
			out, err := obj.get(i.path)
			assert.NoError(t, err)
			assert.EqualValues(t, i.expected, out)
		})
	}
}

// This tests the example in the readme, always copy from here to the readme
func TestObjectReadme(t *testing.T) {
	input := `{ "hello": {"world": ["hi"] } }`
	jsonObj, _ := New([]byte(input)) // or libjson.NewReader(r io.Reader)

	// accessing values
	// fmt.Println(Get[string](&jsonObj, ".hello.world.0")) // hi
	val, err := Get[string](&jsonObj, ".hello.world.0")
	assert.NoError(t, err)
	assert.EqualValues(t, "hi", val)
}

func TestStandardFail(t *testing.T) {
	input := []string{
		`{"a":"b"}/**/`,
		`{"a":"b"}/**//`,
		`{"a":"b"}//`,
		`{"a":"b"}/`,
		`{"a":"b"}#`,
	}
	for _, i := range input {
		t.Run(i, func(t *testing.T) {
			in := []byte(i)
			p := parser{l: lexer{data: in}}
			_, err := p.parse(in)
			assert.Error(t, err)
		})
	}
}
