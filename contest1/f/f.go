package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	alphabetSize = 26
)

func getAlphabetNumber(r byte) int {
	if 'A' <= r && r <= 'Z' {
		return int(r - 'A')
	}
	return int(r - 'a' + alphabetSize)
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func stableCountSort(
	line string,
	curDegree int,
	cnt []int,
	positions []int,
	oldClasses []int,
	newClasses []int,
	newPositions []int,
) {
	for i := range line {
		cnt[oldClasses[i]]++
	}

	for i := 1; i < len(line); i++ {
		cnt[i] += cnt[i-1]
	}

	for i := len(line) - 1; i >= 0; i-- {
		cnt[oldClasses[newPositions[i]]]--
		positions[cnt[oldClasses[newPositions[i]]]] = newPositions[i]
	}

	for i := 1; i < len(line); i++ {
		newClasses[positions[i]] = newClasses[positions[i-1]]
		if oldClasses[positions[i]] != oldClasses[positions[i-1]] ||
			oldClasses[(positions[i]+curDegree)%len(line)] != oldClasses[(positions[i-1]+curDegree)%len(line)] {
			newClasses[positions[i]]++
		}
	}
}

func fill(arr []int, a int) {
	for i := range arr {
		arr[i] = a
	}
}

func fillLCP(line string, positions []int, lcp []int) {
	var pos = make([]int, len(line))
	positions = positions[:len(positions)-1]

	for i := range positions {
		pos[positions[i]] = i
	}

	var curLCP int
	line = line[:len(line)-1]
	for i := range pos {
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
}

func run(line string) int {
	line += "{"

	var (
		cnt          = make([]int, 2*alphabetSize+1)
		positions    = make([]int, len(line))
		classes      = make([]int, len(line))
		newPositions = make([]int, len(line))
		newClasses   = make([]int, len(line))
		lcp          = make([]int, len(line)-1)
	)

	for i := range line {
		classes[i] = getAlphabetNumber(line[i])
		newPositions[i] = i
	}

	stableCountSort(line, 0, cnt, positions, classes, newClasses, newPositions)
	copy(classes, newClasses)

	cnt = make([]int, len(classes))
	for curDegree := 1; curDegree < len(line); curDegree *= 2 {
		for i := range line {
			newPositions[i] = positions[i] - curDegree
			if newPositions[i] < 0 {
				newPositions[i] += len(line)
			}
		}

		fill(cnt, 0)
		fill(newClasses, 0)
		stableCountSort(line, curDegree, cnt, positions, classes, newClasses, newPositions)

		copy(classes, newClasses)
	}

	fillLCP(line, positions, lcp)

	var length int
	for i := range lcp {
		temp := len(line) - positions[i] - lcp[i]
		length += (lcp[i]*temp + ((temp+1)*temp)/2)
	}
	return length
}

func main() {
	var (
		ans int
	)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Scan()

	line := scanner.Text()
	ans = run(line)
	fmt.Println(ans)
}
