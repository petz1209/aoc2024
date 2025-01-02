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
const MIN_DIFF = 1
const MAX_DIFF = 3
const MAX_ISSUES = 1

func main() {

	/*
		- first read the input file as a list of reports
		- check for every report if it adheres to the business logic that determins wether it is save or not
		<addon>
		the new feature states that reports that have a single problematic value, are sill save.
		- extract the unsave reports
		- for each unsave report, iterate over the items with the exeption of one. do it len(report) times.
		- if any combination fits, return true.
		- else: return fasle
	*/

	reports := readInput()
	save, unsave := DivideReportsToSaveAndUnsave(reports, MIN_DIFF, MAX_DIFF)
	fmt.Println("-----------------------------------------------")
	fmt.Printf("Number of Save Reports: %d\n", len(save))
	fmt.Println("-----------------------------------------------")

	saveReportCount := len(save)

	unsaveReports := make([][]int, len(unsave))
	for idx, v := range unsave {
		unsaveReports[idx] = reports[v]
	}

	for idx := range unsaveReports {
		if SaveWithExceptions(unsaveReports[idx], MIN_DIFF, MAX_DIFF) {
			saveReportCount++
		}
	}
	fmt.Println("-----------------------------------------------")
	fmt.Printf("Number of Save Reports: %d\n", saveReportCount)
	fmt.Println("-----------------------------------------------")

}

func readInput() [][]int {
	f, err := os.Open(INPUT)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reports := make([][]int, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), " ")
		tmp := make([]int, 0)
		for _, s := range row {
			num, err := strconv.Atoi(s)
			if err != nil {
				log.Fatal(err)
			}
			tmp = append(tmp, num)
		}
		reports = append(reports, tmp)
	}
	return reports
}

// business logic functions
func CountSaveReports(reports [][]int, mindiff, maxdiff int) int {
	var save_reports int = 0
	for idx := range reports {
		if isSave(reports[idx], mindiff, maxdiff) {
			save_reports++
			//fmt.Printf("[%d]Save Report\n", idx)
		}
	}
	return save_reports
}

func DivideReportsToSaveAndUnsave(reports [][]int, mindiff, maxdiff int) ([]int, []int) {

	save := make([]int, 0)
	unsave := make([]int, 0)
	for idx := range reports {
		if isSave(reports[idx], mindiff, maxdiff) {
			save = append(save, idx)
		} else {
			unsave = append(unsave, idx)
		}
	}
	return save, unsave

}

func SaveWithExceptions(list []int, mindiff, maxdiff int) bool {

	buffer := make([]int, len(list)-1)

	for round := 0; round < len(list); round++ {

		cursor := 0
		for idx := range list {
			if idx == round {
				continue
			}
			buffer[cursor] = list[idx]
			cursor++
		}
		if isSave(buffer, mindiff, maxdiff) {
			return true
		}

	}
	return false

}

func discoverTrend(list []int) int {
	var (
		UPWARDS   = 1
		DOWNWARDS = -1
		cntUp     = 0
		cntDown   = 0
	)
	for idx := range list {
		if idx == 0 {
			continue
		}
		if list[idx] > list[idx-1] {
			cntUp++
		} else if list[idx] < list[idx-1] {
			cntDown++
		}
	}
	if cntUp > cntDown {
		return UPWARDS
	} else if cntUp < cntDown {
		return DOWNWARDS
	}
	return 0

}

// check for a specific list of ints if it adhers to business logic
func isSave(list []int, mindiff, maxdiff int) bool {
	var (
		UPWARDS     = 1
		DOWNWARDS   = -1
		brokenIndex = -1
		last        = -1
	)

	trend := discoverTrend(list)
	for idx, curr := range list {
		if idx == 0 {
			continue
		}

		last = list[idx-1]
		if idx-1 == brokenIndex {
			last = list[idx-2]
		}

		switch trend {
		case UPWARDS:
			if curr < last {
				return false
			}
			if curr-last < mindiff || curr-last > maxdiff {
				return false
			}
		case DOWNWARDS:
			if curr > last {
				return false
			}
			if last-curr < mindiff || last-curr > maxdiff {
				return false
			}
		default:
			return false
		}
	}
	return true

}
