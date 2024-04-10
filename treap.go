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

func (t *Treap[T]) Delete(value T) (deleted bool) {
	for {
		root := t.root.Load()
		newRoot, deleted := root.Delete(value, false)
		if t.root.CompareAndSwap(root, newRoot) {
			return deleted
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

func (t *Treap[T]) BulkDelete(values []T) {
	node := build(t.prioritySource, sortUnique(values, T.Compare))
	for {
		root := t.root.Load()
		newRoot := root.BulkDelete(node, false)
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
