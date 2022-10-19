package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
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

func CreateZFunc(s []int) []int {
	var (
		left  int
		right int
	)
	ans := make([]int, len(s))

	for i := 1; i < len(s); i++ {
		if i <= right {
			ans[i] = ans[i-left]
		}

		if ans[i]+i > right && right >= i {
			ans[i] = right - i
		}

		for i+ans[i] < len(s) && s[i+ans[i]] == s[ans[i]] {
			ans[i]++
		}

		if i+ans[i] > right {
			left = i
			right = i + ans[i]
		}
	}

	return ans
}

func Reverse(arr []int) {
	var (
		left  = 0
		right = len(arr) - 1
	)

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

	var (
		numberCubes int
		// numberColors int
	)

	numberCubes = ScanInt(scanner)
	_ = ScanInt(scanner)

	colors := make([]int, numberCubes)

	for i := 0; i < numberCubes; i++ {
		colors[i] = ScanInt(scanner)
	}

	line := make([]int, 0)
	line = append(line, colors...)
	line = append(line, 0)

	Reverse(colors)
	line = append(line, colors...)

	zFunc := CreateZFunc(line)
	ans := make([]int, 0)
	for i := 0; i < len(colors); i++ {
		pos := len(line) - (i+1)*2
		if zFunc[pos] >= i+1 {
			ans = append(ans, numberCubes-i-1)
		}
	}

	sort.Ints(ans)
	ans = append(ans, numberCubes)

	for _, num := range ans {
		fmt.Print(num, " ")
	}
}
