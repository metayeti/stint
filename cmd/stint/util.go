package main

import (
	"strings"
	"unicode"
)

func stripEmojis(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for _, r := range s {
		if unicode.IsSymbol(r) {
			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}