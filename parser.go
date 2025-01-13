package libjson

import (
	"fmt"
	"strconv"
	"unsafe"
)

type parser struct {
	l       lexer
	c       <-chan token
	cur_tok token
}

func (p *parser) advance() {
	p.cur_tok = <-p.c
}

// parses toks into a valid json representation, thus the return type can be
// either map[string]any, []any, string, nil, false, true or a number
func (p *parser) parse() (any, error) {
	p.advance()
	if val, err := p.expression(); err != nil {
		return nil, err
	} else {
		if p.cur_tok.Type != t_eof && p.cur_tok.Type > t_string {
			return nil, fmt.Errorf("Unexpected non-whitespace character(s) (%s) after JSON data", tokennames[p.cur_tok.Type])
		}
		return val, nil
	}
}

func (p *parser) expression() (any, error) {
	if p.cur_tok.Type == t_left_curly {
		return p.object()
	} else if p.cur_tok.Type == t_left_braket {
		return p.array()
	} else {
		return p.atom()
	}
}

func (p *parser) object() (map[string]any, error) {
	if p.cur_tok.Type != t_left_curly {
		return nil, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_left_curly])
	}
	p.advance()

	m := make(map[string]any, 4)

	if p.cur_tok.Type == t_right_curly {
		p.advance()
		return m, nil
	}

	for p.cur_tok.Type != t_eof && p.cur_tok.Type != t_right_curly {
		if len(m) > 0 {
			if p.cur_tok.Type != t_comma {
				return nil, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_comma])
			}
			p.advance()
		}

		if p.cur_tok.Type != t_string {
			return nil, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_string])
		}
		in := p.l.data[p.cur_tok.Start:p.cur_tok.End]
		key := *(*string)(unsafe.Pointer(&in))
		p.advance()

		if p.cur_tok.Type != t_colon {
			return nil, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_colon])
		}
		p.advance()

		val, err := p.expression()
		if err != nil {
			return nil, err
		}

		// TODO:  think about activating a uniqueness check for object keys,
		// would add an other hashing and a branch for each object key parsed.
		//
		// if _, ok := m[key]; ok {
		// 	return nil, fmt.Errorf("Key %q is already set in this object", key)
		// }

		m[key] = val
	}

	if p.cur_tok.Type != t_right_curly {
		return nil, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_right_curly])
	}
	p.advance()

	return m, nil
}

func (p *parser) array() ([]any, error) {
	if p.cur_tok.Type != t_left_braket {
		return nil, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_left_braket])
	}
	p.advance()

	if p.cur_tok.Type == t_right_braket {
		p.advance()
		return []any{}, nil
	}

	a := make([]any, 0, 8)

	for p.cur_tok.Type != t_eof && p.cur_tok.Type != t_right_braket {
		if len(a) > 0 {
			if p.cur_tok.Type != t_comma {
				return nil, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_comma])
			}
			p.advance()
		}
		node, err := p.expression()
		if err != nil {
			return nil, err
		}
		a = append(a, node)
	}

	if p.cur_tok.Type != t_right_braket {
		return nil, fmt.Errorf("Unexpected %q at this position, expected %q", tokennames[p.cur_tok.Type], tokennames[t_right_braket])
	}

	p.advance()
	return a, nil
}

func (p *parser) atom() (any, error) {
	var r any
	switch p.cur_tok.Type {
	case t_string:
		in := p.l.data[p.cur_tok.Start:p.cur_tok.End]
		r = *(*string)(unsafe.Pointer(&in))
	case t_number:
		in := p.l.data[p.cur_tok.Start:p.cur_tok.End]
		raw := *(*string)(unsafe.Pointer(&in))
		number, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return empty, fmt.Errorf("Invalid floating point number %q: %w", raw, err)
		}
		r = number
	case t_true:
		r = true
	case t_false:
		r = false
	case t_null:
		r = nil
	default:
		return nil, fmt.Errorf("Unexpected %q at this position, expected any of: string, number, true, false or null", tokennames[p.cur_tok.Type])
	}
	p.advance()
	return r, nil
}
