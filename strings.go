package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/zew/logx"
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

// See test cases
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

// Any slice element, that is contained in another slice element
// gets trimmed.
func TrimRedundant(els []string, optArg ...int) []string {

	depth := 0
	if len(optArg) > 0 {
		depth = optArg[0]
	}
	if depth > 20 {
		logx.Fatalf("too deep recursion -%v- lvl %v", strings.Join(els, ","), depth)
		return els
	}

	redundant := []int{}
	for i := 0; i < len(els); i++ {
		for j := 0; j < len(els); j++ {
			if i == j {
				continue
			}
			if strings.Contains(els[i], els[j]) {
				redundant = append(redundant, j)
				// logx.Printf("%v %v - %-20v contains %-20v => red: %v", i, j, els[i], els[j], redundant)
			}
		}
	}

	if len(redundant) == 0 {
		return els
	}

	if len(redundant) > 0 {
		// We cannot remove the redun
		// logx.Printf("redund %v - els -%v-", redundant, strings.Join(els, ","))
		for i := 0; i < len(redundant); i++ {
			els[redundant[i]] = ""
		}
		// logx.Printf("redund %v - els -%v-", redundant, strings.Join(els, ","))
		for i := 0; i < len(els); i++ {
			// logx.Printf("iterating redund %v", i)
			if els[i] == "" {
				head := els[:i]
				tail := els[i+1:]
				// logx.Printf("trimming  %v of %v:  -%v- -%v- ", i, len(els)-1, strings.Join(head, ","), strings.Join(tail, ","))
				els = append(head, tail...)
				// logx.Printf("result   -%v- ", strings.Join(els, ","))
			}
		}
	}
	els = TrimRedundant(els, depth+1)
	return els

}

// normalize spaces
var replNewLines = strings.NewReplacer("\r\n", " ", "\r", " ", "\n", " ")

var replTabs = strings.NewReplacer("\t", " ")
var doubleSpaces = regexp.MustCompile("([ ]+)")

// All kinds of newlines, tabs and double spaces
// are reduced to single space.
// It paves the way for later beautification.
func NormalizeInnerWhitespace(s string) string {
	s = replNewLines.Replace(s)
	s = replTabs.Replace(s)
	s = doubleSpaces.ReplaceAllString(s, " ")
	return s
}

// We want only one newline; not several.
// We also want to catch newlines
// trailed by any kind of white space.
var undoubleNewlines = regexp.MustCompile(`(\r?\n(\s*))+`)

func UndoubleNewlines(s string) string {
	return undoubleNewlines.ReplaceAllString(s, "\n")
}

var removeProtocol = strings.NewReplacer(
	"http://", "",
	"https://", "",
)

func RemoveProtocol(s string) string {
	return removeProtocol.Replace(s)
}

// Func isSpacey is a shortcut func.
// It detects if there is ONLY whitspace,
// but nothing else.
// All combinations of whitespace lead to return true.
func IsSpacey(s string) bool {
	s = strings.TrimSpace(s) // TrimSpace removes leading-trailing \n \r\n
	if s == "" {
		return true
	}
	return false
}
