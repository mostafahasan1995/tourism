package main

import "fmt"

func test(datafetched int) {

	df := datafetched

	df = 10
	fmt.Println(datafetched, df)

}

func SliceChunk[T any](arr []T, size int) [][]T {
	if size == 0 {
		return [][]T{arr}
	}
	c := 0
	result := [][]T{}
	for {
		var a []T
		if c+size > len(arr) {
			a = arr[c:]
		} else {
			a = arr[c : c+size]
		}
		c += len(a)
		result = append(result, a)
		if c == len(arr) {
			break
		}
	}

	return result

}

func main() {
	arr := [][]int{
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4, 5},
	}

	result := SliceChunk(arr, 3)

	fmt.Println(result)
}
