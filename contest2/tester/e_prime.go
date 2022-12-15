package main

// import (
// 	"bufio"
// 	"bytes"
// 	"errors"
// 	"fmt"
// 	"os"
// 	"strconv"
// )

// const (
// 	MODULE = 7340033
// 	ROOT   = 5
// 	INVERT = 4404020
// 	POWER  = 1 << 20
// )

// var (
// 	errZeroFirstCoeff = errors.New("The ears of a dead donkey")
// )

// func ScanInt(scanner *bufio.Scanner) int {
// 	if !scanner.Scan() {
// 		panic("nothing to scan")
// 	}
// 	str := scanner.Text()
// 	n, err := strconv.Atoi(str)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return n
// }

// const whitespaceSymbols = "\t\n\v\f\r "

// func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
// 	if atEOF && len(data) == 0 {
// 		return 0, nil, nil
// 	}

// 	i := bytes.IndexAny(data, whitespaceSymbols)
// 	if i > 0 {
// 		return i + 1, data[:i], nil
// 	}
// 	if i == -1 {
// 		if atEOF {
// 			return len(data), data, nil
// 		}
// 		return 0, nil, nil
// 	}

// 	trimmed := bytes.TrimLeft(data, whitespaceSymbols)
// 	advance = len(data) - len(trimmed)
// 	if atEOF && len(trimmed) != 0 {
// 		trimmed = trimmed[:bytes.IndexAny(trimmed, whitespaceSymbols)]
// 		return advance + len(trimmed), trimmed, nil
// 	}
// 	return advance, nil, nil
// }

// func getReversed(num int, log int) int {
// 	left := 0
// 	right := log - 1
// 	for left < right {
// 		bitLeft := (num & (1 << left)) >> left
// 		bitRight := (num & (1 << right)) >> right
// 		num ^= (bitLeft << left)
// 		num ^= (bitRight << left)
// 		num ^= (bitRight << right)
// 		num ^= (bitLeft << right)
// 		left++
// 		right--
// 	}
// 	return num
// }

// func linearRepr(a, b int) (int, int) {
// 	if b == 1 {
// 		return 0, 1
// 	}
// 	var (
// 		prevX = 0
// 		prevY = 1
// 		x     = 1
// 		y     = -(a / b)
// 	)
// 	y += (MODULE * (-y/MODULE + 1))
// 	a, b = b, a%b
// 	for b != 0 {
// 		q := a / b
// 		if a%b != 0 {
// 			tmpX := x
// 			tmpY := y
// 			x = (prevX - x*q)
// 			if x < 0 {
// 				x += (MODULE * (-x/MODULE + 1))
// 			} else {
// 				x %= MODULE
// 			}
// 			y = (prevY - y*q)
// 			if y < 0 {
// 				y += (MODULE * (-y/MODULE + 1))
// 			} else {
// 				y %= MODULE
// 			}
// 			prevX = tmpX
// 			prevY = tmpY
// 		}
// 		a, b = b, a%b
// 	}
// 	return x, y
// }

// func swap(arr []int, first, second int) {
// 	tmp := arr[first]
// 	arr[first] = arr[second]
// 	arr[second] = tmp
// }

// func deleteZeros(poly []int) {
// 	for len(poly) > 1 && poly[len(poly)-1] == 0 {
// 		poly = poly[:len(poly)-1]
// 	}
// }

// func fft(poly []int, log int, isInvertFFT bool) {
// 	if len(poly) == 1 {
// 		return
// 	}

// 	for i := 0; i < len(poly); i++ {
// 		reversed := getReversed(i, log)
// 		if i < reversed {
// 			swap(poly, i, reversed)
// 		}
// 	}

// 	curOffset := 1
// 	for i := 0; i < log; i++ {
// 		var root int
// 		if isInvertFFT {
// 			root = INVERT
// 		} else {
// 			root = ROOT
// 		}
// 		for j := curOffset * 2; j < POWER; j <<= 1 {
// 			root = (root * root) % MODULE
// 		}
// 		for j := 0; j < len(poly); j += curOffset * 2 {
// 			curRoot := 1
// 			for k := 0; k < curOffset; k++ {
// 				tmp1 := poly[k+j]
// 				tmp2 := poly[j+k+curOffset]
// 				poly[k+j] = tmp1 + (curRoot*tmp2)%MODULE
// 				if poly[k+j] > MODULE {
// 					poly[k+j] -= MODULE
// 				}
// 				poly[k+j+curOffset] = tmp1 - (curRoot*tmp2)%MODULE
// 				if poly[k+j+curOffset] < 0 {
// 					poly[k+j+curOffset] += MODULE
// 				}
// 				curRoot = (curRoot * root) % MODULE
// 			}
// 		}
// 		curOffset *= 2
// 	}
// 	if isInvertFFT {
// 		_, inverted := linearRepr(MODULE, len(poly))
// 		for i := range poly {
// 			poly[i] = (poly[i] * inverted) % MODULE
// 		}
// 	}
// }

// func multiply(first_poly_arg []int, second_poly_arg []int) (deg int, poly []int) {
// 	if len(first_poly_arg) == 0 || len(second_poly_arg) == 0 {
// 		poly = append(poly, 0)
// 		return 0, poly
// 	}
// 	for len(first_poly_arg) > 1 && first_poly_arg[len(first_poly_arg)-1] == 0 {
// 		first_poly_arg = first_poly_arg[:(len(first_poly_arg) - 1)]
// 	}
// 	for len(second_poly_arg) > 1 && second_poly_arg[len(second_poly_arg)-1] == 0 {
// 		second_poly_arg = second_poly_arg[:len(second_poly_arg)-1]
// 	}
// 	if first_poly_arg[len(first_poly_arg)-1] == 0 || second_poly_arg[len(second_poly_arg)-1] == 0 {
// 		poly = append(poly, 0)
// 		return 0, poly
// 	}
// 	var (
// 		first_poly  []int
// 		second_poly []int
// 	)
// 	first_poly = append(first_poly, first_poly_arg...)
// 	second_poly = append(second_poly, second_poly_arg...)
// 	for len(first_poly) < len(second_poly) {
// 		first_poly = append(first_poly, 0)
// 	}

// 	for len(second_poly) < len(first_poly) {
// 		second_poly = append(second_poly, 0)
// 	}

// 	min_deg_two := 1
// 	log := 0
// 	for min_deg_two < len(first_poly) {
// 		log++
// 		min_deg_two <<= 1
// 	}
// 	min_deg_two <<= 1
// 	log++

// 	for len(first_poly) < min_deg_two {
// 		first_poly = append(first_poly, 0)
// 		second_poly = append(second_poly, 0)
// 	}

// 	fft(first_poly, log, false)
// 	fft(second_poly, log, false)

// 	for i := range first_poly {
// 		first_poly[i] = (first_poly[i] * second_poly[i]) % MODULE
// 	}

// 	fft(first_poly, log, true)

// 	for first_poly[len(first_poly)-1] == 0 {
// 		first_poly = first_poly[:(len(first_poly) - 1)]
// 	}
// 	reverse(first_poly)
// 	return len(first_poly) - 1, first_poly
// }

// func reverse(arr []int) {
// 	left := 0
// 	right := len(arr) - 1
// 	for left < right {
// 		arr[left], arr[right] = arr[right], arr[left]
// 		left++
// 		right--
// 	}
// }

// func main() {
// 	scanner := bufio.NewScanner(os.Stdin)
// 	scanner.Buffer(nil, 1<<30)
// 	scanner.Split(splitFunc)

// 	n := ScanInt(scanner)
// 	first_poly := make([]int, n+1)
// 	for i := 0; i <= n; i++ {
// 		first_poly[n-i] = ScanInt(scanner)
// 	}

// 	m := ScanInt(scanner)
// 	second_poly := make([]int, m+1)
// 	for i := 0; i <= m; i++ {
// 		second_poly[m-i] = ScanInt(scanner)
// 	}

// 	degRes, resultPoly := multiply(first_poly, second_poly)
// 	fmt.Printf("%d ", degRes)
// 	for _, coeff := range resultPoly {
// 		fmt.Printf("%d ", coeff)
// 	}
// 	fmt.Println()
// }
