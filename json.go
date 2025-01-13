package libjson

import (
	"io"
)

const chanSize = 50_000

func NewReader(r io.Reader) (JSON, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return JSON{}, err
	}
	return New(data)
}

func New(data []byte) (JSON, error) {
	l := lexer{data: data}
	c := make(chan token, chanSize)
	var lexerErr error
	go func() {
		for {
			if tok, err := l.next(); err == nil {
				if tok.Type == t_eof {
					break
				}
				c <- tok
			} else {
				lexerErr = err
				break
			}
		}
		close(c)
		c = nil
	}()
	p := parser{l: l, c: c}
	obj, err := p.parse()
	if lexerErr != nil {
		return JSON{}, lexerErr
	}
	if err != nil {
		return JSON{}, err
	}
	return JSON{obj}, nil
}
