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

func run(first_poly []complex128, second_poly []complex128) (deg int, poly []int) {
	for len(first_poly) < len(second_poly) {
		first_poly = append(first_poly, 0)
	}

	for len(second_poly) < len(first_poly) {
		second_poly = append(second_poly, 0)
	}

	min_deg_two := 1
	for min_deg_two < len(first_poly) {
		min_deg_two <<= 1
	}
	min_deg_two <<= 1

	for len(first_poly) < min_deg_two {
		first_poly = append(first_poly, 0)
		second_poly = append(second_poly, 0)
	}

	fft(first_poly, false)
	fft(second_poly, false)

	res_poly := make([]int, len(first_poly))
	for i := range first_poly {
		first_poly[i] *= second_poly[i]
	}

	fft(first_poly, true)
	for i := range res_poly {
		res_poly[i] = int(math.Round(real(first_poly[i])))
	}

	for res_poly[len(res_poly)-1] == 0 {
		res_poly = res_poly[:(len(res_poly) - 1)]
	}
	reverse(res_poly)
	return len(res_poly) - 1, res_poly
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

func stupid_algo(first_poly []complex128, second_poly []complex128) []int {
	res := make([]complex128, len(first_poly)+len(second_poly)-1)
	for i := range first_poly {
		for j := range second_poly {
			res[i+j] += first_poly[i] * second_poly[j]
		}
	}
	result := make([]int, len(res))
	for i := range res {
		result[i] = int(real(res[i]))
	}
	reverse(result)
	return result
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)

	n := ScanInt(scanner)
	first_poly := make([]complex128, n+1)
	for i := 0; i <= n; i++ {
		val := ScanInt(scanner)
		first_poly[n-i] = complex(float64(val), 0)
	}

	m := ScanInt(scanner)
	second_poly := make([]complex128, m+1)
	for i := 0; i <= m; i++ {
		val := ScanInt(scanner)
		second_poly[m-i] = complex(float64(val), 0)
	}

	// fmt.Println("right:", stupid_algo(first_poly, second_poly))
	degRes, resultPoly := run(first_poly, second_poly)
	fmt.Printf("%d ", degRes)
	for _, coeff := range resultPoly {
		fmt.Printf("%d ", coeff)
	}
	fmt.Println()
}
