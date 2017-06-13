package util

import (
	"fmt"
	"strconv"
	"strings"
)

func ColorTint(col, channel string, pct float64) string {

	col = strings.TrimPrefix(col, "#")
	col = strings.ToLower(col)
	if len(col) != 3 && len(col) != 6 {
		return ""
	}

	rgb := []string{}
	if len(col) == 3 {
		rgb = strings.Split(col, "")
		rgb[0] = rgb[0] + "0"
		rgb[1] = rgb[1] + "0"
		rgb[2] = rgb[2] + "0"
	}
	if len(col) == 6 {
		rgb = []string{col[:2], col[2:4], col[4:6]}
	}

	channel = strings.ToLower(channel)
	if channel != "r" && channel != "g" && channel != "b" {
		return ""
	}
	if pct > 1 || pct < -1 {
		return ""
	}

	baseHue := int64(0)
	if channel == "r" {
		baseHue, _ = strconv.ParseInt(rgb[0], 16, 9) // nine bit - cause it's a signed int
	}
	if channel == "g" {
		baseHue, _ = strconv.ParseInt(rgb[1], 16, 9)
	}
	if channel == "b" {
		baseHue, _ = strconv.ParseInt(rgb[2], 16, 9)
	}

	if baseHue < 0 || baseHue > 255 {
		return ""
	}
	// if baseHue > 212 && pct > 0 {
	// 	baseHue = 212
	// }
	// if baseHue < 48 && pct < 0 {
	// 	baseHue = 48
	// }
	baseHue8 := uint8(baseHue)

	changeH := baseHue8 + uint8(float64(255-baseHue8)*pct) // changeHue
	if changeH > 255 {
		changeH = 255
	}
	if changeH < 0 {
		changeH = 0
	}

	ret := ""
	if channel == "r" {
		ret = fmt.Sprintf("#%02x%02s%02s", changeH, rgb[1], rgb[2])
	}
	if channel == "g" {
		ret = fmt.Sprintf("#%02s%02x%02s", rgb[0], changeH, rgb[2])
	}
	if channel == "b" {
		ret = fmt.Sprintf("#%02s%02s%02x", rgb[0], rgb[1], changeH)
	}

	// fmt.Printf("from %v - to %v - %v - %v - %v => %v\n", col, rgb, baseHue, baseHue8, changeH, ret)

	return ret

}
