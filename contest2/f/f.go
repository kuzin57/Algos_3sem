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
	module = 7340033
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

func linearRepr(a, b int) (int, int) {
	if b == 1 {
		return 0, 1
	}
	var (
		prevX = 0
		prevY = 1
		x     = 1
		y     = -(a / b)
	)
	y += (module * (-y/module + 1))
	a, b = b, a%b
	for b != 0 {
		q := a / b
		if a%b != 0 {
			tmpX := x
			tmpY := y
			x = (prevX - x*q)
			if x < 0 {
				x += (module * (-x/module + 1))
			} else {
				x %= module
			}
			y = (prevY - y*q)
			if y < 0 {
				y += (module * (-y/module + 1))
			} else {
				y %= module
			}
			prevX = tmpX
			prevY = tmpY
		}
		a, b = b, a%b
	}
	return x, y
}

func deleteZeros(poly []int) {
	for len(poly) > 0 && poly[len(poly)-1] == 0 {
		poly = poly[:(len(poly) - 1)]
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
			poly[i], poly[reversed] = poly[reversed], poly[i]
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
			root = (root * root) % module
		}
		for j := 0; j < len(poly); j += curOffset * 2 {
			curRoot := 1
			for k := 0; k < curOffset; k++ {
				tmp1 := poly[k+j]
				tmp2 := poly[j+k+curOffset]
				poly[k+j] = tmp1 + (curRoot*tmp2)%module
				if poly[k+j] > module {
					poly[k+j] -= module
				}
				poly[k+j+curOffset] = tmp1 - (curRoot*tmp2)%module
				if poly[k+j+curOffset] < 0 {
					poly[k+j+curOffset] += module
				}
				curRoot = (curRoot * root) % module
			}
		}
		curOffset *= 2
	}

	if isInvertFFT {
		_, inverted := linearRepr(module, len(poly))
		for i := range poly {
			poly[i] = (poly[i] * inverted) % module
		}
	}
}

func multiply(firstPoly []int, secondPolyArg []int) []int {
	var (
		secondPoly []int
	)
	secondPoly = append(secondPoly, secondPolyArg...)

	if len(firstPoly) == 0 || len(secondPoly) == 0 {
		firstPoly = []int{0}
		return firstPoly
	}
	for len(firstPoly) > 1 && firstPoly[len(firstPoly)-1] == 0 {
		firstPoly = firstPoly[:(len(firstPoly) - 1)]
	}
	for len(secondPoly) > 1 && secondPoly[len(secondPoly)-1] == 0 {
		secondPoly = secondPoly[:len(secondPoly)-1]
	}
	if firstPoly[len(firstPoly)-1] == 0 || secondPoly[len(secondPoly)-1] == 0 {
		firstPoly = []int{0}
		return firstPoly
	}
	for len(firstPoly) < len(secondPoly) {
		firstPoly = append(firstPoly, 0)
	}

	for len(secondPoly) < len(firstPoly) {
		secondPoly = append(secondPoly, 0)
	}

	minDegTwo := 1
	log := 0
	for minDegTwo < len(firstPoly) {
		log++
		minDegTwo <<= 1
	}
	minDegTwo <<= 1
	log++

	for len(firstPoly) < minDegTwo {
		firstPoly = append(firstPoly, 0)
		secondPoly = append(secondPoly, 0)
	}

	fft(firstPoly, log, false)
	fft(secondPoly, log, false)

	for i := range firstPoly {
		firstPoly[i] = (firstPoly[i] * secondPoly[i]) % module
	}

	fft(firstPoly, log, true)
	deleteZeros(firstPoly)
	return firstPoly
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

	_, invert[0] = linearRepr(module, poly[0])
	for curDeg < deg {
		firstPart := make([]int, 0)
		for i := 0; i < curDeg; i++ {
			firstPart = append(firstPart, poly[i])
		}
		secondPart := make([]int, 0)
		for i := curDeg; i < len(poly); i++ {
			secondPart = append(secondPart, poly[i])
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
				firstPart[i] = (-firstPart[i] - secondPart[i] + 2*module) % module
			} else if i < len(secondPart) {
				firstPart[i] = (-secondPart[i] + module) % module
			} else {
				firstPart[i] = (-firstPart[i] + module) % module
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
			invert[i+curDeg] = (invert[i+curDeg] + firstPart[i]) % module
		}
		curDeg <<= 1
	}
	return invert, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)

	m := ScanInt(scanner)
	n := ScanInt(scanner)
	poly := make([]int, n+1)
	for i := 0; i <= n; i++ {
		poly[i] = ScanInt(scanner)
	}

	result, err := findInvertPoly(poly, m)
	if err != nil {
		fmt.Println(err)
		return
	}
	for len(result) > m {
		result = result[:(len(result) - 1)]
	}
	for len(result) < m {
		result = append(result, 0)
	}
	for _, coeff := range result {
		fmt.Printf("%d ", coeff%module)
	}
	fmt.Println()
}
