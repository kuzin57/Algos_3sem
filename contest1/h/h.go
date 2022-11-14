package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
)

const (
	alphabetSize = 26
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

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func getNumber(r byte) int {
	if '@' <= r && r <= 'Z' {
		return int(r - '@')
	}
	return int(r - 'a' + alphabetSize + 1)
}

func firstInit(line string, cnt []int, positions []int, classes []int) error {
	for _, letter := range line {
		num := getNumber(byte(letter))
		cnt[num]++
	}

	for i := 1; i < len(cnt); i++ {
		cnt[i] += cnt[i-1]
	}
	for i := len(line) - 1; i >= 0; i-- {
		num := getNumber(line[i])
		cnt[num]--
		positions[cnt[num]] = i
	}

	classes[positions[0]] = 0
	for i := 1; i < len(line); i++ {
		classes[positions[i]] = classes[positions[i-1]]
		if line[positions[i]] != line[positions[i-1]] {
			classes[positions[i]]++
		}
	}
	return nil
}

func fill(arr []int, a int) {
	for i := range arr {
		arr[i] = a
	}
}

func run(line string, number int) (string, error) {
	line += "@"

	var (
		cnt          = make([]int, 2*alphabetSize+1)
		positions    = make([]int, len(line))
		classes      = make([]int, len(line))
		newPositions = make([]int, len(line))
		newClasses   = make([]int, len(line))
		pos          = make([]int, len(line))
		lcp          = make([]int, len(line)-1)
	)

	if err := firstInit(line, cnt, positions, classes); err != nil {
		return "", err
	}

	curDegree := 1
	cnt = make([]int, len(classes))
	for curDegree < len(line) {
		for i := range line {
			newPositions[i] = positions[i] - curDegree
			if newPositions[i] < 0 {
				newPositions[i] += len(line)
			}
		}

		fill(cnt, 0)
		for i := range line {
			cnt[classes[newPositions[i]]]++
		}

		for i := 1; i < len(line); i++ {
			cnt[i] += cnt[i-1]
		}

		for i := len(line) - 1; i >= 0; i-- {
			cnt[classes[newPositions[i]]]--
			positions[cnt[classes[newPositions[i]]]] = newPositions[i]
		}

		fill(newClasses, 0)
		for i := 1; i < len(line); i++ {
			newClasses[positions[i]] = newClasses[positions[i-1]]
			if classes[positions[i]] != classes[positions[i-1]] ||
				classes[(positions[i]+curDegree)%len(line)] != classes[(positions[i-1]+curDegree)%len(line)] {
				newClasses[positions[i]]++
			}
		}

		for i := range line {
			classes[i] = newClasses[i]
		}
		curDegree *= 2
	}

	positions = positions[1:]

	for i := range positions {
		pos[positions[i]] = i
	}

	var curLCP int
	line = line[:len(line)-1]
	for i := range line {
		curLCP = max(curLCP-1, 0)
		if pos[i] == len(line)-1 {
			continue
		}
		j := positions[pos[i]+1]
		for i+curLCP < len(line) && j+curLCP < len(line) && line[i+curLCP] == line[j+curLCP] {
			curLCP++
		}
		lcp[pos[i]] = curLCP
	}

	return findKthSubstr(positions, lcp, line, number)
}

func findKthSubstr(suffixArr []int, lcp []int, line string, seqNumber int) (string, error) {
	cnt := len(line) - suffixArr[0]
	if cnt >= seqNumber {
		return line[suffixArr[0]:(suffixArr[0] + seqNumber)], nil
	}
	for i := 1; i < len(suffixArr); i++ {
		curElem := suffixArr[i]
		lcpElem := lcp[i-1]
		if cnt+len(line)-curElem-lcpElem >= seqNumber {
			return line[curElem:(seqNumber - cnt + lcpElem + curElem)], nil
		}
		cnt += len(line) - curElem - lcpElem
	}
	return line[suffixArr[len(suffixArr)-1]:], nil
}

func main() {
	var (
		ans       string
		err       error
		seqNumber int
	)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)
	scanner.Scan()
	line := scanner.Text()
	seqNumber = ScanInt(scanner)

	if ans, err = run(line, seqNumber); err != nil {
		panic(err)
	}

	fmt.Println(ans)
}
