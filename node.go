package pt

import (
	"fmt"
	"io"
)

type _Node[T Ordered[T]] struct {
	value    T
	priority _Priority
	left     *_Node[T]
	right    *_Node[T]
}

type Ordered[T any] interface {
	Compare(T) int
}

func (n *_Node[T]) Upsert(value T, priority _Priority, mutate bool) (ret *_Node[T], existed bool) {
	if n != nil {
		return n.upsertSlow(value, priority, mutate)
	}
	// new node
	return &_Node[T]{
		value:    value,
		priority: priority,
	}, false
}

func (n *_Node[T]) upsertSlow(value T, priority _Priority, mutate bool) (ret *_Node[T], existed bool) {
	switch value.Compare(n.value) {

	case -1:
		left, existed := n.left.Upsert(value, priority, mutate)
		return join(
			n,
			left,
			n.right,
			mutate,
		), existed

	case 1:
		right, existed := n.right.Upsert(value, priority, mutate)
		return join(
			n,
			n.left,
			right,
			mutate,
		), existed

	case 0:
		// exists
		if n.priority == priority {
			// same
			return n, true
		}
		if mutate {
			n.priority = priority
			return join(
				n,
				n.left,
				n.right,
				true,
			), true
		} else {
			return join(
				// new node
				&_Node[T]{
					value:    n.value,
					priority: priority,
					// setting these fields is not required for correctness, but will save a node allocation in join
					left:  n.left,
					right: n.right,
				},
				n.left,
				n.right,
				false,
			), true
		}

	}
	panic("bad Compare result") // NOCOVER
}

func join[T Ordered[T]](middle, left, right *_Node[T], mutate bool) *_Node[T] {
	if middle.priority == minPriority && left == nil && right == nil {
		// leaf node to be deleted
		return nil
	}

	if (left == nil || middle.priority >= left.priority) &&
		(right == nil || middle.priority >= right.priority) {
		// no rotation
		if middle.left == left && middle.right == right {
			// no change
			return middle
		}
		if mutate {
			// use middle
			middle.left = left
			middle.right = right
			return middle
		} else {
			// new node
			return &_Node[T]{
				value:    middle.value,
				priority: middle.priority,
				left:     left,
				right:    right,
			}
		}
	}

	if left != nil && left.priority > middle.priority && (right == nil || left.priority >= right.priority) {
		// rotate right
		if mutate {
			// use left
			left.right = join(
				middle,
				left.right,
				right,
				true,
			)
			return left
		} else {
			// new node
			return &_Node[T]{
				value:    left.value,
				priority: left.priority,
				left:     left.left,
				right: join(
					middle,
					left.right,
					right,
					false,
				),
			}
		}
	}

	if right != nil && right.priority > middle.priority && (left == nil || right.priority >= left.priority) {
		// rotate left
		if mutate {
			// use right
			right.left = join(
				middle,
				left,
				right.left,
				true,
			)
			return right
		} else {
			// new node
			return &_Node[T]{
				value:    right.value,
				priority: right.priority,
				left: join(
					middle,
					left,
					right.left,
					false,
				),
				right: right.right,
			}
		}
	}

	panic("impossible") // NOCOVER
}

func (n *_Node[T]) Height() int {
	if n == nil {
		return 0
	}
	return 1 + max(n.left.Height(), n.right.Height())
}

func (n *_Node[T]) Dump(out io.Writer, level int) {
	if n == nil {
		return
	}
	for i := 0; i < level; i++ {
		out.Write([]byte("\t"))
	}
	fmt.Fprintf(out, "value %v, priority %v\n", n.value, n.priority)
	n.left.Dump(out, level+1)
	n.right.Dump(out, level+1)
}

func (n *_Node[T]) Delete(value T, mutate bool) (ret *_Node[T], deleted bool) {
	return n.Upsert(value, minPriority, mutate)
}

func (n *_Node[T]) Split(value T, mutate bool) (ret *_Node[T], existed bool) {
	return n.Upsert(value, maxPriority, mutate)
}

func (n *_Node[T]) Union(n2 *_Node[T], mutate bool) *_Node[T] {
	if n2 == nil {
		return n
	}
	if n == nil {
		return n2
	}
	if n2.priority > n.priority {
		return n2.Union(n, mutate)
	}
	n2Split, _ := n2.Split(n.value, mutate)
	return join(
		n,
		n.left.Union(n2Split.left, mutate),
		n.right.Union(n2Split.right, mutate),
		mutate,
	)
}

func (n *_Node[T]) Length() int {
	if n == nil {
		return 0
	}
	return 1 + n.left.Length() + n.right.Length()
}

func (n *_Node[T]) Get(pivot T) (ret T, ok bool) {
	for {
		if n == nil {
			return
		}
		switch pivot.Compare(n.value) {
		case 0:
			return n.value, true
		case -1:
			n = n.left
		case 1:
			n = n.right
		default:
			panic("bad Compare result") // NOCOVER
		}
	}
}

func build[T Ordered[T]](source _PrioritySource, slice []T) *_Node[T] {
	if len(slice) == 0 {
		return nil
	}
	return buildSlow(source, slice)
}

func buildSlow[T Ordered[T]](source _PrioritySource, slice []T) *_Node[T] {
	i := len(slice) / 2
	left := slice[:i]
	right := slice[i+1:]
	ret := &_Node[T]{
		value:    slice[i],
		priority: source(),
		left:     build(source, left),
		right:    build(source, right),
	}
	heapify(ret)
	return ret
}

func heapify[T Ordered[T]](n *_Node[T]) {
	max := n
	if n.left != nil && n.left.priority > max.priority {
		max = n.left
	}
	if n.right != nil && n.right.priority > max.priority {
		max = n.right
	}
	if max != n {
		max.priority, n.priority = n.priority, max.priority
		heapify(max)
	}
}

func (t *_Node[T]) BulkDelete(t2 *_Node[T], mutate bool) *_Node[T] {
	if t2 == nil || t == nil {
		return t
	}
	t2Split, exists := t2.Split(t.value, mutate)
	if !exists {
		return join(
			t,
			t.left.BulkDelete(t2Split.left, mutate),
			t.right.BulkDelete(t2Split.right, mutate),
			mutate,
		)
	}
	if mutate {
		t.priority = minPriority
		return join(
			t,
			t.left.BulkDelete(t2Split.left, mutate),
			t.right.BulkDelete(t2Split.right, mutate),
			mutate,
		)
	}
	return join(
		&_Node[T]{
			value:    t.value,
			priority: minPriority,
			left:     t.left,
			right:    t.right,
		},
		t.left.BulkDelete(t2Split.left, mutate),
		t.right.BulkDelete(t2Split.right, mutate),
		mutate,
	)
}
