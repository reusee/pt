package pt

type node[T ordered[T]] struct {
	value    T
	priority int
	left     *node[T]
	right    *node[T]
}

func (n *node[T]) split(pivot *node[T]) (left, middle, right *node[T]) {
	if n == nil {
		return nil, nil, nil
	}

	switch pivot.value.Compare(n.value) {

	case 0:
		return n.left, n, n.right

	case -1:
		// pivot in left
		left, middle, right = n.left.split(pivot)
		// left < middle < (right < n < n.right)
		// merge right and n to a new node
		newNode := *n
		newNode.left = right
		return left, middle, &newNode

	case 1:
		// pivot in right
		left, middle, right = n.right.split(pivot)
		// (n.left < n < left) < middle < right
		// merge n and left to a new node
		newNode := *n
		newNode.right = left
		return &newNode, middle, right

	default:
		panic("bad compare result")
	}

}
