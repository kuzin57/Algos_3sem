package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
)

const (
	epsilon = 1e-7
)

var (
	errLinesParallel  = fmt.Errorf("lines are parallel")
	errPointsAreEqual = fmt.Errorf("points are equal")
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
	a, b, c float64
}

type vector struct {
	x, y float64
}

func equal(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func ScanVector(scanner *bufio.Scanner) *vector {
	x, y := ScanInt(scanner), ScanInt(scanner)
	return &vector{x: float64(x), y: float64(y)}
}

func (v *vector) areCollinear(other *vector) bool {
	return equal(v.x*other.y, other.x*v.y)
}

func subVectors(first, second *vector) *vector {
	return &vector{
		x: first.x - second.x,
		y: first.y - second.y,
	}
}

func (l *line) areParallel(other *line) bool {
	return equal(l.a*other.b, other.a*l.b)
}

func getIntersection(first, second *line) (*vector, error) {
	if first.areParallel(second) {
		return nil, errLinesParallel
	}

	toDivide := first.a*second.b - second.a*first.b
	return &vector{
		x: (second.c*first.b - first.c*second.b) / toDivide,
		y: (second.a*first.c - first.a*second.c) / toDivide,
	}, nil
}

func buildLine(firstPoint, secondPoint *vector) (*line, error) {
	if firstPoint.x == secondPoint.x && secondPoint.y == firstPoint.y {
		return nil, errPointsAreEqual
	}

	return &line{
		a: firstPoint.y - secondPoint.y,
		b: secondPoint.x - firstPoint.x,
		c: firstPoint.x*secondPoint.y -
			secondPoint.x*firstPoint.y,
	}, nil
}

func linesEqual(firstLine, secondLine *line) bool {
	return equal(firstLine.a*secondLine.b, firstLine.b*secondLine.a) &&
		equal(firstLine.a*secondLine.c, firstLine.c*secondLine.a)
}

func betweenTwoNumbers(a, b, c float64) bool {
	return min(a, b) <= c+epsilon && c <= max(a, b)+epsilon
}

func betweenTwoPoints(firstPoint, secondPoint, thirdPoint *vector) bool {
	if !subVectors(secondPoint, firstPoint).areCollinear(
		subVectors(secondPoint, thirdPoint),
	) {
		return false
	}

	return betweenTwoNumbers(firstPoint.x, thirdPoint.x, secondPoint.x) &&
		betweenTwoNumbers(firstPoint.y, thirdPoint.y, secondPoint.y)
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
	firstLine, err := buildLine(a, b)
	if errors.Is(errPointsAreEqual, err) {
		if betweenTwoPoints(c, a, d) {
			return true
		}
		return false
	}

	secondLine, err := buildLine(c, d)
	if errors.Is(errPointsAreEqual, err) {
		if betweenTwoPoints(a, c, b) {
			return true
		}
		return false
	}

	intersection, err := getIntersection(firstLine, secondLine)
	if errors.Is(errLinesParallel, err) {
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
