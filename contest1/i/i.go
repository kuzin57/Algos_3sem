package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	alphabetSize  = 26
	minimalSymbol = 'A' - 1
)

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func getAlphabetNumber(r byte) int {
	if minimalSymbol <= r && r <= 'Z' {
		return int(r - minimalSymbol)
	}
	return int(r - 'a' + alphabetSize + 1)
}

func firstInit(line string, cnt []int, positions []int, classes []int) error {
	for _, letter := range line {
		num := getAlphabetNumber(byte(letter))
		cnt[num]++
	}

	for i := 1; i < len(cnt); i++ {
		cnt[i] += cnt[i-1]
	}
	for i := len(line) - 1; i >= 0; i-- {
		num := getAlphabetNumber(line[i])
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

func run(line string) (int, error) {
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
		return 0, err
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

	return getKey(positions, lcp, line)
}

func getKey(suffArr []int, lcp []int, line string) (int, error) {
	var (
		ans      int
		minBegin int
		maxBegin int
		prevLCP  = -1
	)

	minBegin = suffArr[0]
	maxBegin = suffArr[0]
	for i := 1; i < len(suffArr); i++ {
		minBegin = min(minBegin, suffArr[i])
		maxBegin = max(maxBegin, suffArr[i])
		if lcp[i-1] != prevLCP {
			minBegin = min(suffArr[i-1], suffArr[i])
			maxBegin = max(suffArr[i-1], suffArr[i])
			prevLCP = lcp[i-1]
		}
		ans = max(ans, abs(minBegin-maxBegin)+prevLCP*prevLCP+prevLCP)
	}
	return max(ans, len(line)), nil
}

func maxSupref(line string) int {
	var curPref string
	var ans int
	for i := range line {
		curPref = line[:min(len(line), i+1)]
		if strings.HasSuffix(line, curPref) && line != curPref {
			ans = max(ans, len(curPref))
		}
	}
	return ans
}

func stupidAlgo(line string) int {
	var ans int
	for i := 0; i < len(line); i++ {
		for j := i + 1; j <= len(line); j++ {
			supref := maxSupref(line[i:j])
			ans = max(ans, j-i+supref*supref)
		}
	}
	return ans
}

func main() {
	var (
		ans  int
		line string
		err  error
	)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Scan()
	line = scanner.Text()

	if ans, err = run(line); err != nil {
		panic(err)
	}

	fmt.Println(ans)
}
