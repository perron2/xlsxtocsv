package main

import (
	"golang.org/x/text/encoding/charmap"
	"io"
	"unicode/utf8"
)

type ansiWriter struct {
	w        io.Writer
	fallback byte
}

func (w ansiWriter) Write(p []byte) (n int, err error) {
	var buf []byte
	for len(p) > 0 {
		r, size := utf8.DecodeRune(p)
		if r == utf8.RuneError && size == 1 {
			buf = append(buf, w.fallback)
		} else if e, ok := charmap.Windows1252.EncodeRune(r); ok {
			buf = append(buf, e)
		} else {
			buf = append(buf, w.fallback)
		}
		p = p[size:]
		n += size
	}

	_, err = w.w.Write(buf)
	if err != nil {
		return n, err
	}

	return n, nil
}
