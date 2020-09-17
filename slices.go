package util

import (
	"log"
	"strings"
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
		log.Fatalf("too deep recursion -%v- lvl %v", strings.Join(els, ","), depth)
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
				// log.Printf("%v %v - %-20v contains %-20v => red: %v", i, j, els[i], els[j], redundant)
			}
		}
	}

	if len(redundant) == 0 {
		return els
	}

	if len(redundant) > 0 {
		// We cannot remove the redun
		// log.Printf("redund %v - els -%v-", redundant, strings.Join(els, ","))
		for i := 0; i < len(redundant); i++ {
			els[redundant[i]] = ""
		}
		// log.Printf("redund %v - els -%v-", redundant, strings.Join(els, ","))
		for i := 0; i < len(els); i++ {
			// log.Printf("iterating redund %v", i)
			if els[i] == "" {
				head := els[:i]
				tail := els[i+1:]
				// log.Printf("trimming  %v of %v:  -%v- -%v- ", i, len(els)-1, strings.Join(head, ","), strings.Join(tail, ","))
				els = append(head, tail...)
				// log.Printf("result   -%v- ", strings.Join(els, ","))
			}
		}
	}
	els = TrimRedundant(els, depth+1)
	return els

}
