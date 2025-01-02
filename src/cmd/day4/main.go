package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Field struct {
	XMax   int
	YMax   int
	XMin   int
	YMin   int
	Matrix [][]byte
}
type Point struct {
	X     int
	Y     int
	Field *Field
}

func (p *Point) Up() error {
	if p.Y == p.Field.YMin {
		return fmt.Errorf("cannot move up. already on top of the field")
	}
	p.Y--
	return nil
}
func (p *Point) Down() error {
	if p.Y == p.Field.YMax {
		return fmt.Errorf("cannot move down. already on bottom of the field")
	}
	p.Y++
	return nil
}

func (p *Point) Left() error {
	if p.X == p.Field.XMin {
		return fmt.Errorf("cannot move left. already on left boarder of the field")
	}
	p.X--
	return nil
}

func (p *Point) Right() error {
	if p.X == p.Field.XMax {
		return fmt.Errorf("cannot move right. already on right boarder of the field")
	}
	p.X++
	return nil
}

func (p *Point) GetValue() byte {
	return p.Field.Matrix[p.Y][p.X]
}

func (p *Point) Copy() Point {
	newPoint := Point{X: p.X, Y: p.Y, Field: p.Field}
	return newPoint
}

type XFinder struct {
	Center    *Point
	StartChar byte
	EndChar   byte
}

func (m *XFinder) findStartPoint() int {

	type Points struct {
		TopLeft  Point
		BotRight Point
		TopRight Point
		BotLeft  Point
	}
	if m.Center.Y == m.Center.Field.YMax || m.Center.Y == m.Center.Field.YMin ||
		m.Center.X == m.Center.Field.XMax || m.Center.X == m.Center.Field.XMin {
		return 0
	}
	tr := m.Center.Copy()
	tr.Right()
	tr.Up()
	tl := m.Center.Copy()
	tl.Left()
	tl.Up()
	br := m.Center.Copy()
	br.Right()
	br.Down()
	bl := m.Center.Copy()
	bl.Left()
	bl.Down()
	points := Points{TopLeft: tl, TopRight: tr, BotLeft: bl, BotRight: br}

	switch {
	case points.TopLeft.GetValue() == m.StartChar:
		if points.BotRight.GetValue() != m.EndChar {
			return 0
		}
	case points.TopLeft.GetValue() == m.EndChar:
		if points.BotRight.GetValue() != m.StartChar {
			return 0
		}
	default:
		return 0
	}
	switch {
	case points.TopRight.GetValue() == m.StartChar:
		if points.BotLeft.GetValue() != m.EndChar {
			return 0
		}
	case points.TopRight.GetValue() == m.EndChar:
		if points.BotLeft.GetValue() != m.StartChar {
			return 0
		}
	default:
		return 0
	}

	return 1
}

const INPUT = "input.txt"

func main() {
	matrix := ReadInput2()
	xmasCnt := CountXmas2(matrix)
	fmt.Printf("Rows: %d, Columns: %d\n", len(matrix), len(matrix[0]))
	fmt.Printf("XmasCount: %d\n", xmasCnt)
}

func CountXmas(matrix [][]byte) int {

	var xmasCnt int = 0
	var MyField Field = Field{}
	var lastChar byte = 'S'
	MyField.YMax = len(matrix) - 1
	MyField.XMax = len(matrix[0]) - 1
	MyField.Matrix = matrix
	startChords := collectStartPoints(matrix, 'X')

	startPoints := make([]Point, 0)

	for i := range startChords {
		p := Point{X: startChords[i][1],
			Y:     startChords[i][0],
			Field: &MyField,
		}
		startPoints = append(startPoints, p)
	}
	fmt.Printf("Total StartPoints: %d\n", len(startPoints))

	for _, p := range startPoints {
		xmasCnt += Left2Right(p, lastChar)
		xmasCnt += Right2Left(p, lastChar)
		xmasCnt += Top2Down(p, lastChar)
		xmasCnt += Bot2Up(p, lastChar)
		xmasCnt += LT2RB(p, lastChar)
		xmasCnt += LB2RT(p, lastChar)
		xmasCnt += RT2LB(p, lastChar)
		xmasCnt += RB2LT(p, lastChar)

	}
	return xmasCnt

}

func CountXmas2(matrix [][]byte) int {
	var xmasCnt int = 0
	var MyField Field = Field{}
	MyField.YMax = len(matrix) - 1
	MyField.XMax = len(matrix[0]) - 1
	MyField.Matrix = matrix
	startChords := collectStartPoints(matrix, 'A')
	startPoints := make([]Point, 0)

	for i := range startChords {
		p := Point{X: startChords[i][1],
			Y:     startChords[i][0],
			Field: &MyField,
		}
		startPoints = append(startPoints, p)
	}

	for _, p := range startPoints {
		finder := XFinder{Center: &p, StartChar: 'M', EndChar: 'S'}
		xmasCnt += finder.findStartPoint()
	}

	return xmasCnt

}

func ReadInput() [][]byte {

	f, err := os.Open(INPUT)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	matrix := make([][]byte, 0)

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		matrix = append(matrix, sc.Bytes())
	}
	return matrix
}

func ReadInput2() [][]byte {
	f, err := os.Open(INPUT)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	matrix := make([][]byte, 0)

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		matrix = append(matrix, []byte(sc.Text()))
	}
	return matrix
}

func collectStartPoints(matrix [][]byte, startChar byte) [][]int {
	startPoints := make([][]int, 0)

	for y := range matrix {
		for x := range matrix[y] {
			//fmt.Printf("Letter: %s, [%d, %d]\n", string(matrix[y][x]), y, x)
			if matrix[y][x] == startChar {
				startPoints = append(startPoints, []int{y, x})
			}
		}

	}
	return startPoints

}

func NextValueMatch(cchar byte) byte {
	var Map map[byte]byte = map[byte]byte{'X': 'M', 'M': 'A', 'A': 'S'}
	return Map[cchar]

}

func Left2Right(p Point, lastChar byte) int {
	if p.GetValue() == lastChar {
		return 1
	}
	nV := NextValueMatch(p.GetValue())
	// check if we can even move to the right
	if err := p.Right(); err != nil {
		return 0
	}
	if p.GetValue() != nV {
		return 0
	}
	return Left2Right(p, lastChar)
}

func Right2Left(p Point, lastChar byte) int {
	if p.GetValue() == lastChar {
		return 1
	}
	nV := NextValueMatch(p.GetValue())
	// check if we can even move to the right
	if err := p.Left(); err != nil {
		return 0
	}
	if p.GetValue() != nV {
		return 0
	}
	return Right2Left(p, lastChar)
}

func Top2Down(p Point, lastChar byte) int {
	if p.GetValue() == lastChar {
		return 1
	}
	nV := NextValueMatch(p.GetValue())
	// check if we can even move to the right
	if err := p.Down(); err != nil {
		return 0
	}
	if p.GetValue() != nV {
		return 0
	}
	return Top2Down(p, lastChar)
}

func Bot2Up(p Point, lastChar byte) int {
	if p.GetValue() == lastChar {
		return 1
	}
	nV := NextValueMatch(p.GetValue())
	// check if we can even move to the right
	if err := p.Up(); err != nil {
		return 0
	}
	if p.GetValue() != nV {
		return 0
	}
	return Bot2Up(p, lastChar)
}

// leftTop 2 RightBot
func LT2RB(p Point, lastChar byte) int {
	if p.GetValue() == lastChar {
		return 1
	}
	nV := NextValueMatch(p.GetValue())
	// check if we can even move to the right
	if err := p.Down(); err != nil {
		return 0
	}
	if err := p.Right(); err != nil {
		return 0
	}
	if p.GetValue() != nV {
		return 0
	}
	return LT2RB(p, lastChar)
}

func LB2RT(p Point, lastChar byte) int {
	if p.GetValue() == lastChar {
		return 1
	}
	nV := NextValueMatch(p.GetValue())
	// check if we can even move to the right
	if err := p.Up(); err != nil {
		return 0
	}
	if err := p.Right(); err != nil {
		return 0
	}
	if p.GetValue() != nV {
		return 0
	}
	return LB2RT(p, lastChar)
}

func RT2LB(p Point, lastChar byte) int {
	if p.GetValue() == lastChar {
		return 1
	}
	nV := NextValueMatch(p.GetValue())
	// check if we can even move to the right
	if err := p.Down(); err != nil {
		return 0
	}
	if err := p.Left(); err != nil {
		return 0
	}
	if p.GetValue() != nV {
		return 0
	}
	return RT2LB(p, lastChar)
}

func RB2LT(p Point, lastChar byte) int {
	if p.GetValue() == lastChar {
		return 1
	}
	nV := NextValueMatch(p.GetValue())
	// check if we can even move to the right
	if err := p.Up(); err != nil {
		return 0
	}
	if err := p.Left(); err != nil {
		return 0
	}
	if p.GetValue() != nV {
		return 0
	}
	return RB2LT(p, lastChar)
}
