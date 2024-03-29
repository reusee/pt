package pt

import "cmp"

type Int int

func (i Int) Compare(i2 Int) int {
	return cmp.Compare(i, i2)
}
