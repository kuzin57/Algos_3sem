package main

import "fmt"

func CreatePrefFunc(s string) []int {
	ans := make([]int, len(s))
	ans[0] = 0

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

	prefFunc := CreatePrefFunc(t + "#" + s)
	for i, pref := range prefFunc {
		if len(t) == pref {
			fmt.Print(i-len(t)-len(t), " ")
		}
	}
}
