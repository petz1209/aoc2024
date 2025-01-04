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
	// test case description for task2 in aoc is wrong. They made a mistake.
	// The correct result for input_test.txt = 105
	solvePuzzle2()

}

func solvePuzzle1() {
	rules, instructions := readInput()
	rulesMap := NewRulesMap(rules)
	validInstructionIndexes, _ := validateAllInstructions(rulesMap, instructions)
	validInstructions := make([][]int, len(validInstructionIndexes))
	for idx, i := range validInstructionIndexes {
		validInstructions[idx] = instructions[i]
	}
	checkSum := calculateCheckSum(validInstructions)
	fmt.Println("----------------------------------")
	fmt.Printf("result: %d\n", checkSum)
	fmt.Println("----------------------------------")

}

func solvePuzzle2() {
	rules, instructions := readInput()
	rulesMap := NewRulesMap(rules)
	_, ivIdxs := validateAllInstructions(rulesMap, instructions)

	invalidInstructions := make([][]int, len(ivIdxs))
	for idx, i := range ivIdxs {
		invalidInstructions[idx] = instructions[i]
	}

	for idx := range invalidInstructions {
		sortInstruction(rulesMap, &invalidInstructions[idx])
	}

	checkSum := calculateCheckSum(invalidInstructions)
	fmt.Println("----------------------------------")
	fmt.Printf("result: %d\n", checkSum)
	fmt.Println("----------------------------------")

}

func readInput() ([][2]int, [][]int) {
	f, err := os.Open(INPUT)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	rules := make([][2]int, 0)
	instructions := make([][]int, 0)
	scanner := bufio.NewScanner(f)
	firstPart := true
	for scanner.Scan() {

		row := scanner.Text()

		switch firstPart {
		case true:
			firstPart = newRule(&rules, row)
		default:
			newInstruction(&instructions, row)
		}

	}
	return rules, instructions

}

func newRule(rules *[][2]int, row string) bool {
	if row == "" {
		return false
	}
	tmp := strings.Split(row, "|")

	n1, _ := strconv.Atoi(tmp[0])
	n2, _ := strconv.Atoi(tmp[1])
	rule := [2]int{n1, n2}
	*rules = append(*rules, rule)
	return true
}

func newInstruction(instructions *[][]int, row string) {
	tmp := strings.Split(row, ",")
	instruction := make([]int, len(tmp))
	for i, v := range tmp {
		num, _ := strconv.Atoi(v)
		instruction[i] = num
	}
	*instructions = append(*instructions, instruction)
}

func NewRulesMap(rules [][2]int) map[int][]int {
	m := make(map[int][]int, 0)
	for _, r := range rules {

		// if it is the first asigned follower
		if m[r[0]] == nil {
			v := []int{r[1]}
			m[r[0]] = v
		} else {
			m[r[0]] = append(m[r[0]], r[1])
		}
	}
	return m

}

func validateAllInstructions(rules map[int][]int, instructions [][]int) ([]int, []int) {
	validInstructions := make([]int, 0)
	invalidInstructions := make([]int, 0)
	for idx, inst := range instructions {
		if instructionValid(rules, inst) {
			validInstructions = append(validInstructions, idx)
		} else {
			invalidInstructions = append(invalidInstructions, idx)
		}
	}
	return validInstructions, invalidInstructions

}

func instructionValid(rules map[int][]int, instruction []int) bool {
	for idx := range instruction {
		if idx == 0 {
			continue
		}
		/*
			if rulesMap[inst[idx]] == nil {
				continue
			}*/
		for _, follower := range rules[instruction[idx]] {
			if instruction[idx-1] == follower {
				return false
			}
		}
	}
	return true
}

func positionValid(rules map[int][]int, instruction []int, idx int) bool {
	if idx == 0 {
		return true
	}
	for _, v := range rules[instruction[idx]] {
		if v == instruction[idx-1] {
			return false
		}
	}
	return true

}

func sortInstruction(rules map[int][]int, inst *[]int) {
	// brute force this thing
	failedIdx := -1
	for idx := range *inst {
		if !positionValid(rules, *inst, idx) {
			failedIdx = idx
			break
		}
	}
	if failedIdx != -1 {
		swap := (*inst)[failedIdx]
		(*inst)[failedIdx] = (*inst)[failedIdx-1]
		(*inst)[failedIdx-1] = swap
		sortInstruction(rules, inst)
	}
}

func calculateCheckSum(instructions [][]int) int {

	checkSum := 0
	for _, inst := range instructions {
		e := len(inst)/2 - 1
		mod := len(inst) % 2
		if mod > 0 {
			e++
		}
		checkSum += inst[e]
	}

	return checkSum
}
