package main

import (
	"bufio"
	"bytes"
	"fmt"
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)
	var (
		n       = ScanInt(scanner)
		array   = make([]int, n)
		simples bool
		ones    int
	)

	for i := 0; i < n; i++ {
		array[i] = ScanInt(scanner)
	}

	minDist := n
	for i := 0; i < n; i++ {
		cur := array[i]
		if cur == 1 {
			simples = true
			ones++
		}
		for j := i + 1; j < n; j++ {
			cur = gcd(cur, array[j])
			if cur == 1 {
				simples = true
				minDist = min(minDist, j-i)
				break
			}
		}
	}
	switch simples {
	case true:
		if ones > 0 {
			fmt.Println(n - ones)
		} else {
			fmt.Println(minDist + n - 1)
		}
	case false:
		fmt.Println(-1)
	}
}
