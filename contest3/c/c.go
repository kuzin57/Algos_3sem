package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
)

const (
	epsilon = float64(1) / 10000000
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

type line struct {
	A, B, C float64
}

type vector struct {
	x, y float64
}

func Equal(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func ScanVector(scanner *bufio.Scanner) *vector {
	x, y := ScanInt(scanner), ScanInt(scanner)
	return &vector{x: float64(x), y: float64(y)}
}

func isCollinear(first, second *vector) bool {
	return Equal(first.x*second.y, second.x*first.y)
}

func subVectors(first, second *vector) *vector {
	return &vector{
		x: first.x - second.x,
		y: first.y - second.y,
	}
}

func isParallel(first, second *line) bool {
	return Equal(first.A*second.B, second.A*first.B)
}

func getIntersection(first, second *line) (*vector, error) {
	if isParallel(first, second) {
		return nil, fmt.Errorf("line are parallel")
	}

	toDivide := first.A*second.B - second.A*first.B
	return &vector{
		x: (second.C*first.B - first.C*second.B) / toDivide,
		y: (second.A*first.C - first.A*second.C) / toDivide,
	}, nil
}

func buildLine(firstPoint, secondPoint *vector) (*line, error) {
	if firstPoint.x == secondPoint.x && secondPoint.y == firstPoint.y {
		return nil, fmt.Errorf("points are equal")
	}

	return &line{
		A: firstPoint.y - secondPoint.y,
		B: secondPoint.x - firstPoint.x,
		C: firstPoint.x*secondPoint.y -
			secondPoint.x*firstPoint.y,
	}, nil
}

func linesEqual(firstLine, secondLine *line) bool {
	return Equal(firstLine.A*secondLine.B, firstLine.B*secondLine.A) &&
		Equal(firstLine.A*secondLine.C, firstLine.C*secondLine.A)
}

func betweenTwoNumbers(a, b, c float64) bool {
	return min(a, b) <= c+epsilon && c <= max(a, b)+epsilon
}

func betweenTwoPoints(firstPoint, secondPoint, thirdPoint *vector) bool {
	if !isCollinear(
		subVectors(secondPoint, firstPoint),
		subVectors(secondPoint, thirdPoint),
	) {
		return false
	}

	return betweenTwoNumbers(firstPoint.x, thirdPoint.x, secondPoint.x) &&
		betweenTwoNumbers(firstPoint.y, thirdPoint.y, secondPoint.y)
}

func buildLineWithCheck(
	firstPoint,
	secondPoint,
	firstEndOtherSegment,
	secondEndOtherSegment *vector,
) (*line, bool, error) {
	firstLine, err := buildLine(firstPoint, secondPoint)
	if err != nil {
		if betweenTwoPoints(
			firstEndOtherSegment, firstPoint, secondEndOtherSegment) {
			return nil, true, err
		}
		return nil, false, err
	}
	return firstLine, false, nil
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a < b {
		return b
	}
	return a
}

func checkLocationOnLine(a, b, c, d *vector, firstLine, secondLine *line) bool {
	if !linesEqual(firstLine, secondLine) {
		return false
	}
	return betweenTwoPoints(c, a, d) || betweenTwoPoints(c, b, d) || betweenTwoPoints(a, c, b) ||
		betweenTwoPoints(a, d, b)
}

func run(a, b, c, d *vector) bool {
	firstLine, ans, err := buildLineWithCheck(a, b, c, d)
	if err != nil {
		return ans
	}

	secondLine, ans, err := buildLineWithCheck(c, d, a, b)
	if err != nil {
		return ans
	}

	intersection, err := getIntersection(firstLine, secondLine)
	if err != nil {
		return checkLocationOnLine(a, b, c, d, firstLine, secondLine)
	}

	return betweenTwoPoints(a, intersection, b) && betweenTwoPoints(c, intersection, d)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)

	ans := run(
		ScanVector(scanner),
		ScanVector(scanner),
		ScanVector(scanner),
		ScanVector(scanner),
	)

	if ans {
		fmt.Println("YES")
		return
	}
	fmt.Println("NO")
}
