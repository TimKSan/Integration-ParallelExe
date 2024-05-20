package main

import "fmt"

func main() {

	arr1 := [4]int{3, 1, 4, 2}
	arr2 := [5]int{9, 5, 8, 6, 7}

	mergedArr3 := mergeArr(arr1[:], arr2[:])
	fmt.Println(mergedArr3) //[3 1 4 2 9 5 8 6 7]
	sortedFinalArr := sortBubble(mergedArr3)
	fmt.Println(sortedFinalArr) //[1 2 3 4 5 6 7 8 9]
}

// Слияние массивов
func mergeArr(arr1 []int, arr2 []int) []int {
	arr3 := [9]int{}
	copy(arr3[:], arr1)
	copy(arr3[len(arr1):], arr2)

	return arr3[:]
}

// Сортировка пузырьком
func sortBubble(arr []int) []int {
	for i := 0; i < len(arr)-1; i++ {
		for j := 0; j < len(arr)-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}
