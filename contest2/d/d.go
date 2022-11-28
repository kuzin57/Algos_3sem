package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
)

const (
	limit = 10_000_000
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

func getPrimes() []int {
	var (
		mind   = make([]int, limit+1)
		primes []int
	)
	for i := 2; i <= limit; i++ {
		mind[i] = i
	}
	for i := 2; i <= limit; i++ {
		if mind[i] == i {
			primes = append(primes, i)
		}
		for _, prime := range primes {
			if prime*i > limit || prime > mind[i] {
				break
			}
			mind[i*prime] = prime
		}
	}
	return primes
}

func findMinNotUsedPrimeIndex(index int, primes []int, primesUsed []bool) int {
	for i := index; i < len(primes); i++ {
		if !primesUsed[i] {
			return i
		}
	}
	return -1
}

func getNextDivisibleByPrime(num int, primes []int, primesUsed []bool, iteratePrimes func([]int, []bool, int, bool) bool) int {
	var withoutUsedPrimes bool
	for {
		withoutUsedPrimes = !(iteratePrimes(primes, primesUsed, num, true))
		if !withoutUsedPrimes {
			num++
			continue
		}
		return num
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)
	var (
		n              = ScanInt(scanner)
		arr            = make([]int, n)
		primes         = getPrimes()
		primesUsed     = make([]bool, len(primes))
		takeOnlyPrimes bool
		index          int
	)

	iterateThroughPrimes := func(primes []int, primesUsed []bool, number int, firstCase bool) bool {
		for j, prime := range primes {
			if number < prime {
				break
			}
			if number%prime == 0 && primesUsed[j] && firstCase {
				return true
			}
			if number%prime == 0 && !firstCase {
				primesUsed[j] = true
			}
		}
		return false
	}

	for i := 0; i < n; i++ {
		arr[i] = ScanInt(scanner)
	}
	for i := 0; i < n; i++ {
		if takeOnlyPrimes {
			index = findMinNotUsedPrimeIndex(index, primes, primesUsed)
			primesUsed[index] = true
			arr[i] = primes[index]
			continue
		}

		takeOnlyPrimes = iterateThroughPrimes(primes, primesUsed, arr[i], true)
		if takeOnlyPrimes {
			arr[i] = getNextDivisibleByPrime(arr[i], primes, primesUsed, iterateThroughPrimes)
		}

		iterateThroughPrimes(primes, primesUsed, arr[i], false)
	}
	for _, num := range arr {
		fmt.Printf("%d ", num)
	}
	fmt.Println()
}
