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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)
	var (
		N      = ScanInt(scanner)
		primes []int
		mind   = make([]int, N+1)
		sum    int
	)

	for i := 2; i <= N; i++ {
		mind[i] = i
	}
	for i := 2; i <= N; i++ {
		if mind[i] == i {
			primes = append(primes, i)
		}
		for _, prime := range primes {
			if prime*i > N || prime > mind[i] {
				break
			}
			mind[i*prime] = prime
		}
	}
	for i := 2; i <= N; i++ {
		if mind[i] != i {
			sum += mind[i]
		}
	}
	fmt.Println(sum)
}
