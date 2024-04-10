package pt

import "sort"

func sortUnique[T any](slice []T, cmp func(T, T) int) []T {
	if len(slice) == 0 {
		return nil
	}
	sort.Slice(slice, func(i, j int) bool {
		return cmp(slice[i], slice[j]) < 0
	})
	res := slice[:0]
	for i := 0; i < len(slice)-1; i++ {
		if cmp(slice[i], slice[i+1]) != 0 {
			res = append(res, slice[i])
		}
	}
	res = append(res, slice[len(slice)-1])
	return res
}
