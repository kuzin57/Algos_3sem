package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
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

func euclidCoeffs(a, b, module int) (int, int) {
	var (
		prevX, prevY = 0, 0
		x, y         = 1, -(a / b)
	)
	a, b = b, a%b
	for b != 0 {
		q := a / b
		if a%b != 0 {
			tmpX, tmpY := x, y
			x = (prevX - x*q) % module
			y = (prevY - y*q) % module
			prevX, prevY = tmpX, tmpY
		}
		a, b = b, a%b
	}
	return x, y
}

type modular struct {
	value  int
	module int
}

type modularFabric struct {
	module int
}

func initFabric(module int) modularFabric {
	return modularFabric{module: module}
}

func (f modularFabric) buildModular(number int) modular {
	ans := modular{value: number % f.module}
	ans.normalize()
	return ans
}

func (m *modular) normalize() {
	if m.value >= 0 {
		m.value %= m.module
		return
	}
	m.value = (m.value + m.module*(-m.value/m.module+1)) % m.module
}

func sumModulars(first modular, second modular) modular {
	ans := modular{value: first.value + second.value}
	ans.normalize()
	return ans
}

func substrModulars(first modular, second modular) modular {
	ans := modular{value: first.value - second.value}
	ans.normalize()
	return ans
}

func multModulars(first modular, second modular) modular {
	ans := modular{value: first.value * second.value}
	ans.normalize()
	return ans
}

func findInvert(m modular) modular {
	_, invert := euclidCoeffs(m.module, m.value, m.module)
	ans := modular{value: invert}
	ans.normalize()
	return ans
}

func (m modular) String() string {
	return fmt.Sprintf("%d", m.value)
}

func ScanModular(scanner *bufio.Scanner, f modularFabric) modular {
	return f.buildModular(ScanInt(scanner))
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Split(splitFunc)

	fabric := initFabric(1000000007)
	var (
		a        = ScanModular(scanner, fabric)
		b        = ScanModular(scanner, fabric)
		c        = ScanModular(scanner, fabric)
		d        = ScanModular(scanner, fabric)
		inverseB modular
		inverseD modular
	)

	a.normalize()
	b.normalize()
	c.normalize()
	d.normalize()
	inverseB = findInvert(b)
	inverseD = findInvert(d)
	fmt.Println(sumModulars(multModulars(a, inverseB), multModulars(c, inverseD)))
}
