package pt

type node[T ordered[T]] struct {
	value    T
	priority Priority
	left     *node[T]
	right    *node[T]
}

func (n *node[T]) upsert(value T, priority Priority) *node[T] {
	if n != nil {
		return n.upsertSlow(value, priority)
	}
	return &node[T]{
		value:    value,
		priority: priority,
	}
}

func (n *node[T]) upsertSlow(value T, priority Priority) *node[T] {
	switch value.Compare(n.value) {

	case -1:
		return join(
			n,
			n.left.upsert(value, priority),
			n.right,
		)

	case 1:
		return join(
			n,
			n.left,
			n.right.upsert(value, priority),
		)

	case 0:
		// exists
		if n.priority == priority {
			// same
			return n
		}
		return join(
			&node[T]{
				value:    n.value,
				priority: priority,
			},
			n.left,
			n.right,
		)

	}
	panic("bad Compare result")
}

func join[T ordered[T]](middle, left, right *node[T]) *node[T] {
	if (left == nil || middle.priority >= left.priority) &&
		(right == nil || middle.priority >= right.priority) {
		// no rotation
		return &node[T]{
			value:    middle.value,
			priority: middle.priority,
			left:     left,
			right:    right,
		}
	}

	if left != nil && left.priority > middle.priority && (right == nil || left.priority > right.priority) {
		// rotate right
		return &node[T]{
			value:    left.value,
			priority: left.priority,
			left:     left.left,
			right: join(
				middle,
				left.right,
				right,
			),
		}
	}

	if right != nil && right.priority > middle.priority && (left == nil || right.priority > left.priority) {
		// rotate left
		return &node[T]{
			value:    right.value,
			priority: right.priority,
			left: join(
				middle,
				left,
				right.left,
			),
			right: right.right,
		}
	}

	panic("impossible")
}
