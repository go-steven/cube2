package util

import (
	"fmt"
	"sort"
	"strings"
)

func Trim(s string) string {
	return strings.TrimSpace(s)
}

func LowerTrim(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func UpperTrim(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

func InArrayStr(s string, list []string) bool {
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})

	index := sort.SearchStrings(list, s)
	return index >= 0 && index < len(list) && list[index] == s
}

func Uint64Join(vals []uint64) string {
	str_vals := []string{}
	for _, v := range vals {
		str_vals = append(str_vals, fmt.Sprintf("%v", v))
	}
	return strings.Join(str_vals, ",")
}
