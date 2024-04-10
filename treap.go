package pt

import "sync/atomic"

type Treap[T Ordered[T]] struct {
	root           atomic.Pointer[_Node[T]]
	prioritySource _PrioritySource
}

func NewTreap[T Ordered[T]]() *Treap[T] {
	treap := &Treap[T]{
		prioritySource: newPrioritySource(),
	}
	treap.root.Store((*_Node[T])(nil))
	return treap
}

func (t *Treap[T]) Upsert(value T) (existed bool) {
	for {
		root := t.root.Load()
		newRoot, existed := root.Upsert(value, t.prioritySource(), false)
		if t.root.CompareAndSwap(root, newRoot) {
			return existed
		}
	}
}

func (t *Treap[T]) Remove(value T) (removed bool) {
	for {
		root := t.root.Load()
		newRoot, removed := root.Remove(value, false)
		if t.root.CompareAndSwap(root, newRoot) {
			return removed
		}
	}
}

func (t *Treap[T]) Union(t2 *Treap[T]) {
	for {
		root := t.root.Load()
		newRoot := root.Union(t2.root.Load(), false)
		if t.root.CompareAndSwap(root, newRoot) {
			return
		}
	}
}

func (t *Treap[T]) Length() int {
	return t.root.Load().Length()
}

func (t *Treap[T]) Get(pivot T) (ret T, ok bool) {
	return t.root.Load().Get(pivot)
}

func (t *Treap[T]) BulkRemove(values []T) {
	node := build(t.prioritySource, sortUnique(values, T.Compare))
	for {
		root := t.root.Load()
		newRoot := root.BulkRemove(node, false)
		if t.root.CompareAndSwap(root, newRoot) {
			return
		}
	}
}

func (t *Treap[T]) BulkUpsert(values []T) {
	node := build(t.prioritySource, sortUnique(values, T.Compare))
	for {
		root := t.root.Load()
		newRoot := root.Union(node, false)
		if t.root.CompareAndSwap(root, newRoot) {
			return
		}
	}
}
