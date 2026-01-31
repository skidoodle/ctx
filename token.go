package main

import (
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

type TokenCounter struct {
	w     io.Writer
	Count int64
	Err   error

	leftover []byte
	inWord   bool
	wordLen  int
}

func (tc *TokenCounter) Write(p []byte) (int, error) {
	if tc.Err != nil {
		return 0, tc.Err
	}

	data := p
	if len(tc.leftover) > 0 {
		data = make([]byte, len(tc.leftover)+len(p))
		copy(data, tc.leftover)
		copy(data[len(tc.leftover):], p)
	}

	totalProcessed := 0

	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)

		if r == utf8.RuneError && size == 1 {
			if len(data) < utf8.UTFMax {
				tc.leftover = data
				break
			}
		}

		data = data[size:]
		totalProcessed += size
		tc.leftover = nil

		isSpace := unicode.IsSpace(r)
		isAlpha := unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_'

		if isAlpha {
			if !tc.inWord {
				tc.Count++
				tc.inWord = true
				tc.wordLen = 1
			} else {
				tc.wordLen++
				if tc.wordLen > 4 {
					tc.Count++
					tc.wordLen = 1
				}
			}
		} else if isSpace {
			tc.inWord = false
			tc.wordLen = 0
		} else {
			tc.Count++
			tc.inWord = false
			tc.wordLen = 0
		}
	}

	var n int
	if tc.w != nil {
		n, tc.Err = tc.w.Write(p)
	} else {
		n = len(p)
	}

	return n, tc.Err
}

func (tc *TokenCounter) WriteByte(c byte) error {
	_, err := tc.Write([]byte{c})
	return err
}

func (tc *TokenCounter) Printf(format string, a ...any) {
	if tc.Err != nil {
		return
	}
	_, tc.Err = fmt.Fprintf(tc, format, a...)
}

func (tc *TokenCounter) Println(a ...any) {
	if tc.Err != nil {
		return
	}
	_, tc.Err = fmt.Fprintln(tc, a...)
}
