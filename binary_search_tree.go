package main

import "fmt"

type Node struct {
	Left  *Node
	Right *Node
	Load  *Load
}

func (current *Node) Min() *Node {
	if current == nil {
		return current
	}

	for current.Left != nil {
		current = current.Left
	}

	return current
}

func (current *Node) Insert(node *Node) *Node {
	if current == nil {
		return node
	}

	if current.Load.Less(node.Load) {
		current.Right = current.Right.Insert(node)
	} else {
		current.Left = current.Left.Insert(node)
	}

	return current
}

func (current *Node) Delete(node *Node) *Node {
	if current == nil {
		return current
	}

	if current.Load.Less(node.Load) {
		current.Right = current.Right.Delete(node)
	} else if node.Load.Less(current.Load) {
		current.Left = current.Left.Delete(node)
	} else {
		if current.Left == nil {
			return current.Right
		} else if current.Right == nil {
			return current.Left
		}

		minNode := current.Right.Min()
		current.Load = minNode.Load
		current.Right = current.Right.Delete(minNode)
	}

	return current
}

func (n *Node) Search(point *Point) *Node {
	previous := n
	current := n

	for current != nil {
		previous = current

		if current.Load.pickup.Less(point) {
			current = current.Right
		} else if point.Less(current.Load.pickup) {
			current = current.Left
		} else {
			return current
		}
	}

	return previous
}

func (n *Node) String() string {
	var s string

	if n.Left != nil {
		s += fmt.Sprintf("<%s>", n.Left.Load.id)
	}

	s += fmt.Sprintf(" %s ", n.Load.id)

	if n.Right != nil {
		s += fmt.Sprintf("<%s>", n.Right.Load.id)
	}

	return s
}
