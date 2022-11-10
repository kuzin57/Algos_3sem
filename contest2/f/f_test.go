package main

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func check(poly []int, invert []int, m int) bool {
	res := multiply(poly, invert)
	if res[0] != 1 {
		return false
	}
	if len(res) == 1 {
		return true
	}
	for len(res) < m {
		res = append(res, 0)
	}
	for i := 1; i < m; i++ {
		if res[i] != 0 {
			fmt.Println("res", len(res), len(poly), len(invert))
			return false
		}
	}
	return true
}

func genTest() (int, int, []int) {
	m := rand.Int()%1000 + 1
	n := rand.Int() % 1000
	poly := make([]int, n+1)
	for i := 0; i < n+1; i++ {
		poly[i] = rand.Int() % MODULE
	}
	return m, n, poly
}

func TestRun(t *testing.T) {
	for i := 0; i < 1000; i++ {
		m, _, poly := genTest()
		deleteZeros(poly)
		res, err := findInvertPoly(poly, m)
		if poly[0] == 0 {
			assert.Error(t, err)
		} else {
			s := check(poly, res, m)
			assert.Equal(t, true, s)
			if !s {
				// fmt.Println("m, n, poly res", res)
				return
			}
		}
	}
}
