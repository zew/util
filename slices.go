package util

import (
	"strings"

	"github.com/zew/logx"
)

func Contains(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
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
