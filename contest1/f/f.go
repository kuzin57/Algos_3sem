package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	alphabetSize = 26
)

func getNumber(r byte) int {
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

func run(line string) (int, error) {
	line += "{"

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
			cnt[classes[i]]++
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

	positions = positions[:len(positions)-1]

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

	var length int
	for i := range lcp {
		temp := len(line) - positions[i] - lcp[i]
		length += (lcp[i]*temp + ((temp+1)*temp)/2)
	}
	return length, nil
}

func main() {
	var (
		ans int
		err error
	)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Scan()

	line := scanner.Text()
	if ans, err = run(line); err != nil {
		panic(err)
	}
	fmt.Println(ans)
}
