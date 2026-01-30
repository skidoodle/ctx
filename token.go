package main

import (
	"fmt"
	"io"
	"unicode"
)

type TokenCounter struct {
	w       io.Writer
	Count   int64
	Err     error
	inWord  bool
	inSpace bool
}

func (tc *TokenCounter) Write(p []byte) (int, error) {
	if tc.Err != nil {
		return 0, tc.Err
	}

	for _, b := range p {
		r := rune(b)
		isSpace := unicode.IsSpace(r)
		isAlpha := unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_'

		if isAlpha {
			if !tc.inWord {
				tc.Count++
				tc.inWord = true
				tc.inSpace = false
			}
		} else if isSpace {
			if !tc.inSpace {
				tc.Count++
				tc.inSpace = true
				tc.inWord = false
			}
		} else {
			tc.Count++
			tc.inWord = false
			tc.inSpace = false
		}
	}

	n, err := tc.w.Write(p)
	tc.Err = err
	return n, err
}

func (tc *TokenCounter) WriteByte(c byte) error {
	_, err := tc.Write([]byte{c})
	return err
}

func (tc *TokenCounter) Printf(format string, a ...any) {
	if tc.Err != nil {
		return
	}
	_, _ = fmt.Fprintf(tc, format, a...)
}

func (tc *TokenCounter) Println(a ...any) {
	if tc.Err != nil {
		return
	}
	_, _ = fmt.Fprintln(tc, a...)
}
