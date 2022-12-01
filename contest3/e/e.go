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
	epsilon   = float64(1) / 10000000
	precision = float64(10000000)
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
	x, y float64
}

func NewVector(x, y float64) *Vector {
	return &Vector{x: x, y: y}
}

type Polygon struct {
	vertices []*Vector
	edges    []*Edge
}

func ScanVector(scanner *bufio.Scanner) *Vector {
	return &Vector{x: float64(ScanInt(scanner)), y: float64(ScanInt(scanner))}
}

func ScanPolygon(scanner *bufio.Scanner, vertices int) *Polygon {
	newPolynom := &Polygon{}
	for i := 0; i < vertices; i++ {
		newPolynom.vertices = append(newPolynom.vertices, ScanVector(scanner))
	}
	return newPolynom
}

type Edge struct {
	line      *Line
	firstEnd  *Vector
	secondEnd *Vector
}

func NewEdge(line *Line, firstEnd *Vector, secondEnd *Vector) *Edge {
	return &Edge{
		line:      line,
		firstEnd:  firstEnd,
		secondEnd: secondEnd,
	}
}

type Line struct {
	A, B, C float64
}

func (l *Line) checkHalfPlanes(firstPoint, secondPoint *Vector) float64 {
	if (l.A*firstPoint.x+l.B*firstPoint.y+l.C)*
		(l.A*secondPoint.x+l.B*secondPoint.y+l.C) < 0 {
		return 1
	}
	return 0
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

func buildLine(firstVector, secondVector *Vector) (*Line, error) {
	if Equal(firstVector.x, secondVector.x) && Equal(secondVector.y, firstVector.y) {
		return nil, fmt.Errorf("Vectors are equal")
	}

	return &Line{
		A: firstVector.y - secondVector.y,
		B: secondVector.x - firstVector.x,
		C: firstVector.x*secondVector.y -
			secondVector.x*firstVector.y,
	}, nil
}

func betweenTwoNumbers(a, b, c float64) bool {
	return min(a, b) <= c+epsilon && c <= max(a, b)+epsilon
}

func Equal(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func betweenTwoVectors(firstVector, secondVector, thirdVector *Vector) bool {
	return betweenTwoNumbers(firstVector.x, thirdVector.x, secondVector.x) &&
		betweenTwoNumbers(firstVector.y, thirdVector.y, secondVector.y)
}

func (p *Polygon) initEdges() {
	verticesNumber := len(p.vertices)
	for i := range p.vertices {
		newLine, err := buildLine(p.vertices[i], p.vertices[(i+1)%verticesNumber])
		if err != nil {
			panic(err)
		}

		if len(p.edges) > 0 && p.edges[len(p.edges)-1].line.A == 0 && newLine.A == 0 {
			p.edges[len(p.edges)-1].secondEnd = p.vertices[(i+1)%verticesNumber]
			continue
		}
		newEdge := NewEdge(newLine, p.vertices[i], p.vertices[(i+1)%verticesNumber])
		p.edges = append(p.edges, newEdge)
	}
}

func (p *Polygon) isInside(point *Vector) int {
	processTwoEdges := func(line *Line, firstEdge *Edge, secondEdge *Edge) int {
		return int(line.checkHalfPlanes(
			firstEdge.firstEnd, secondEdge.secondEnd,
		))
	}
	var cntEdges int
	edgesNumber := len(p.edges)
	horizontalLine := &Line{A: 0, B: 1, C: -point.y}
	for i, edge := range p.edges {
		if Equal(edge.line.A, 0) && point.y != edge.firstEnd.y {
			continue
		}
		if Equal(edge.line.A, 0) && point.x <= edge.secondEnd.x {
			cntEdges += processTwoEdges(
				edge.line, p.edges[(i-1+edgesNumber)%edgesNumber], p.edges[(i+1)%edgesNumber],
			)
			continue
		}
		if Equal(edge.secondEnd.y, point.y) && point.x <= edge.secondEnd.x {
			cntEdges += processTwoEdges(
				horizontalLine, edge, p.edges[(i+1)%edgesNumber],
			)
			continue
		}
		if Equal(edge.firstEnd.y, point.y) {
			continue
		}
		x := (-edge.line.C - edge.line.B*point.y) / edge.line.A
		if betweenTwoVectors(edge.firstEnd, NewVector(x, point.y), edge.secondEnd) && point.x <= x {
			cntEdges++
		}
	}
	return cntEdges % 2
}

func (p *Polygon) checkBoundary(point *Vector) bool {
	for _, edge := range p.edges {
		if Equal(edge.line.A*point.x+edge.line.B*point.y+edge.line.C, 0) &&
			betweenTwoVectors(edge.firstEnd, point, edge.secondEnd) {
			return true
		}
	}
	return false
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)

	n, m := ScanInt(scanner), ScanInt(scanner)
	polygon := ScanPolygon(scanner, n)
	points := make([]*Vector, m)
	for i := 0; i < m; i++ {
		points[i] = ScanVector(scanner)
	}

	polygon.initEdges()
	for _, point := range points {
		if polygon.checkBoundary(point) {
			fmt.Println("BOUNDARY")
			continue
		}
		if polygon.isInside(point) == 1 {
			fmt.Println("INSIDE")
			continue
		}
		fmt.Println("OUTSIDE")
	}
}
