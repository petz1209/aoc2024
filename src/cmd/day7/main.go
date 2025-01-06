/*
ein kombinatorik problem
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

func main() {
	res1 := puzzle1(INPUT)
	fmt.Println("---------------------------------------------")
	fmt.Printf("day6-puzzle1 result: %d\n", res1)
	fmt.Println("---------------------------------------------")

	res2 := puzzle2(INPUT)
	fmt.Println("---------------------------------------------")
	fmt.Printf("day6-puzzle2 result: %d\n", res2)
	fmt.Println("---------------------------------------------")

}

func puzzle1(inputPath string) int {
	data := readInput(inputPath)
	result := 0
	for _, d := range data {
		//fmt.Printf("expected result: %d   inputs: %v\n", d[0], d[1:])
		passed := testEngine(d[0], d[1:], []byte{'*', '+'})
		if passed {
			result += d[0]
		}
	}
	return result
}

func puzzle2(inputPath string) int {
	data := readInput(inputPath)
	result := 0
	for _, d := range data {
		//fmt.Printf("expected result: %d   inputs: %v\n", d[0], d[1:])
		passed := testEngine(d[0], d[1:], []byte{'*', '+', '|'})
		if passed {
			result += d[0]
		}
	}
	return result
}

func testEngine(expected int, inputs []int, calculations []byte) bool {
	// somehow I need to figure out how I can iterate on different calc operators throughout the stack

	// based on the length of the inputs I need to create some rule system

	//caluclationConfigs := make([][]byte, 0)
	combinations := generateCombinations(len(inputs)-1, calculations)
	for _, combo := range combinations {
		r := calculateRound(inputs, combo)
		if r == expected {
			return true
		}
	}
	return false

}

func generateCombinations(length int, states []byte) [][]byte {
	var combinations [][]byte

	// Backtracking-Funktion
	var backtrack func(current []byte)
	backtrack = func(current []byte) {
		if len(current) == length {
			// Eine Kopie von `current` hinzufügen
			temp := make([]byte, len(current))
			copy(temp, current)
			combinations = append(combinations, temp)
			return
		}
		for _, state := range states {
			current = append(current, state)   // Zustand hinzufügen
			backtrack(current)                 // Rekursion
			current = current[:len(current)-1] // Zustand entfernen
		}
	}

	backtrack([]byte{})
	return combinations
}

func calculateRound(inputs []int, calcs []byte) int {

	res := 0
	for i, v := range inputs {
		if i == 0 {
			res = v
			continue
		}
		switch calcs[i-1] {
		case '*':
			res *= v
		case '+':
			res += v
		case '|':
			new := fmt.Sprintf("%d%d", res, v)
			res, _ = strconv.Atoi(new)
		}
	}
	return res

}

func readInput(path string) [][]int {

	result := make([][]int, 0)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		row := strings.Fields(scanner.Text())
		tmp := make([]int, len(row))
		for i := range row {
			v, _ := strconv.Atoi(strings.TrimSuffix(row[i], ":"))
			tmp[i] = v
		}
		result = append(result, tmp)
	}
	return result

}
