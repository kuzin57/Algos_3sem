package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
)

const module = 1000000007

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
	a, b = b, a%b
	for b != 0 {
		q := a / b
		if a%b != 0 {
			tmpX := x
			tmpY := y
			x = (prevX - x*q) % module
			y = (prevY - y*q) % module
			prevX = tmpX
			prevY = tmpY
		}
		a, b = b, a%b
	}
	return x, y
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)

	var (
		a         = ScanInt(scanner)
		b         = ScanInt(scanner)
		c         = ScanInt(scanner)
		d         = ScanInt(scanner)
		inverse_b int
		inverse_d int
	)

	if a < 0 {
		a += module
	}
	if b < 0 {
		b += module
	}
	if c < 0 {
		c += module
	}
	if d < 0 {
		d += module
	}
	_, inverse_b = linearRepr(module, b)
	if inverse_b < 0 {
		inverse_b += module
	}
	_, inverse_d = linearRepr(module, d)
	if inverse_d < 0 {
		inverse_d += module
	}
	fmt.Println(((a*inverse_b)%module + (c*inverse_d)%module) % module)
}
