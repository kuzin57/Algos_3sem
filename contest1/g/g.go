package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	alphabetSize  = 26
	minimalSymbol = 'A' - 1
)

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

func run(line string) ([]int, error) {
	line += "@"

	var (
		cnt          = make([]int, 2*alphabetSize+1)
		positions    = make([]int, len(line))
		classes      = make([]int, len(line))
		newPositions = make([]int, len(line))
		newClasses   = make([]int, len(line))
	)

	if err := firstInit(line, cnt, positions, classes); err != nil {
		return nil, err
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

	return positions, nil
}

func main() {
	var (
		ans []int
		err error
	)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Scan()

	line := scanner.Text()
	if ans, err = run(line); err != nil {
		panic(err)
	}

	for i := range ans {
		fmt.Print(ans[i]+1, " ")
	}
	fmt.Println()
}
