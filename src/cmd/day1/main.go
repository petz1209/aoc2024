/*
Day1
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const INPUT = "input.txt"

func readInput() ([]int, []int) {
	f, err := os.Open(INPUT)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	left := make([]int, 0)
	right := make([]int, 0)
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), "   ")
		toSlice(&left, row[0])
		toSlice(&right, row[1])
	}
	return left, right
}

func quickSort(arr []int, low, high int) []int {
	if low < high {
		var p int
		arr, p = partition(arr, low, high)
		arr = quickSort(arr, low, p-1)
		arr = quickSort(arr, p+1, high)
	}
	return arr
}
func quickSortStart(arr []int) []int {
	return quickSort(arr, 0, len(arr)-1)
}

func partition(arr []int, low, high int) ([]int, int) {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return arr, i
}

func toSlice(arr *[]int, str string) {
	num, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	*arr = append(*arr, num)
}

func toSet(arr []int) []int {
	set := make([]int, 0)
	for idx := range arr {
		if idx == 0 {
			set = append(set, arr[idx])
		}
		if arr[idx] > set[len(set)-1] {
			set = append(set, arr[idx])
		}
	}
	return set
}

func countOccurences(val int, arr []int) int {
	c := 0
	for _, v := range arr {
		if v == val {
			c++
		}
	}
	return c

}

func main() {
	/*
		to achieve this task, we first split the input into 2 lists, left and right
		then we sort those lists ascending, using quicksort algorithm
		then we iterate and over the list index and compare the values

		2nd
		now we create a set of right and check the occurences
	*/

	left, right := readInput()
	// sort
	left = quickSortStart(left)
	right = quickSortStart(right)
	// cacluate for ch1
	totalDiff := 0
	for idx := range left {
		if len(right)-1 < idx {
			break
		}
		if left[idx] > right[idx] {
			totalDiff += (left[idx] - right[idx])
		} else {
			totalDiff += (right[idx] - left[idx])
		}
	}
	fmt.Println("-------------------------------------------")
	fmt.Printf("result: %d\n", totalDiff)
	fmt.Println("-------------------------------------------")

	// For each number in the left list, find out how many times it occurs in the right one
	lSet := toSet(left)
	// calculate for c2
	res := 0
	for _, v := range lSet {
		c := countOccurences(v, right)
		res += (v * c)
	}
	fmt.Println("-------------------------------------------")
	fmt.Printf("result2: %d\n", res)
	fmt.Println("-------------------------------------------")

}
