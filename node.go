package pt

type node[T ordered[T]] struct {
	value    T
	priority int
	left     *node[T]
	right    *node[T]
}
