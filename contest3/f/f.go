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

func (v *Vector) String() string {
	return fmt.Sprintf("%d %d", int(v.x), int(v.y))
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

func ScanPolygon(scanner *bufio.Scanner, vertices int) *Polygon {
	newPolynom := &Polygon{}
	for i := 0; i < vertices; i++ {
		newPolynom.vertices = append(newPolynom.vertices, ScanVector(scanner))
	}
	return newPolynom
}

func subVectors(first, second *Vector) *Vector {
	return &Vector{
		x: first.x - second.x,
		y: first.y - second.y,
	}
}

func (p *Polygon) sortVertices() {
	sort.Slice(p.vertices, func(i, j int) bool {
		return p.vertices[i].x < p.vertices[j].x ||
			(p.vertices[i].x == p.vertices[j].x &&
				p.vertices[i].y < p.vertices[j].y)
	})
}

func (p *Polygon) convexHull() ([]*Vector, []*Vector) {
	p.sortVertices()
	var upperPart, lowerPart []*Vector
	upperPart = append(upperPart, p.vertices[0])
	upperPart = append(upperPart, p.vertices[1])
	lowerPart = append(lowerPart, p.vertices[0])
	lowerPart = append(lowerPart, p.vertices[1])

	processVertex := func(part *[]*Vector, isUpper, index int) {
		for len(*part) > 1 &&
			vectorProduct(
				subVectors((*part)[len(*part)-2], (*part)[len(*part)-1]),
				subVectors(p.vertices[index], (*part)[len(*part)-1]),
			)*isUpper <= 0 {
			*part = (*part)[:len(*part)-1]
		}
	}

	for i := 2; i < len(p.vertices); i++ {
		processVertex(&upperPart, 1, i)
		upperPart = append(upperPart, p.vertices[i])
		processVertex(&lowerPart, -1, i)
		lowerPart = append(lowerPart, p.vertices[i])
	}

	lowerPart = lowerPart[1 : len(lowerPart)-1]
	return upperPart, lowerPart
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)

	N := ScanInt(scanner)
	polygon := ScanPolygon(scanner, N)

	upperPart, lowerPart := polygon.convexHull()
	fmt.Println(len(upperPart) + len(lowerPart))
	for _, vertex := range upperPart {
		fmt.Println(vertex)
	}
	for i := len(lowerPart) - 1; i >= 0; i-- {
		fmt.Println(lowerPart[i])
	}
}
