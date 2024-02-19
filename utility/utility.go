package utility

import "sort"

func UpperBoundInt(ar []int, num int) (pos int) {
	return sort.Search(len(ar), func(i int) bool {
		return !(ar[i] <= num)
	})
}

func LowerBoundInt(ar []int, num int) (pos int) {
	return sort.Search(len(ar), func(i int) bool {
		return !(ar[i] < num)
	})
}
