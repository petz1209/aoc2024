package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

const RENDER = false
const RENDER_SPEED = 200
const RENDER_INFO_SPEED = 2000
const INPUT = "input.txt"

type Point struct {
	x int
	y int
}

type Vector struct {
	From Point
	To   Point
}

type Officer struct {
	Point
	Moves       int
	Direction   string
	field       Field
	OnField     bool
	Positions   []Point
	Vectors     []Vector
	MovesInLoop bool
	OnPathToOut bool
}

func NewOfficer(position Point, direction string, field Field) Officer {
	return Officer{Point: position, Direction: direction,
		OnField:   true,
		field:     field,
		Positions: []Point{Point{x: position.x, y: position.y}},
	}
}

type Field struct {
	XMin  int
	XMax  int
	YMin  int
	YMax  int
	field [][]byte
}

func (f *Field) IsObstacle(p Point) bool {
	if f.field[p.y][p.x] == '#' {
		return true
	}
	return false

}
func NewField(matrix [][]byte) Field {
	return Field{XMin: 0, XMax: len(matrix[0]) - 1,
		YMin: 0, YMax: len(matrix) - 1,
		field: matrix,
	}
}

func (o *Officer) turnRight() {
	vpoint := Point{x: o.Point.x, y: o.Point.y}
	// make sure that a vector is not within the same point  a situation of two turns back to back
	if vpoint.x == o.Vectors[len(o.Vectors)-1].From.x && vpoint.y == o.Vectors[len(o.Vectors)-1].From.y {
	} else {
		o.Vectors[len(o.Vectors)-1].To = vpoint
		loop := VectorEqual(o.Vectors[:len(o.Vectors)-1], o.Vectors[len(o.Vectors)-1])
		if loop {
			o.MovesInLoop = true
		}
		o.Vectors = append(o.Vectors, Vector{From: vpoint})
	}

	directions := map[string]string{
		"up":    "right",
		"right": "down",
		"down":  "left",
		"left":  "up",
	}
	//fmt.Printf("turned from %s to %s\n", o.Direction, directions[o.Direction])
	o.Direction = directions[o.Direction]
}
func (o *Officer) move() {

	o.IsOnPathToOut()
	if o.Moves == 0 {
		o.Vectors = append(o.Vectors, Vector{From: Point{x: o.Point.x, y: o.Point.y}})
	}
	o.Moves++
	var nextPoint Point
	switch o.Direction {

	case "up":
		nextPoint.x = o.Point.x
		nextPoint.y = o.Point.y - 1
	case "right":
		nextPoint.x = o.Point.x + 1
		nextPoint.y = o.Point.y
	case "down":
		nextPoint.x = o.Point.x
		nextPoint.y = o.Point.y + 1
	case "left":
		nextPoint.x = o.Point.x - 1
		nextPoint.y = o.Point.y
	}
	if o.LeftField(nextPoint) {
		return
	}
	if o.field.IsObstacle(nextPoint) {
		o.turnRight()
		return
	}
	o.Point.x = nextPoint.x
	o.Point.y = nextPoint.y
	o.AddPosition(nextPoint)
}

func (o *Officer) AddPosition(p Point) {

	for i := range o.Positions {
		if o.Positions[i].x == p.x && o.Positions[i].y == p.y {
			return
		}
	}
	o.Positions = append(o.Positions, p)

}

func (o *Officer) LeftField(p Point) bool {
	if p.x < o.field.XMin || p.x > o.field.XMax || p.y < o.field.YMin || p.y > o.field.YMax {
		o.OnField = false
		return true
	}
	return false
}

func (o *Officer) Draw() {
	snapshot := make([][]byte, len(o.field.field))
	for i := range o.field.field {
		tmp := make([]byte, len(o.field.field[i]))
		for j := range o.field.field[i] {
			if o.field.field[i][j] == '^' {
				tmp[j] = '.'
			} else {
				tmp[j] = o.field.field[i][j]
			}
		}
		snapshot[i] = tmp
	}

	cur := '^'
	switch o.Direction {
	case "right":
		cur = '>'
	case "down":
		cur = 'v'
	case "left":
		cur = '<'
	}
	snapshot[o.Point.y][o.Point.x] = byte(cur)
	drawField(snapshot)
}

func (o *Officer) IsOnPathToOut() bool {

	switch o.Direction {
	case "up":
		for i := o.y; i > 0; i-- {
			if o.field.field[i][o.x] == '#' {
				o.OnPathToOut = true
				return true
			}
		}
		return false

	case "right":
		for i := o.x; i < len(o.field.field[o.y]); i++ {
			if o.field.field[o.y][i] == '#' {
				o.OnPathToOut = true
				return true
			}
		}
		return false
	case "down":
		for i := range o.field.field {
			if i < o.y {
				continue
			}
			if o.field.field[i][o.x] == '#' {
				o.OnPathToOut = true
				return true
			}
		}
		return false

	case "left":
		for i := o.x; i > 0; i-- {
			if o.field.field[o.y][i] == '#' {
				o.OnPathToOut = true
				return true
			}
		}
		return false
	}
	return false

}

func main() {
	res := solvePuzzle1(INPUT)
	fmt.Println("---------------------------------------------")
	fmt.Printf("day6-puzzle1 result: %d\n", res)
	fmt.Println("---------------------------------------------")
	res2 := solvePuzzle2(INPUT)
	fmt.Println("---------------------------------------------")
	fmt.Printf("day6-puzzle2 result: %d\n", res2)
	fmt.Println("---------------------------------------------")
}

func solvePuzzle1(inputPath string) int {
	result := 0
	matrix := readInput(inputPath)
	field := NewField(matrix)
	//drawField(matrix)
	officerStartPosition := FindOfficerStartPosition(matrix)
	officer := NewOfficer(officerStartPosition, "up", field)
	for {
		officer.move()
		if !officer.OnField {
			break
		}
	}
	result = len(officer.Positions)
	return result

}

func solvePuzzle2(inputPath string) int {
	/* I have to place one additional # so the officer gets stuck in a forever loop
	- cant be the start position
	- should be only positions that aren't # already

	ideas:
	build vectors everytime the guy runs a straight line
	if the same vector is followed again. We know that it must be a loop
	otherwise he falls out and than it's over as well
	*/
	result := 0
	matrix := readInput(inputPath)
	officerStartPosition := FindOfficerStartPosition(matrix)
	extractPositions := func(matrix [][]byte, sp Point) []Point {
		field := NewField(matrix)
		officer := NewOfficer(sp, "up", field)
		for {
			officer.move()
			if !officer.OnField {
				break
			}
		}
		return officer.Positions
	}(matrix, officerStartPosition)

	for idx := range extractPositions {
		if extractPositions[idx].x == officerStartPosition.x && extractPositions[idx].y == officerStartPosition.y {
			continue
		}
		// Adjust the position in the matrix to #
		// create new officer
		newMatrix := readInput(inputPath)
		newMatrix[extractPositions[idx].y][extractPositions[idx].x] = '#'
		field := NewField(newMatrix)
		officer := NewOfficer(officerStartPosition, "up", field)
		//fmt.Printf("Round: %d\n", idx)
		// let him run
		// if he gets stuck in a loop great
		for {
			if RENDER {
				officer.Draw()
				time.Sleep(RENDER_SPEED * time.Millisecond)
				fmt.Print("\033[2J")
			}

			officer.move()

			if officer.OnPathToOut {
				//fmt.Println("Is on Path Out")
				//break
			}
			if !officer.OnField {
				break
			}
			if officer.MovesInLoop {
				result++
				break
			}

		}
	}
	return result

}

func FindOfficerStartPosition(matrix [][]byte) Point {
	officerStartPosition := Point{}
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] == '^' {
				officerStartPosition.x = j
				officerStartPosition.y = i
				return officerStartPosition
			}
		}
	}
	return officerStartPosition
}

func readInput(inputPath string) [][]byte {

	matrix := make([][]byte, 0)
	f, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		row := []byte(scanner.Text())
		matrix = append(matrix, row)
	}
	return matrix
}

func drawField(matrix [][]byte) {
	width := len(matrix[0])*2 + 2
	for i := 0; i < width; i++ {
		fmt.Print("-")
	}
	fmt.Println()

	for _, row := range matrix {
		fmt.Print("|")
		for _, v := range row {
			fmt.Printf("%s ", string(v))
		}
		fmt.Print("|")
		fmt.Println()
	}
	for i := 0; i < width; i++ {
		fmt.Print("-")
	}
	fmt.Println()

}

func VectorEqual(vecArr []Vector, vec Vector) bool {
	//fmt.Printf("len(vecArr): %d\n", len(vecArr))
	//fmt.Printf("vectArr: %++v\n", vecArr)
	for idx := range vecArr {
		//fmt.Printf("vec[%d]: %++v,  comp: %++v\n", idx, vecArr[idx], vec)
		if vecArr[idx].From.x == vec.From.x && vecArr[idx].From.y == vec.From.y &&
			vecArr[idx].To.x == vec.To.x && vecArr[idx].To.y == vec.To.y {
			//fmt.Printf("Same rout twice %++v\n", vec)
			return true
		}
	}
	return false
}
