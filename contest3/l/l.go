package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"
)

var (
	errorVectorsEqual = fmt.Errorf("vectors are equal")
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

func NewVector(x, y int) *Vector {
	return &Vector{x: x, y: y}
}

type Polygon struct {
	vertices []*Vector
	edges    []*Edge
}

func ScanVector(scanner *bufio.Scanner) *Vector {
	return &Vector{x: int(ScanInt(scanner)), y: int(ScanInt(scanner))}
}

func ScanPolygon(scanner *bufio.Scanner, vertices int) *Polygon {
	newPolynom := &Polygon{}
	for i := 0; i < vertices; i++ {
		newPolynom.vertices = append(newPolynom.vertices, ScanVector(scanner))
	}
	return newPolynom
}

type Edge struct {
	line       *Line
	firstEnd   *Vector
	secondEnd  *Vector
	edgeVector *Vector
}

func NewEdge(line *Line, firstEnd, secondEnd, edgeVector *Vector) *Edge {
	return &Edge{
		line:       line,
		firstEnd:   firstEnd,
		secondEnd:  secondEnd,
		edgeVector: edgeVector,
	}
}

type Line struct {
	A, B, C int
}

func buildLine(firstVector, secondVector *Vector) (*Line, error) {
	if Equal(firstVector.x, secondVector.x) && Equal(secondVector.y, firstVector.y) {
		return nil, errorVectorsEqual
	}

	return &Line{
		A: firstVector.y - secondVector.y,
		B: secondVector.x - firstVector.x,
		C: firstVector.x*secondVector.y -
			secondVector.x*firstVector.y,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func betweenTwoNumbers(a, b, c int) bool {
	return min(a, b) <= c && c <= max(a, b)
}

func Equal(a, b int) bool {
	return a == b
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

		newEdge := NewEdge(
			newLine,
			p.vertices[i],
			p.vertices[(i+1)%verticesNumber],
			subVectors(p.vertices[(i+1)%verticesNumber], p.vertices[i]),
		)
		p.edges = append(p.edges, newEdge)
	}
}

func reverse[Type *Vector | *Edge](arr []Type) {
	left, right := 0, len(arr)-1
	for left < right {
		arr[left], arr[right] = arr[right], arr[left]
		left++
		right--
	}
}

func (p *Polygon) defineClockwise() {
	N := len(p.vertices)
	var sum int
	for i := 0; i < N; i++ {
		sum += int(p.vertices[(i+1)%N].x-p.vertices[i].x) * int(p.vertices[(i+1)%N].y+p.vertices[i].y)
	}
	if sum <= 0 {
		return
	}
	reverse(p.vertices)
}

func (p *Polygon) sortVertices() {
	sort.Slice(p.vertices, func(i, j int) bool {
		return p.vertices[i].x < p.vertices[j].x ||
			(p.vertices[i].x == p.vertices[j].x &&
				p.vertices[i].y < p.vertices[j].y)
	})
}

func twoVectorsEqual(first, second *Vector) bool {
	return Equal(first.x, second.x) && Equal(first.y, second.y)
}

func subVectors(first, second *Vector) *Vector {
	return &Vector{
		x: first.x - second.x,
		y: first.y - second.y,
	}
}

func crossProduct(first, second *Vector) int {
	return first.x*second.y - first.y*second.x
}

func (p *Polygon) findEdgeWithFirstEnd(firstEnd *Vector) int {
	for i, edge := range p.edges {
		if twoVectorsEqual(edge.firstEnd, firstEnd) {
			return i
		}
	}
	return -1
}

func isInsideTriangle(point, v1, v2, v3 *Vector) bool {
	d1 := crossProduct(subVectors(point, v2), subVectors(v1, v2))
	d2 := crossProduct(subVectors(point, v3), subVectors(v2, v3))
	d3 := crossProduct(subVectors(point, v1), subVectors(v3, v1))

	hasNegative := (d1 < 0) || (d2 < 0) || (d3 < 0)
	hasPositive := (d1 > 0) || (d2 > 0) || (d3 > 0)

	return !(hasNegative && hasPositive)
}

func checkOneLine(first, second, third *Vector) bool {
	line, _ := buildLine(first, third)
	return Equal(line.A*second.x+line.B*second.y+line.C, 0)
}

func (p *Polygon) isPointInside(point *Vector, mostLeftPoint int) bool {
	left := 1
	right := len(p.vertices) - 1
	for right-left > 1 {
		mid := (left + right) / 2
		line, _ := buildLine(p.vertices[mostLeftPoint], p.vertices[(mostLeftPoint+mid)%len(p.vertices)])
		if line.A*point.x+line.B*point.y+line.C < 0 {
			right = mid
		} else {
			left = mid
		}
	}

	if checkOneLine(
		p.vertices[mostLeftPoint],
		p.vertices[(mostLeftPoint+left)%len(p.vertices)],
		p.vertices[(mostLeftPoint+right)%len(p.vertices)],
	) {
		return checkOneLine(p.vertices[mostLeftPoint], point, p.vertices[(mostLeftPoint+right)%len(p.vertices)]) &&
			betweenTwoVectors(
				p.vertices[mostLeftPoint],
				point,
				p.vertices[(mostLeftPoint+right)%len(p.vertices)],
			)
	}

	if Equal(point.x, p.vertices[mostLeftPoint].x) && Equal(point.y, p.vertices[mostLeftPoint].y) {
		return true
	}

	return isInsideTriangle(
		point,
		p.vertices[mostLeftPoint],
		p.vertices[(mostLeftPoint+left)%len(p.vertices)],
		p.vertices[(mostLeftPoint+right)%len(p.vertices)],
	)
}

func sumVectors(first, second *Vector) *Vector {
	return &Vector{
		x: first.x + second.x,
		y: first.y + second.y,
	}
}

func sumPolygons(first, second *Polygon) *Polygon {
	result := &Polygon{}
	defer result.initEdges()

	firstPtr := first.findEdgeWithFirstEnd(first.vertices[0])
	secondPtr := second.findEdgeWithFirstEnd(second.vertices[0])

	result.vertices = append(
		result.vertices,
		sumVectors(first.vertices[0], second.vertices[0]),
	)

	addVertexByVector := func(vector *Vector) {
		result.vertices = append(
			result.vertices,
			sumVectors(result.vertices[len(result.vertices)-1], vector),
		)
	}

	incByModule := func(i *int, module int) {
		*i = (*i + 1) % module
	}

	firstFinish, secondFinish := firstPtr, secondPtr
	var startedFirst, startedSecond bool
	for {
		if startedFirst && startedSecond && firstPtr == firstFinish && secondPtr == secondFinish {
			break
		}

		if startedFirst && firstPtr == firstFinish {
			addVertexByVector(second.edges[secondPtr].edgeVector)
			incByModule(&secondPtr, len(second.edges))
			continue
		}

		if startedSecond && secondPtr == secondFinish {
			addVertexByVector(first.edges[firstPtr].edgeVector)
			incByModule(&firstPtr, len(first.edges))
			continue
		}

		crossProd := crossProduct(first.edges[firstPtr].edgeVector, second.edges[secondPtr].edgeVector)
		switch {
		case crossProd == 0:
			addVertexByVector(sumVectors(first.edges[firstPtr].edgeVector, second.edges[secondPtr].edgeVector))
			incByModule(&firstPtr, len(first.edges))
			incByModule(&secondPtr, len(second.edges))
			startedFirst, startedSecond = true, true
		case crossProd > 0:
			addVertexByVector(first.edges[firstPtr].edgeVector)
			incByModule(&firstPtr, len(first.edges))
			startedFirst = true
		case crossProd < 0:
			addVertexByVector(second.edges[secondPtr].edgeVector)
			incByModule(&secondPtr, len(second.edges))
			startedSecond = true
		}
	}
	result.vertices = result.vertices[:len(result.vertices)-1]
	return result
}

func (p *Polygon) findExtremeVertex() int {
	var index int
	curPoint := p.vertices[0]
	for i, vertex := range p.vertices {
		if vertex.x < curPoint.x ||
			(vertex.x == curPoint.x && vertex.y < curPoint.y) {
			curPoint = vertex
			index = i
		}
	}
	return index
}

func run(scanner *bufio.Scanner) {
	n := ScanInt(scanner)
	firstPolygon := ScanPolygon(scanner, n)
	firstPolygon.defineClockwise()
	firstPolygon.initEdges()

	n = ScanInt(scanner)
	secondPolygon := ScanPolygon(scanner, n)
	secondPolygon.defineClockwise()
	secondPolygon.initEdges()

	n = ScanInt(scanner)
	thirdPolygon := ScanPolygon(scanner, n)
	thirdPolygon.defineClockwise()
	thirdPolygon.initEdges()

	firstPolygon.sortVertices()
	secondPolygon.sortVertices()
	thirdPolygon.sortVertices()

	sumPolygons := sumPolygons(sumPolygons(firstPolygon, secondPolygon), thirdPolygon)
	mostLeftPoint := sumPolygons.findExtremeVertex()
	q := ScanInt(scanner)
	for t := 0; t < q; t++ {
		point := ScanVector(scanner)
		point.x *= 3
		point.y *= 3
		if sumPolygons.isPointInside(point, mostLeftPoint) {
			fmt.Println("YES")
			continue
		}
		fmt.Println("NO")
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)
	run(scanner)
}
