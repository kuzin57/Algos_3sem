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
	root   = 5
	invert = 4404020
	power  = 1 << 20
)

var (
	errZeroFirstCoeff = errors.New("The ears of a dead donkey")
	fabric            modularFabric
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

type modular struct {
	value  int
	module int
}

type modularFabric struct {
	module int
}

func initFabric(module int) modularFabric {
	return modularFabric{module: module}
}

func (f modularFabric) buildModular(number int) modular {
	ans := modular{value: number % f.module}
	ans.normalize()
	return ans
}

func (m *modular) normalize() {
	if m.value >= 0 {
		m.value %= m.module
		return
	}
	m.value = (m.value + m.module*(-m.value/m.module+1)) % m.module
}

func sumModulars(first modular, second modular) modular {
	ans := modular{value: first.value + second.value}
	ans.normalize()
	return ans
}

func subModulars(first modular, second modular) modular {
	ans := modular{value: first.value - second.value}
	ans.normalize()
	return ans
}

func multModulars(first modular, second modular) modular {
	ans := modular{value: first.value * second.value}
	ans.normalize()
	return ans
}

func findInvert(m modular) modular {
	_, invert := euclidCoeffs(m.module, m.value)
	ans := modular{value: invert}
	ans.normalize()
	return ans
}

func (m modular) String() string {
	return fmt.Sprintf("%d", m.value)
}

func ScanModular(scanner *bufio.Scanner, f modularFabric) modular {
	return f.buildModular(ScanInt(scanner))
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

func euclidCoeffs(a, b int) (int, int) {
	var (
		prevX = fabric.buildModular(0)
		prevY = fabric.buildModular(1)
		x     = fabric.buildModular(1)
		y     = fabric.buildModular(-(a / b))
	)
	a, b = b, a%b
	for b != 0 {
		q := fabric.buildModular(a / b)
		if a%b != 0 {
			tmpX := x
			tmpY := y
			x = subModulars(prevX, multModulars(x, q))
			y = subModulars(prevY, multModulars(y, q))
			prevX = tmpX
			prevY = tmpY
		}
		a, b = b, a%b
	}
	return x.value, y.value
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
		var curRoot int
		if isInvertFFT {
			curRoot = invert
		} else {
			curRoot = root
		}
		for j := curOffset * 2; j < power; j <<= 1 {
			curRoot = (curRoot * curRoot) % module
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
				curRoot = (curRoot * curRoot) % module
			}
		}
		curOffset *= 2
	}

	if isInvertFFT {
		_, inverted := euclidCoeffs(module, len(poly))
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
	for len(firstPoly) > 0 && firstPoly[len(firstPoly)-1] == 0 {
		firstPoly = firstPoly[:(len(firstPoly) - 1)]
	}
	return firstPoly
}

func findInvertPoly(poly []int, deg int, fabric *modularFabric) ([]int, error) {
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

	_, invert[0] = euclidCoeffs(module, poly[0])
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

	fabric = initFabric(module)
	result, err := findInvertPoly(poly, m, &fabric)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(result) > m {
		result = result[:m]
	}
	for len(result) < m {
		result = append(result, 0)
	}
	for _, coeff := range result {
		fmt.Printf("%d ", coeff%module)
	}
	fmt.Println()
}
