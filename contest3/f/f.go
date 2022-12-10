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

type Vector struct {
	x, y int
}

func vectorProduct(v1 *Vector, v2 *Vector) int {
	return v1.x*v2.y - v1.y*v2.x
}

type Polygon struct {
	vertices []*Vector
}

func ScanVector(scanner *bufio.Scanner) *Vector {
	return &Vector{x: ScanInt(scanner), y: ScanInt(scanner)}
}

func ScanPoints(scanner *bufio.Scanner, vertices int) []*Vector {
	var points []*Vector
	for i := 0; i < vertices; i++ {
		points = append(points, ScanVector(scanner))
	}
	return points
}

func subVectors(first, second *Vector) *Vector {
	return &Vector{
		x: first.x - second.x,
		y: first.y - second.y,
	}
}

func sortVertices(points []*Vector) {
	sort.Slice(points, func(i, j int) bool {
		return points[i].x < points[j].x ||
			(points[i].x == points[j].x &&
				points[i].y < points[j].y)
	})
}

func convexHull(points []*Vector) ([]*Vector, []*Vector) {
	sortVertices(points)
	var upperPart, lowerPart []*Vector
	upperPart = append(upperPart, points[0])
	upperPart = append(upperPart, points[1])
	lowerPart = append(lowerPart, points[0])
	lowerPart = append(lowerPart, points[1])

	processVertex := func(part []*Vector, isUpper, index int) []*Vector {
		crossProduct := vectorProduct(
			subVectors((part)[len(part)-2], part[len(part)-1]),
			subVectors(points[index], part[len(part)-1]),
		)
		for len(part) > 1 && crossProduct*isUpper <= 0 {
			part = part[:len(part)-1]
		}
		return part
	}

	for i := 2; i < len(points); i++ {
		upperPart = processVertex(upperPart, 1, i)
		upperPart = append(upperPart, points[i])
		lowerPart = processVertex(lowerPart, -1, i)
		lowerPart = append(lowerPart, points[i])
	}

	lowerPart = lowerPart[1 : len(lowerPart)-1]
	return upperPart, lowerPart
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)

	N := ScanInt(scanner)
	points := ScanPoints(scanner, N)

	upperPart, lowerPart := convexHull(points)
	fmt.Println(len(upperPart) + len(lowerPart))
	for _, vertex := range upperPart {
		fmt.Printf("%d %d\n", int(vertex.x), int(vertex.y))
	}
	for i := len(lowerPart) - 1; i >= 0; i-- {
		fmt.Printf("%d %d\n", int(lowerPart[i].x), int(lowerPart[i].y))
	}
}
