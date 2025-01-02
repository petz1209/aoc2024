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

var active bool = true

func main() {

	var result int = 0
	instructions := readInput()

	v1InstructionCount := 0
	v2InstructionCount := 0
	for idx := range instructions {
		cmds2 := extractNoise(instructions[idx])
		cmds := extactNoiseCondionally(instructions[idx])
		v1InstructionCount += len(cmds2)
		v2InstructionCount += len(cmds)

		for _, cmd := range cmds {
			instrucionRes := executor(cmd)
			result += instrucionRes
		}
	}

	fmt.Println("------------------------------------------")
	fmt.Printf("Result: %d\n", result)
	fmt.Println("------------------------------------------")

	fmt.Printf("V1Count: %d  v2Count:%d\n", v1InstructionCount, v2InstructionCount)

}

func readInput() []string {

	instructions := make([]string, 0)

	f, err := os.Open(INPUT)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		instructions = append(instructions, scanner.Text())
	}
	return instructions

}

func extractNoise(row string) []string {

	store := make([]string, 0)
	for idx := range row {
		switch row[idx] {
		case 'm':
			fmt.Printf("We got m -> %s\n", row[idx:])
			res, ok := handleSMulStartEvent(row, idx)
			if ok {
				store = append(store, res)
			}
		default:
			continue
		}
	}
	return store
}

func extactNoiseCondionally(row string) []string {
	//active = true
	store := make([]string, 0)
	for idx := range row {
		switch row[idx] {
		case 'm':
			if !active {
				continue
			}
			res, ok := handleSMulStartEvent(row, idx)
			if ok {
				store = append(store, res)
			}

		case 'd':
			HandleActivationEvent(row, idx)

		default:
			continue
		}
	}
	return store

}

func handleSMulStartEvent(row string, index int) (string, bool) {
	lcursor := index
	rcursor := index + 3
	// check if we call the correct function
	if row[lcursor:rcursor+1] != "mul(" {
		return "", false
	}
	rcursor++
	lcursor = rcursor
	// now we need to extract the params section
	num1, err := findNumber(row[lcursor:], ',')
	if err != nil {
		return "", false
	}
	// move cursor
	for i := lcursor; ; i++ {
		rcursor = i
		if row[i] == ',' {
			break
		}
	}
	lcursor = rcursor
	num2, err := findNumber(row[lcursor+1:], ')')
	if err != nil {
		return "", false
	}
	return "mul(" + num1 + "," + num2 + ")", true

}

func findNumber(str string, breakSymbol byte) (string, error) {
	num := ""
	for i := 0; i < len(str); i++ {
		switch i {
		case 0:
			if !isNumb(str[i]) {
				return "", fmt.Errorf("inputs must be numeric")
			}
			num = fmt.Sprintf("%s%s", num, string(str[i]))

		case 3:
			if str[i] != breakSymbol {
				return "", fmt.Errorf("invalid number length")
			}
			return num, nil

		default:
			if !isNumb(str[i]) && str[i] != breakSymbol {
				return "", fmt.Errorf("invalid symbol")
			}
			if str[i] == breakSymbol {
				return num, nil
			}
			if isNumb(str[i]) {
				num = fmt.Sprintf("%s%s", num, string(str[i]))
			}

		}

	}
	return "", fmt.Errorf("an error Occurred")
}

func isNumb(c byte) bool {
	//var NUMBERS = []rune{'0', '1', '2', '3','4', '5', '6','7', '8','9'}
	NUMBERS := "0123456789"
	for idx := range NUMBERS {
		if c == NUMBERS[idx] {
			return true
		}
	}
	return false

}

func executor(cmd string) int {
	switch {
	case strings.HasPrefix(cmd, "mul"):
		return mul(cmd)
	default:
		return 0
	}

}

func mul(cmd string) int {
	cmd = strings.TrimPrefix(cmd, "mul(")
	cmd = strings.TrimSuffix(cmd, ")")
	params := strings.Split(cmd, ",")
	num1, _ := strconv.Atoi(params[0])
	num2, _ := strconv.Atoi(params[1])
	return num1 * num2

}

func HandleActivationEvent(str string, idx int) {
	fmt.Println("---------------------------------")
	fmt.Println(str[idx : idx+7])
	fmt.Println("---------------------------------")
	if str[idx:idx+7] == "don't()" {
		active = false
	} else if str[idx:idx+4] == "do()" {
		active = true
	}
	fmt.Println(active)
}
