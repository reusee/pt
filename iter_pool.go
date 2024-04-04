package pt

import (
	"reflect"
	"sync"
)

var iterPool = new(sync.Map) // reflect.Type -> *sync.Pool

func getIter[T Ordered[T]]() *Iter[T] {
	t := reflect.TypeFor[T]()
	v, ok := iterPool.Load(t)
	if ok {
		iter := v.(*sync.Pool).Get().(*Iter[T])
		return iter
	}
	v, _ = iterPool.LoadOrStore(t, &sync.Pool{
		New: func() any {
			return new(Iter[T])
		},
	})
	pool := v.(*sync.Pool)
	iter := pool.Get().(*Iter[T])
	return iter
}

func putIter[T Ordered[T]](iter *Iter[T]) {
	iter.stack = iter.stack[:0]
	t := reflect.TypeFor[T]()
	v, ok := iterPool.Load(t)
	if ok {
		v.(*sync.Pool).Put(iter)
		return
	}
	v, _ = iterPool.LoadOrStore(t, &sync.Pool{
		New: func() any {
			return new(Iter[T])
		},
	})
	pool := v.(*sync.Pool)
	pool.Put(iter)
}
