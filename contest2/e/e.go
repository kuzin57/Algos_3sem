package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
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

func fft(poly []complex128, isInvert bool) {
	if len(poly) == 1 {
		return
	}

	var even_degrees, odd_degrees []complex128
	for i, coeff := range poly {
		switch i % 2 {
		case 0:
			even_degrees = append(even_degrees, coeff)
		case 1:
			odd_degrees = append(odd_degrees, coeff)
		}
	}

	fft(even_degrees, isInvert)
	fft(odd_degrees, isInvert)

	angle := (2 * math.Pi) / float64(len(poly))
	if isInvert {
		angle *= (-1)
	}

	first_root := complex(math.Cos(angle), math.Sin(angle))
	cur_root := complex(1, 0)
	for i := 0; i < len(poly)/2; i++ {
		poly[i] = even_degrees[i] + cur_root*odd_degrees[i]
		poly[i+len(poly)/2] = even_degrees[i] - cur_root*odd_degrees[i]
		cur_root *= first_root
		if isInvert {
			poly[i] /= 2
			poly[i+len(poly)/2] /= 2
		}
	}
}

func multiply(firstPoly []complex128, secondPoly []complex128) []int { // equal lengths are expected
	minDegTwo := 1
	for minDegTwo < len(firstPoly) {
		minDegTwo <<= 1
	}
	minDegTwo <<= 1

	for len(firstPoly) < minDegTwo {
		firstPoly = append(firstPoly, 0)
		secondPoly = append(secondPoly, 0)
	}

	fft(firstPoly, false)
	fft(secondPoly, false)

	helpPoly := make([]complex128, len(firstPoly))
	for i := range firstPoly {
		helpPoly[i] = firstPoly[i] * secondPoly[i]
	}

	fft(helpPoly, true)
	resPoly := make([]int, len(firstPoly))
	for i := range resPoly {
		resPoly[i] = int(math.Round(real(helpPoly[i])))
	}

	for len(resPoly) > 0 && resPoly[len(resPoly)-1] == 0 {
		resPoly = resPoly[:(len(resPoly) - 1)]
	}
	reverse(resPoly)
	return resPoly
}

func reverse(arr []int) {
	left := 0
	right := len(arr) - 1
	for left < right {
		arr[left], arr[right] = arr[right], arr[left]
		left++
		right--
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)

	n := ScanInt(scanner)
	firstPoly := make([]complex128, n+1)
	for i := 0; i <= n; i++ {
		val := ScanInt(scanner)
		firstPoly[n-i] = complex(float64(val), 0)
	}

	m := ScanInt(scanner)
	secondPoly := make([]complex128, m+1)
	for i := 0; i <= m; i++ {
		val := ScanInt(scanner)
		secondPoly[m-i] = complex(float64(val), 0)
	}

	for len(firstPoly) < len(secondPoly) {
		firstPoly = append(firstPoly, 0)
	}

	for len(secondPoly) < len(firstPoly) {
		secondPoly = append(secondPoly, 0)
	}
	resultPoly := multiply(firstPoly, secondPoly)
	fmt.Printf("%d ", len(resultPoly)-1)
	for _, coeff := range resultPoly {
		fmt.Printf("%d ", coeff)
	}
	fmt.Println()
}
