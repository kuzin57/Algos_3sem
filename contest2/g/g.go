package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	MODULE = 7340033
	MOD7   = 7
	ROOT   = 5
	INVERT = 4404020
	POWER  = 1 << 20
)

var (
	errZeroFirstCoeff = errors.New("The ears of a dead donkey")
)

func ScanInt(scanner *bufio.Scanner) int {
	if !scanner.Scan() {
		panic("nothing to scan")
	}
	str := scanner.Text()
	n, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return n
}

const whitespaceSymbols = "\t\n\v\f\r "

func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	i := bytes.IndexAny(data, whitespaceSymbols)
	if i > 0 {
		return i + 1, data[:i], nil
	}
	if i == -1 {
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	}

	trimmed := bytes.TrimLeft(data, whitespaceSymbols)
	advance = len(data) - len(trimmed)
	if atEOF && len(trimmed) != 0 {
		trimmed = trimmed[:bytes.IndexAny(trimmed, whitespaceSymbols)]
		return advance + len(trimmed), trimmed, nil
	}
	return advance, nil, nil
}

func getReversed(num int, log int) int {
	left := 0
	right := log - 1
	for left < right {
		bitLeft := (num & (1 << left)) >> left
		bitRight := (num & (1 << right)) >> right
		num ^= (bitLeft << left)
		num ^= (bitRight << left)
		num ^= (bitRight << right)
		num ^= (bitLeft << right)
		left++
		right--
	}
	return num
}

func gcd(a, b int) (int, int, int) {
	var x, y int
	if a == 0 {
		return 0, 1, b
	}
	newX, newY, d := gcd(b%a, a)
	x = (newY - ((b / a) * newX))
	y = newX
	return x, y, d
}

func getInvert(a int, mod int) int {
	x, _, _ := gcd(a, mod)
	for x < 0 {
		x += mod
	}
	return x
}

func swap(arr []int, first, second int) {
	tmp := arr[first]
	arr[first] = arr[second]
	arr[second] = tmp
}

func deleteZeros(poly *[]int) {
	for len(*poly) > 1 && (*poly)[len(*poly)-1] == 0 {
		*poly = (*poly)[:(len(*poly) - 1)]
	}
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func fft(poly []int, log int, isInvertFFT bool) {
	if len(poly) == 1 {
		return
	}

	for i := 0; i < len(poly); i++ {
		reversed := getReversed(i, log)
		if i < reversed {
			swap(poly, i, reversed)
		}
	}

	curOffset := 1
	for i := 0; i < log; i++ {
		var root int
		if isInvertFFT {
			root = INVERT
		} else {
			root = ROOT
		}
		for j := curOffset * 2; j < POWER; j <<= 1 {
			root = (root * root) % MODULE
		}
		for j := 0; j < len(poly); j += curOffset * 2 {
			curRoot := 1
			for k := 0; k < curOffset; k++ {
				tmp1 := poly[k+j]
				tmp2 := poly[j+k+curOffset]
				poly[k+j] = tmp1 + (curRoot*tmp2)%MODULE
				if poly[k+j] > MODULE {
					poly[k+j] -= MODULE
				}
				poly[k+j+curOffset] = tmp1 - (curRoot*tmp2)%MODULE
				if poly[k+j+curOffset] < 0 {
					poly[k+j+curOffset] += MODULE
				}
				curRoot = (curRoot * root) % MODULE
			}
		}
		curOffset *= 2
	}
	if isInvertFFT {
		inverted := getInvert(len(poly), MODULE)
		for i := range poly {
			poly[i] = (poly[i] * inverted) % MODULE
		}
	}
}

func multiply(first_poly []int, second_poly_arg []int) []int {
	var (
		second_poly []int
	)
	second_poly = append(second_poly, second_poly_arg...)
	if len(first_poly) == 0 || len(second_poly) == 0 {
		first_poly = []int{0}
		return first_poly
	}
	for len(first_poly) > 1 && first_poly[len(first_poly)-1] == 0 {
		first_poly = first_poly[:(len(first_poly) - 1)]
	}
	for len(second_poly) > 1 && second_poly[len(second_poly)-1] == 0 {
		second_poly = second_poly[:len(second_poly)-1]
	}
	if first_poly[len(first_poly)-1] == 0 || second_poly[len(second_poly)-1] == 0 {
		first_poly = []int{0}
		return first_poly
	}
	for len(first_poly) < len(second_poly) {
		first_poly = append(first_poly, 0)
	}

	for len(second_poly) < len(first_poly) {
		second_poly = append(second_poly, 0)
	}

	min_deg_two := 1
	log := 0
	for min_deg_two < len(first_poly) {
		log++
		min_deg_two <<= 1
	}
	min_deg_two <<= 1
	log++

	for len(first_poly) < min_deg_two {
		first_poly = append(first_poly, 0)
		second_poly = append(second_poly, 0)
	}

	fft(first_poly, log, false)
	fft(second_poly, log, false)

	for i := range first_poly {
		first_poly[i] = (first_poly[i] * second_poly[i]) % MODULE
	}

	fft(first_poly, log, true)
	deleteZeros(&first_poly)
	for i := range first_poly {
		first_poly[i] %= MOD7
	}
	return first_poly
}

func findInvertPoly(poly []int, deg int) ([]int, error) {
	if poly[0] == 0 {
		return nil, errZeroFirstCoeff
	}
	for len(poly) < deg {
		extra := make([]int, deg-len(poly))
		poly = append(poly, extra...)
	}
	var (
		invert = make([]int, len(poly))
		curDeg = 1
	)

	invert[0] = getInvert(poly[0], MOD7)
	for curDeg < deg {
		firstPart := make([]int, 0)
		for i := 0; i < curDeg; i++ {
			firstPart = append(firstPart, poly[i]%MOD7)
		}
		secondPart := make([]int, 0)
		for i := curDeg; i < len(poly); i++ {
			secondPart = append(secondPart, poly[i]%MOD7)
		}

		firstPart = multiply(firstPart, invert)
		if len(firstPart) > deg {
			firstPart = firstPart[:deg+1]
		}
		if curDeg >= len(firstPart) {
			firstPart = make([]int, 0)
		} else {
			firstPart = firstPart[curDeg:]
		}
		secondPart = multiply(secondPart, invert)
		if len(secondPart) > deg {
			secondPart = secondPart[:deg+1]
		}
		for i := 0; i < max(len(firstPart), len(secondPart)); i++ {
			if i == len(firstPart) {
				firstPart = append(firstPart, 0)
			}
			if i < len(firstPart) && i < len(secondPart) {
				firstPart[i] = (-firstPart[i] - secondPart[i] + 2*MOD7) % MOD7
			} else if i < len(secondPart) {
				firstPart[i] = (-secondPart[i] + MOD7) % MOD7
			} else {
				firstPart[i] = (-firstPart[i] + MOD7) % MOD7
			}
		}
		firstPart = multiply(firstPart, invert)
		if len(firstPart) > deg {
			firstPart = firstPart[:deg+1]
		}
		for i := range firstPart {
			if i+curDeg == len(invert) {
				invert = append(invert, 0)
			}
			invert[i+curDeg] = (invert[i+curDeg] + firstPart[i]) % MOD7
		}
		curDeg <<= 1
	}
	deleteZeros(&invert)
	return invert, nil
}

func revPoly(poly []int) {
	left := 0
	right := len(poly) - 1
	for left < right {
		poly[left], poly[right] = poly[right], poly[left]
		left++
		right--
	}
}

func divPoly(firstPoly []int, secondPoly []int) ([]int, []int) {
	deg := len(firstPoly) - len(secondPoly) + 1
	invertedSecondPoly, _ := findInvertPoly(secondPoly, deg)
	invertedSecondPoly = invertedSecondPoly[:deg]
	revQ := multiply(invertedSecondPoly, firstPoly)
	revQ = revQ[:deg]
	revPoly(revQ)
	var R []int
	revPoly(firstPoly)
	fmt.Println("secondPoly:", secondPoly)
	revPoly(secondPoly)
	fmt.Println("secondPoly:", secondPoly)
	secondPoly = multiply(secondPoly, revQ)
	for i := range firstPoly {
		last := firstPoly[i] % MOD7
		if i < len(secondPoly) {
			last = (last + MOD7 - secondPoly[i]) % MOD7
		}
		R = append(R, last%MOD7)
	}
	return revQ, R
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)

	m := ScanInt(scanner)
	firstPoly := make([]int, m+1)
	for i := 0; i <= m; i++ {
		firstPoly[i] = ScanInt(scanner)
	}
	n := ScanInt(scanner)
	secondPoly := make([]int, n+1)
	for i := 0; i <= n; i++ {
		secondPoly[i] = ScanInt(scanner)
	}
	if len(firstPoly) < len(secondPoly) {
		fmt.Println("0 0")
		fmt.Printf("%d ", len(firstPoly)-1)
		for _, coeff := range firstPoly {
			fmt.Printf("%d ", coeff)
		}
		fmt.Println()
		return
	}
	Q, R := divPoly(firstPoly, secondPoly)
	zerosFinished := false
	for i := len(Q) - 1; i >= 0; i-- {
		if Q[i] != 0 && !zerosFinished {
			zerosFinished = true
			fmt.Printf("%d ", i)
		}
		if !zerosFinished {
			continue
		}
		fmt.Printf("%d ", Q[i])
	}
	if !zerosFinished {
		fmt.Print("0 0")
	}
	fmt.Println()
	zerosFinished = false
	for i := len(R) - 1; i >= 0; i-- {
		if R[i] != 0 && !zerosFinished {
			zerosFinished = true
			fmt.Printf("%d ", i)
		}
		if !zerosFinished {
			continue
		}
		fmt.Printf("%d ", R[i])
	}
	if !zerosFinished {
		fmt.Print("0 0")
	}
	fmt.Println()
}
