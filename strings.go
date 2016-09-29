package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

func UpTo(s string, numChars int) string {
	if len(s) > numChars {
		return s[:numChars]
	}
	return s
}

func UpToR(s string, numChars int) string {
	if len(s) > numChars {
		return s[len(s)-numChars:]
	}
	return s
}

func Ellipsoider(s string, maxChars int) string {
	if len(s) > maxChars {
		maxChars -= 6
		return s[:maxChars/2] + " ... " + s[len(s)-maxChars/2:]
	}
	return s
}

func IndentedDump(v interface{}) string {

	firstColLeftMostPrefix := " "
	byts, err := json.MarshalIndent(v, firstColLeftMostPrefix, "\t")
	if err != nil {
		s := fmt.Sprintf("error indent: %v\n", err)
		return s
	}

	byts = bytes.Replace(byts, []byte(`\u003c`), []byte("<"), -1)
	byts = bytes.Replace(byts, []byte(`\u003e`), []byte(">"), -1)
	byts = bytes.Replace(byts, []byte(`\n`), []byte("\n"), -1)

	return string(byts)
}

func EnsureUtf8(haystack string) string {
	ret := bytes.Buffer{}
	for _, codepoint := range haystack {
		ret.WriteRune(codepoint)
	}
	return ret.String()
}

func HumanizeFloat(f float64) string {
	str := fmt.Sprintf("%v", f)
	strs := strings.Split(str, ".")
	if len(strs) == 1 {
		return str
	}
	if pos := strings.Index(strs[1], "0000"); pos > -1 {
		if pos == 0 {
			return strs[0]
		}
		// pos > 0
		strs[1] = strs[1][0:pos]
		return strings.Join(strs, ".")
	}

	if pos := strings.Index(strs[1], "9999"); pos > -1 {
		// 1234.99999;  -1234.99999
		if pos == 0 {
			if f >= 0 {
				return fmt.Sprintf("%v", math.Ceil(f))
			}
			return fmt.Sprintf("%v", math.Floor(f))
		}
		// 1234.599999;  -1234.599999
		strs[1] = strs[1][0:pos]
		rn, sz := utf8.DecodeLastRuneInString(strs[1])
		if rn != utf8.RuneError && sz == 1 {
			// log.Printf("%v - %c => %c\n", strs[1], rn, rn+1)
			rn++
			strs[1] = strs[1][0 : len(strs[1])-1]
			strs[1] = fmt.Sprintf("%v%c", strs[1], rn)
		}

		return strings.Join(strs, ".")
	}

	return strings.Join(strs, ".")
}
