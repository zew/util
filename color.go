package util

import (
	"fmt"
	"strings"
)

func ColorTint(channel string, baseHue int, pct float64) string {

	channel = strings.ToLower(channel)

	if channel != "r" && channel != "g" && channel != "b" {
		return ""
	}
	if pct > 1 || pct < -1 {
		return ""
	}
	if baseHue < 0 || baseHue > 255 {
		return ""
	}

	changeH := baseHue + int(float64(255-baseHue)*pct) // changeHue
	if changeH > 255 {
		changeH = 255
	}
	if changeH < 0 {
		changeH = 0
	}
	if channel == "r" {
		return fmt.Sprintf("RGB(%02d,%02d,%02d)", changeH, baseHue, baseHue)
	}
	if channel == "g" {
		return fmt.Sprintf("RGB(%02d,%02d,%02d)", baseHue, changeH, baseHue)
	}
	if channel == "b" {
		return fmt.Sprintf("RGB(%02d,%02d,%02d)", baseHue, baseHue, changeH)
	}
	return ""

}
