package main

import (
	"fmt"
	"math/rand"
	"sort"

	tm "github.com/nsf/termbox-go"
	sprite "github.com/pdevine/go-asciisprite"
)

const l0 = `xxxx[]
[][][]`

const l1 = `xx[]
xx[]
xx[][]`

const l2 = `
[][][]
[]`

const l3 = `[][]
xx[]
xx[]`

const j0 = `[]
[][][]`

const j1 = `xx[][]
xx[]
xx[]`

const j2 = `
[][][]
xxxx[]`

const j3 = `xx[]
xx[]
[][]`

const t0 = `
[][][]
xx[]`

const t1 = `xx[]
[][]
xx[]`

const t2 = `xx[]
[][][]`

const t3 = `xx[]
xx[][]
xx[]`

const i0 = `

[][][][]`

const i1 = `xxxx[]
xxxx[]
xxxx[]
xxxx[]`

const s0 = `
xx[][]
[][]`

const s1 = `xx[]
xx[][]
xxxx[]`

const z0 = `
[][]
xx[][]`

const z1 = `xxxx[]
xx[][]
xx[]`

const sq0 = `
xx[][]
xx[][]`

// Frames Per "Gridcell"
var levelFPG = map[int]int{
	0:  48,
	1:  43,
	2:  38,
	3:  33,
	4:  28,
	5:  23,
	6:  18,
	7:  13,
	8:  8,
	9:  6,
	10: 5,
	11: 5,
	12: 5,
	13: 4,
	14: 4,
	15: 4,
	16: 3,
	17: 3,
	18: 3,
	19: 2,
	20: 2,
	21: 2,
	22: 2,
	23: 2,
	24: 2,
	25: 2,
	26: 2,
	27: 2,
	28: 2,
	29: 1,
	30: 1,
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

type TetrominoBlock struct {
	sprite.BaseSprite
	ClearDirection int
	Xoffset        int
	Yoffset        int
}

func NewTetrominoBlock(x, y int) *TetrominoBlock {
	c := sprite.NewCostume("[]", 'x')

	b := &TetrominoBlock{BaseSprite: sprite.BaseSprite{
		X:              x,
		Y:              y,
		Height:         1,
		Width:          2,
		Costumes:       []sprite.Costume{c},
		CurrentCostume: 0,
		Visible:        true,
		Dead:           false,
	},
		Xoffset: x,
		Yoffset: y,
	}
	return b
}

func (s *TetrominoBlock) Update() {
	if s.ClearDirection != 0 {
		if s.ClearDirection < 0 {
			s.X--
			if s.X-20 < 0 {
				s.Visible = false
				s.Dead = true
			}
		} else {
			s.X++
			if s.X-20 > 20 {
				s.Visible = false
				s.Dead = true
			}
		}
	}
}

type Tetromino struct {
	sprite.BaseSprite
	TimeOut int
	Timer   int
	Xoffset int
	Yoffset int
	Stopped bool
	Type    string
}

func NewTetromino(ttype string) *Tetromino {
	t := &Tetromino{BaseSprite: sprite.BaseSprite{
		Alpha:   'x',
		Visible: true,
	},
		Xoffset: 20 + 2,
		Yoffset: 10,
		Timer:   0,
		TimeOut: levelFPG[0],
		Stopped: false,
		Type:    ttype,
	}

	switch {
	case ttype == "t":
		t.AddCostume(sprite.NewCostume(t0, 'x'))
		t.AddCostume(sprite.NewCostume(t1, 'x'))
		t.AddCostume(sprite.NewCostume(t2, 'x'))
		t.AddCostume(sprite.NewCostume(t3, 'x'))
	case ttype == "j":
		t.AddCostume(sprite.NewCostume(j0, 'x'))
		t.AddCostume(sprite.NewCostume(j1, 'x'))
		t.AddCostume(sprite.NewCostume(j2, 'x'))
		t.AddCostume(sprite.NewCostume(j3, 'x'))
	case ttype == "z":
		t.AddCostume(sprite.NewCostume(z0, 'x'))
		t.AddCostume(sprite.NewCostume(z1, 'x'))
	case ttype == "sq":
		t.AddCostume(sprite.NewCostume(sq0, 'x'))
	case ttype == "s":
		t.AddCostume(sprite.NewCostume(s0, 'x'))
		t.AddCostume(sprite.NewCostume(s1, 'x'))
	case ttype == "l":
		t.AddCostume(sprite.NewCostume(l0, 'x'))
		t.AddCostume(sprite.NewCostume(l1, 'x'))
		t.AddCostume(sprite.NewCostume(l2, 'x'))
		t.AddCostume(sprite.NewCostume(l3, 'x'))
	case ttype == "i":
		t.AddCostume(sprite.NewCostume(i0, 'x'))
		t.AddCostume(sprite.NewCostume(i1, 'x'))
	}

	return t
}

func getRandTetromino(src rand.Source, bg *Well) *Tetromino {
	r := rand.New(src)
	ttypes := []string{"t", "j", "z", "sq", "s", "l", "i"}

	i := r.Intn(len(ttypes))

	t := NewTetromino(ttypes[i])

	t.X = 3
	t.Y = 10
	t.Stopped = true

	return t
}

func (s *Tetromino) BottomEdgeHeight() int {
	return s.Y + s.Costumes[s.CurrentCostume].BottomEdge() - s.Yoffset
}

func (s *Tetromino) PlaceInWell() {
	stats[s.Type].IncVal()
	s.Y = background.Y - s.Costumes[s.CurrentCostume].TopEdge()
	s.X = background.X + 10

	for _, b := range allBlocks {
		for _, c := range s.Costumes[s.CurrentCostume].Blocks {
			if c.X+s.X == b.X && c.Y+s.Y == b.Y {
				gamemode = gameover
				background.TimeOut = 20
				return
			}
		}
	}

	s.Stopped = false
	s.TimeOut = levelFPG[CurrentLevel.Val]
}

func (s *Tetromino) Update() {
	if !s.Stopped {
		s.Timer = s.Timer + 1
		if s.Timer >= s.TimeOut {
			findBottomEdge(s)
			if !s.Stopped {
				s.Y = s.Y + 1
				s.Timer = 0
			}
		}
	}
}

func (s *Tetromino) RotateClockwise() {
	s.CurrentCostume = s.CurrentCostume + 1
	if s.CurrentCostume >= len(s.Costumes) {
		s.CurrentCostume = 0
	}
	/*
		if s.X < s.Xoffset+2 {
			if s.X-s.Xoffset < -findLeftEdge(s) {
				s.X = s.X + 2
			}
		}
	*/
}

func (s *Tetromino) RotateAnticlockwise() {
	s.CurrentCostume = s.CurrentCostume - 1
	if s.CurrentCostume < 0 {
		s.CurrentCostume = len(s.Costumes) - 1
	}
}

func findLeftEdge(s *Tetromino) {
	c := s.Costumes[s.CurrentCostume]

	furthestLeft := c.LeftEdge()
	m := c.LeftEdgeByRow()

	// if we're against the edge, return
	if s.X+furthestLeft-s.Xoffset <= 0 {
		return
	}

	for y, x := range m {
		for _, b := range allBlocks {
			if y+s.Y == b.Y {
				if b.X+b.Width == s.X+x {
					return
				}
			}
		}
	}
	s.X = s.X - 2
}

func findRightEdge(s *Tetromino) {
	c := s.Costumes[s.CurrentCostume]

	furthestRight := c.RightEdge()
	m := c.RightEdgeByRow()

	// if we're against the edge, return
	if s.X+furthestRight-s.Xoffset+1 >= 20 {
		return
	}

	for y, x := range m {
		for _, b := range allBlocks {
			if y+s.Y == b.Y {
				if b.X == s.X+x+1 {
					return
				}
			}
		}
	}
	s.X = s.X + 2
}

func findBottomEdge(s *Tetromino) {
	c := s.Costumes[s.CurrentCostume]
	m := c.BottomEdgeByColumn()

	for k, v := range m {
		if s.Stopped {
			break
		}

		// if we're at the bottom, stop
		if s.Y+1+v-s.Yoffset > 20 {
			s.Stopped = true
			s.convertToBlocks()
			break
		}
		// check if we're touching a block
		for _, bs := range allBlocks {
			if s.X+k == bs.X && s.Y+1+v == bs.Y {
				s.Stopped = true
				s.convertToBlocks()
				break
			}
		}
	}
}

func (s *Tetromino) convertToBlocks() {
	// XXX - this is ugly, but saves us having two data representations
	//       we should probably just omit any odd X values
	var cnt int
	for _, b := range s.Costumes[s.CurrentCostume].Blocks {
		if b.Char == '[' {
			cnt++
			nb := NewTetrominoBlock(s.X+b.X, s.Y+b.Y)
			allBlocks = append(allBlocks, nb)
			allSprites.Sprites = append(allSprites.Sprites, nb)
		}
	}
	// remove the piece
	for cnt, cs := range allSprites.Sprites {
		if s == cs {
			allSprites.Sprites = append(allSprites.Sprites[:cnt], allSprites.Sprites[cnt+1:]...)
			break
		}
	}
}

func (s *Tetromino) MoveLeft() {
	findLeftEdge(s)
}

func (s *Tetromino) MoveRight() {
	findRightEdge(s)
}

func CheckRows() {
	var rowsCleared uint8

	rows := make(map[int][]*TetrominoBlock)

	for _, b := range allBlocks {
		rows[b.Y] = append(rows[b.Y], b)
	}

	// Clear any lines
	for _, row := range rows {
		if len(row) == 10 {
			rowsCleared++
			TotalLines++
			for _, b := range row {
				if b.X-20 > 10 {
					b.ClearDirection = 1
				} else {
					b.ClearDirection = -1
				}
			}
		}
	}
	if TotalLines >= (CurrentLevel.Val+1)*10 {
		CurrentLevel.IncVal()
	}

	switch {
	case rowsCleared == 1:
		scoreText.AddVal(40 * (CurrentLevel.Val + 1))
	case rowsCleared == 2:
		scoreText.AddVal(100 * (CurrentLevel.Val + 1))
	case rowsCleared == 3:
		scoreText.AddVal(300 * (CurrentLevel.Val + 1))
	case rowsCleared == 4:
		scoreText.AddVal(1200 * (CurrentLevel.Val + 1))
	}

	scoreText.AddVal(nextScore)
}

func Vaccuum() {
	rows := make(map[int][]*TetrominoBlock)
	rowKeys := []int{}

	for _, b := range allBlocks {
		rows[b.Y] = append(rows[b.Y], b)
	}

	for y, _ := range rows {
		rowKeys = append(rowKeys, y)
	}
	sort.Sort(sort.IntSlice(rowKeys))

	for _, y := range rowKeys {
		dead := true
		for _, blk := range rows[y] {
			if !blk.Dead {
				dead = false
			}
		}
		if dead {
			for _, blk := range rows[y] {
				removeBlock(blk)
			}
			for _, b := range allBlocks {
				if y > b.Y {
					b.Y++
				}
			}
		}
	}
}

func printRows(rows map[int][]int) {
	tm.SetCursor(0, 0)
	for y, xrows := range rows {
		fmt.Printf("row %d: %d - ", y, len(xrows))
		/*
			fmt.Printf("row %d: ", y)
			for _, x := range xrows {
				fmt.Printf("(%d,%d) ", x, y)
			}
		*/
	}
	fmt.Printf("******")
}

func getBlockIndex(b *TetrominoBlock) (bool, int) {
	for cnt, blk := range allBlocks {
		if blk == b {
			return true, cnt
		}
	}
	return false, -1
}

func removeBlock(b *TetrominoBlock) bool {
	ok, idx := getBlockIndex(b)
	if ok {
		allSprites.Remove(b)
		copy(allBlocks[idx:], allBlocks[idx+1:])
		allBlocks[len(allBlocks)-1] = nil
		allBlocks = allBlocks[:len(allBlocks)-1]
	}
	return ok
}
