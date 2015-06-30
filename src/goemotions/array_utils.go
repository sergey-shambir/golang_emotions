package main

import "unicode/utf8"

const (
	MAX_RUNE_SIZE = 8
)

func array_append_rune(dst []byte, src rune) []byte {
	buf := make([]byte, MAX_RUNE_SIZE, MAX_RUNE_SIZE)
	written := utf8.EncodeRune(buf, src)
	return append(dst, buf[0:written]...)
}
