package pt

type node[T ordered[T]] struct {
	value    T
	priority int64
	left     *node[T]
	right    *node[T]
}

func (n *node[T]) insert(value T, priority int64) *node[T] {
	if n == nil {
		return &node[T]{
			value:    value,
			priority: priority,
		}
	}

	switch value.Compare(n.value) {

	case 0:
		// exists
		if n.priority == priority {
			// same
			return n
		}

		ret := *n
		ret.priority = priority
		if (ret.left == nil || ret.left.priority <= ret.priority) &&
			(ret.right == nil || ret.right.priority <= ret.priority) {
			// no rotation
			return &ret
		}

		panic("fixme")

	case -1:
		// insert left
		left := n.left.insert(value, priority)
		if left == n.left {
			// no insert
			return n
		}
		if n.priority >= left.priority {
			// no rotation
			return &node[T]{
				value:    n.value,
				priority: n.priority,
				left:     left,
				right:    n.right,
			}
		}
		// rotate
		return &node[T]{
			value:    left.value,
			priority: left.priority,
			left:     left.left,
			right: &node[T]{
				value:    n.value,
				priority: n.priority,
				left:     left.right,
				right:    n.right,
			},
		}

	case 1:
		// insert right
		right := n.right.insert(value, priority)
		if right == n.right {
			// no insert
			return n
		}
		if n.priority >= right.priority {
			// no rotation
			return &node[T]{
				value:    n.value,
				priority: n.priority,
				left:     n.left,
				right:    right,
			}
		}
		// rotate
		return &node[T]{
			value:    right.value,
			priority: right.priority,
			left: &node[T]{
				value:    n.value,
				priority: n.priority,
				left:     n.left,
				right:    right.left,
			},
			right: right.right,
		}

	}
	panic("bad Compare result")
}
