package main

import "fmt"

func BuildPrefFunc(s string) []int {
	ans := make([]int, len(s))

	for i := 1; i < len(s); i++ {
		res := ans[i-1]

		for res > 0 && s[res] != s[i] {
			res = ans[res-1]
		}

		if s[i] == s[res] {
			res++
		}
		ans[i] = res
	}

	return ans
}

func main() {
	var (
		s string
		t string
	)

	fmt.Scan(&s)
	fmt.Scan(&t)

	prefFunc := BuildPrefFunc(t + "#" + s)
	for i, pref := range prefFunc {
		if len(t) == pref {
			fmt.Printf("%d ", i-len(t)-len(t))
		}
	}
}
