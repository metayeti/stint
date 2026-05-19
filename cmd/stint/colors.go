package main

import (
	"fmt"
	"math"
	"strings"
)

const (
	reset   = "\033[0m"
	orange  = "\033[38;2;255;165;0m"
)

func hueShiftedString(input string) string {
	var result strings.Builder
	length := len(input)

	for i, char := range input {
		hue := 30 + (float64(i)/float64(length))*25 // orange to yellow

		r, g, b := hsvToRGB(hue, 0.85, 0.98)

		result.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c", r, g, b, char))
	}

	result.WriteString(reset)
	return result.String()
}

func hsvToRGB(h, s, v float64) (r, g, b int) {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	var r1, g1, b1 float64

	switch {
	case h < 60:
		r1, g1, b1 = c, x, 0
	case h < 120:
		r1, g1, b1 = x, c, 0
	case h < 180:
		r1, g1, b1 = 0, c, x
	case h < 240:
		r1, g1, b1 = 0, x, c
	case h < 300:
		r1, g1, b1 = x, 0, c
	default:
		r1, g1, b1 = c, 0, x
	}

	r = int((r1 + m) * 255)
	g = int((g1 + m) * 255)
	b = int((b1 + m) * 255)
	return
}