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

func getNumber(r int, linesNumber int) int {
	return r + linesNumber
}

func firstInit(line []int, cnt []int, positions []int, classes []int, linesNumber int) {
	for _, letter := range line {
		num := getNumber(letter, linesNumber)
		cnt[num]++
	}

	for i := 1; i < len(cnt); i++ {
		cnt[i] += cnt[i-1]
	}
	for i := len(line) - 1; i >= 0; i-- {
		num := getNumber(line[i], linesNumber)
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
}

func fill(arr []int, a int) {
	for i := range arr {
		arr[i] = a
	}
}

func run(line []int, lengths []int, linesNumber int) []int {
	var (
		cnt          = make([]int, 1+alphabetSize+linesNumber)
		positions    = make([]int, len(line))
		classes      = make([]int, len(line))
		newPositions = make([]int, len(line))
		newClasses   = make([]int, len(line))
		pos          = make([]int, len(line))
		lcp          = make([]int, len(line)-1)
	)

	firstInit(line, cnt, positions, classes, linesNumber)

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
	return solve(positions, lcp, line, lengths, linesNumber)
}

func joinSets(first map[int]struct{}, second map[int]struct{}) {
	for key := range second {
		first[key] = struct{}{}
	}
}

func findCommon(suffArr []int, lcp []int, line []int, ans []int, lengths []int, linesNumber int) {
	var (
		stackVal     = make([]int, 0)
		stackSets    = make([]map[int]struct{}, 0)
		from         = make([]int, 0)
		cnt          int
		indexLengths int
		curLen       int
	)

	curLen = lengths[indexLengths]
	for i := 0; i < len(line); i++ {
		from = append(from, cnt)
		if i+1 == curLen && indexLengths < len(lengths) {
			indexLengths++
			curLen += lengths[indexLengths]
			cnt++
			continue
		}
		if i+1 == curLen {
			indexLengths++
			cnt++
		}
	}
	for i := 0; i < len(lcp); i++ {
		switch {
		case i == 0 && lcp[0] > 0:
			stackVal = append(stackVal, lcp[0])
			stackSets = append(stackSets, make(map[int]struct{}))
			stackSets[0][from[suffArr[0]]] = struct{}{}
			stackSets[0][from[suffArr[1]]] = struct{}{}
		case lcp[i] > lcp[i-1]:
			stackVal = append(stackVal, lcp[i])
			stackSets = append(stackSets, make(map[int]struct{}))
			stackSets[len(stackSets)-1][from[suffArr[i]]] = struct{}{}
			stackSets[len(stackSets)-1][from[suffArr[i+1]]] = struct{}{}
		case lcp[i] < lcp[i-1]:
			curSet := make(map[int]struct{})
			stackValLen := len(stackVal)
			stackSetsLen := len(stackSets)
			for stackValLen > 0 && stackVal[stackValLen-1] > lcp[i] {
				joinSets(curSet, stackSets[stackSetsLen-1])
				if len(stackSets) >= 2 {
					joinSets(stackSets[stackSetsLen-2], stackSets[stackSetsLen-1])
				}
				ans[len(stackSets[stackSetsLen-1])] = max(
					ans[len(stackSets[stackSetsLen-1])],
					stackVal[stackValLen-1],
				)
				stackVal = stackVal[:(stackValLen - 1)]
				stackSets = stackSets[:(stackSetsLen - 1)]
				stackValLen--
				stackSetsLen--
			}
			if lcp[i] != 0 {
				stackVal = append(stackVal, lcp[i])
				stackSets = append(stackSets, make(map[int]struct{}))
				stackSets[len(stackSets)-1][from[suffArr[i]]] = struct{}{}
				stackSets[len(stackSets)-1][from[suffArr[i+1]]] = struct{}{}
				joinSets(stackSets[len(stackSets)-1], curSet)
			}
		default:
			if len(stackSets) > 0 {
				stackSets[len(stackSets)-1][from[suffArr[i+1]]] = struct{}{}
			}
		}
	}
}

func solve(suffArr []int, lcp []int, line []int, lengths []int, linesNumber int) []int {
	var (
		ans = make([]int, linesNumber+1)
	)

	findCommon(suffArr, lcp, line, ans, lengths, linesNumber)
	return ans
}

func completeAnswer(ans []int) {
	for i := len(ans) - 1; i > 0; i-- {
		if ans[i] > ans[i-1] {
			ans[i-1] = ans[i]
		}
	}
}

func main() {
	var (
		ans         []int
		bigline     []int
		lengths     []int
		linesNumber int
	)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(splitFunc)
	scanner.Buffer(nil, 1<<30)
	linesNumber = ScanInt(scanner)
	for i := 0; i < linesNumber; i++ {
		scanner.Scan()
		line := scanner.Text()
		for _, b := range line {
			bigline = append(bigline, int(b-'a'))
		}
		lengths = append(lengths, len(line)+1)
		bigline = append(bigline, -i-1)
	}

	run(bigline, lengths, linesNumber)
	completeAnswer(ans)
	for i := 2; i < len(ans); i++ {
		fmt.Println(ans[i])
	}
}
