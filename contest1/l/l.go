package main

import (
	"bufio"
	"fmt"
	"os"
)

type node struct {
	next        []int
	link        int
	len         int
	left, right int
	parent      int
}

func createNode(len int) *node {
	ans := &node{next: make([]int, 27), len: len, link: -1, parent: -1}
	for i := 0; i < 27; i++ {
		ans.next[i] = -1
	}
	return ans
}

func (n *node) get(b byte) int {
	if b == '$' {
		return n.next[0]
	}
	return n.next[b-'a'+1]
}

func (n *node) set(b byte, val int) {
	if b == '$' {
		n.next[0] = val
		return
	}
	n.next[b-'a'+1] = val
}

type Automata struct {
	nodes []*node
	last  int
}

func NewAutomata() *Automata {
	start := createNode(0)
	automata := &Automata{nodes: make([]*node, 0), last: 0}
	automata.nodes = append(automata.nodes, start)
	return automata
}

func (a *Automata) Add(c byte, cnt int) {
	var exit bool
	cur := createNode(a.nodes[a.last].len + 1)
	cur.left = a.nodes[a.last].left
	cur.right = a.nodes[a.last].right + 1
	a.nodes = append(a.nodes, cur)
	curNode := a.last
	a.last = len(a.nodes) - 1
	for curNode != -1 {
		to := a.nodes[curNode].get(c)
		switch to {
		case -1:
			a.nodes[curNode].set(c, len(a.nodes)-1)
		default:
			if a.nodes[to].len == a.nodes[curNode].len+1 {
				cur.link = to
				exit = true
			} else {
				newNode := createNode(a.nodes[curNode].len + 1)
				for i, nodeTo := range a.nodes[to].next {
					if nodeTo == -1 {
						continue
					}
					if i == 0 {
						newNode.set(byte('$'), nodeTo)
						continue
					}
					newNode.set(byte('a'+i-1), nodeTo)
				}
				a.nodes = append(a.nodes, newNode)
				newNode.link = a.nodes[to].link
				a.nodes[to].link = len(a.nodes) - 1
				cur.link = len(a.nodes) - 1
				suffLink := curNode
				var flag bool
				for suffLink != -1 {
					ok := a.nodes[suffLink].get(c)
					if ok == to {
						if !flag {
							if suffLink != 0 {
								newNode.left = cur.right - newNode.len
								newNode.right = cur.right
							} else {
								newNode.left = cur.right - 1
								newNode.right = cur.right
							}
							flag = true
						}
						a.nodes[suffLink].set(c, len(a.nodes)-1)
					} else {
						break
					}
					suffLink = a.nodes[suffLink].link
				}
				exit = true
			}
		}
		if exit {
			break
		}
		curNode = a.nodes[curNode].link
	}
	if curNode == -1 {
		cur.link = 0
	}
}

func reverse(line string) string {
	var (
		left  = 0
		right = len(line) - 1
		bytes = []byte(line)
	)
	for left < right {
		bytes[left], bytes[right] = bytes[right], bytes[left]
		left++
		right--
	}
	return string(bytes)
}

func run(line string) {
	automata := NewAutomata()
	rline := reverse(line)
	for i, b := range rline {
		automata.Add(byte(b), i)
	}
	newAutomata := NewAutomata()
	for i := 0; i < len(automata.nodes)-1; i++ {
		newAutomata.nodes = append(newAutomata.nodes, createNode(0))
	}

	for i, node := range automata.nodes {
		linkNode := node.link
		if linkNode == -1 {
			continue
		}
		start := len(line) - (node.left + node.len - automata.nodes[linkNode].len - 1) - 1
		newAutomata.nodes[linkNode].set(line[start], i)
		newAutomata.nodes[i].left = start
		newAutomata.nodes[i].len = node.len - automata.nodes[linkNode].len
		newAutomata.nodes[i].parent = linkNode
	}
	used := make([]bool, len(newAutomata.nodes))
	fmt.Println(len(newAutomata.nodes))
	var a int
	numbers := make([]int, len(newAutomata.nodes))
	display(0, newAutomata, used, numbers, &a)
}

func display(cur int, automata *Automata, used []bool, numbers []int, cnt *int) {
	used[cur] = true
	curNode := automata.nodes[cur]
	numbers[cur] = *cnt
	(*cnt)++
	if curNode.parent != -1 {
		fmt.Println(numbers[curNode.parent], curNode.left, curNode.left+curNode.len)
	}
	for i := range curNode.next {
		if curNode.next[i] != -1 && !used[curNode.next[i]] {
			display(curNode.next[i], automata, used, numbers, cnt)
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(nil, 1<<30)
	scanner.Scan()
	line := scanner.Text()
	run(line)
}
