package libjson

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexerWhitespace(t *testing.T) {
	json := "\n\r\t      "
	l := lexer{}
	toks, err := l.lex(strings.NewReader(json))
	assert.Error(t, err)
	assert.Empty(t, toks)
}

func TestLexerStructure(t *testing.T) {
	json := "{}[],:"
	l := lexer{}
	toks, err := l.lex(strings.NewReader(json))
	if err != nil {
		t.Error(err)
	}
	tList := []token{
		{Type: t_left_curly},
		{Type: t_right_curly},
		{Type: t_left_braket},
		{Type: t_right_braket},
		{Type: t_comma},
		{Type: t_colon},
	}
	assert.EqualValues(t, tList, toks)
}

func TestLexerAtoms(t *testing.T) {
	json := `
    "string""" "🤣"
    true false null
    1 0 12.5 1e15 -1929 -0
    -1.4E+5 -129.1928e-19028
    `
	l := lexer{}
	toks, err := l.lex(strings.NewReader(json))
	assert.NoError(t, err)
	tList := []struct {
		Type t_json
		Val  string
	}{
		{Type: t_string, Val: "string"},
		{Type: t_string, Val: ""},
		{Type: t_string, Val: "🤣"},
		{Type: t_true},
		{Type: t_false},
		{Type: t_null},
		{Type: t_number, Val: "1"},
		{Type: t_number, Val: "0"},
		{Type: t_number, Val: "12.5"},
		{Type: t_number, Val: "1e15"},
		{Type: t_number, Val: "-1929"},
		{Type: t_number, Val: "-0"},
		{Type: t_number, Val: "-1.4E+5"},
		{Type: t_number, Val: "-129.1928e-19028"},
	}
	for i, tok := range tList {
		got := toks[i]
		wanted := tok
		assert.EqualValues(t, wanted.Type, got.Type)
		assert.EqualValues(t, wanted.Val, string(json[got.Start:got.End]))
	}
}

func TestLexer(t *testing.T) {
	json := `
    {
        "key": "value",
        "arrayOfDataTypes": ["string", 1, true, false, null],
        "subobject": { "key": "value" }
    }
    `
	l := lexer{}
	toks, err := l.lex(strings.NewReader(json))
	if err != nil {
		t.Error(err)
	}

	tList := []struct {
		Type t_json
		Val  string
	}{
		{Type: t_left_curly},

		{Type: t_string, Val: "key"},
		{Type: t_colon},
		{Type: t_string, Val: "value"},
		{Type: t_comma},

		{Type: t_string, Val: "arrayOfDataTypes"},
		{Type: t_colon},
		{Type: t_left_braket},
		{Type: t_string, Val: "string"},
		{Type: t_comma},
		{Type: t_number, Val: "1"},
		{Type: t_comma},
		{Type: t_true},
		{Type: t_comma},
		{Type: t_false},
		{Type: t_comma},
		{Type: t_null},
		{Type: t_right_braket},
		{Type: t_comma},

		{Type: t_string, Val: "subobject"},
		{Type: t_colon},
		{Type: t_left_curly},
		{Type: t_string, Val: "key"},
		{Type: t_colon},
		{Type: t_string, Val: "value"},
		{Type: t_right_curly},
		{Type: t_right_curly},
	}
	for i, tok := range tList {
		got := toks[i]
		assert.Equal(t, tok.Type, got.Type)
		if tok.Val != "" {
			assert.EqualValues(t, tok.Val, string(json[got.Start:got.End]))
		}
	}
}

func TestLexerFail(t *testing.T) {
	input := []string{
		"",
		`"`,
		"'",
		`0xFF`,
		string([]byte{0x0C}),
		`{"test": 'value'}`,
		"🤣",
		`{"a":"b"}/**/`,
	}
	for _, in := range input {
		t.Run(in, func(t *testing.T) {
			l := &lexer{}
			toks, err := l.lex(strings.NewReader(in))
			assert.Error(t, err)
			assert.Empty(t, toks)
		})
	}
}
