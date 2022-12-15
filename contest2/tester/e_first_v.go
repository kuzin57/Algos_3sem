package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
)

const MODULE = 7340033

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

func get_reversed(num int, size int) int {
	left := 0
	right := size - 1
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

func swap(arr []complex128, first, second int) {
	tmp := arr[first]
	arr[first] = arr[second]
	arr[second] = tmp
}

func fft(poly []complex128, size int, isInvertFFT bool) {
	if len(poly) == 1 {
		return
	}

	for i := 0; i < len(poly); i++ {
		reversed := get_reversed(i, size)
		if i < reversed {
			swap(poly, i, reversed)
		}
	}

	curOffset := 1
	for i := 0; i < size; i++ {
		fmt.Println("hello")
		angle := (2 * math.Pi) / float64(curOffset*2)
		if isInvertFFT {
			angle *= -1
		}
		root := complex(math.Cos(angle), math.Sin(angle))
		for j := 0; j < len(poly); j += curOffset * 2 {
			curRoot := complex(1, 0)
			for k := 0; k < curOffset; k++ {
				tmp1 := poly[k+j]
				tmp2 := poly[j+k+curOffset]
				poly[k+j] = tmp1 + curRoot*tmp2
				poly[k+j+curOffset] = tmp1 - curRoot*tmp2
				curRoot *= root
				if isInvertFFT {
					poly[k+j] /= 2
					poly[k+j+curOffset] /= 2
				}
			}
		}
		curOffset *= 2
	}
	fmt.Println("##############################")
}

func run(first_poly []complex128, second_poly []complex128) (deg int, poly []int) {
	fmt.Println("first", first_poly, "second", second_poly)
	for len(first_poly) < len(second_poly) {
		first_poly = append(first_poly, 0)
	}

	for len(second_poly) < len(first_poly) {
		second_poly = append(second_poly, 0)
	}

	min_deg_two := 1
	size := 0
	for min_deg_two < len(first_poly) {
		size++
		min_deg_two <<= 1
	}
	min_deg_two <<= 1
	size++

	for len(first_poly) < min_deg_two {
		first_poly = append(first_poly, 0)
		second_poly = append(second_poly, 0)
	}

	fft(first_poly, size, false)
	fft(second_poly, size, false)

	res_poly := make([]int, len(first_poly))
	for i := range first_poly {
		first_poly[i] *= second_poly[i]
	}

	fft(first_poly, size, true)
	for i := range res_poly {
		res_poly[i] = int(math.Round(real(first_poly[i])))
	}

	fmt.Println("res", res_poly)

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

	degRes, resultPoly := run(first_poly, second_poly)
	fmt.Printf("%d ", degRes)
	for _, coeff := range resultPoly {
		fmt.Printf("%d ", coeff)
	}
	fmt.Println()
}
