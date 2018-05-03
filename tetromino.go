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
	0: 48,
	1: 43,
	2: 38,
	3: 33,
	4: 28,
	5: 23,
	6: 18,
	7: 13,
	8: 8,
	9: 6,
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
	Xoffset int
	Yoffset int
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
	},
		Xoffset: x,
		Yoffset: y,
	}
	return b
}

type Tetromino struct {
	sprite.BaseSprite
	TimeOut int
	Timer   int
	Xoffset int
	Yoffset int
	Stopped bool
}

func NewTetromino() *Tetromino {
	t := &Tetromino{BaseSprite: sprite.BaseSprite{
		Alpha:   'x',
		Visible: true,
	},
		Xoffset: 20 + 2,
		Yoffset: 10,
		Timer:   0,
		TimeOut: levelFPG[0],
		Stopped: false,
	}
	return t
}

func getRandTetromino(src rand.Source, bg sprite.BaseBackground) *Tetromino {
	r := rand.New(src)
	i := r.Intn(7)

	var t *Tetromino

	switch {
	case i == 0:
		t = NewL()
	case i == 1:
		t = NewJ()
	case i == 2:
		t = NewT()
	case i == 3:
		t = NewS()
	case i == 4:
		t = NewZ()
	case i == 5:
		t = NewI()
	case i == 6:
		t = NewSq()
	}

	t.X = 3
	t.Y = 10
	t.Stopped = true

	return t
}

func (s *Tetromino) PlaceInWell() {
	s.Y = background.Y - s.Costumes[s.CurrentCostume].TopEdge()
	s.X = background.X + 10
	s.Stopped = false
	s.TimeOut = levelFPG[CurrentLevel]
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

func CheckRows() {
	// 1. get all of the blocks in row
	// 2. check if there are no gaps
	// 3. if the row is complete:
	// 4.   delete the row
	tm.SetCursor(0, 0)

	rows := make(map[int][]int)
	var rowKeys []int

	for _, b := range allBlocks {
		rows[b.Y] = append(rows[b.Y], b.X)
	}

	for r := range rows {
		rowKeys = append(rowKeys, r)
	}

	// remove blocks from the top first
	sort.Sort(sort.IntSlice(rowKeys))

	for _, y := range rowKeys {
		if len(rows[y]) == 10 {
			TotalLines++
			// remove blocks
			for _, x := range rows[y] {
				removeBlock(x, y)
			}

			// move blocks by one row
			for _, b := range allBlocks {
				if b.Y < y {
					b.Y = b.Y + 1
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

func getBlockIndex(x, y int) (bool, int) {
	for cnt, blk := range allBlocks {
		if blk.X == x && blk.Y == y {
			return true, cnt
		}
	}
	return false, -1
}

func removeBlock(x, y int) {
	ok, idx := getBlockIndex(x, y)
	blk := allBlocks[idx]
	if ok {
		allSprites.Remove(blk)
		copy(allBlocks[idx:], allBlocks[idx+1:])
		allBlocks[len(allBlocks)-1] = nil
		allBlocks = allBlocks[:len(allBlocks)-1]
	} else {
		fmt.Printf("couldn't remove block at %d, %d", x, y)
	}
}

func (s *Tetromino) MoveRight() {
	findRightEdge(s)
}

func NewL() *Tetromino {
	l := NewTetromino()
	l.AddCostume(sprite.NewCostume(l0, 'x'))
	l.AddCostume(sprite.NewCostume(l1, 'x'))
	l.AddCostume(sprite.NewCostume(l2, 'x'))
	l.AddCostume(sprite.NewCostume(l3, 'x'))
	return l
}

func NewJ() *Tetromino {
	j := NewTetromino()
	j.AddCostume(sprite.NewCostume(j0, 'x'))
	j.AddCostume(sprite.NewCostume(j1, 'x'))
	j.AddCostume(sprite.NewCostume(j2, 'x'))
	j.AddCostume(sprite.NewCostume(j3, 'x'))
	return j
}

func NewT() *Tetromino {
	t := NewTetromino()
	t.AddCostume(sprite.NewCostume(t0, 'x'))
	t.AddCostume(sprite.NewCostume(t1, 'x'))
	t.AddCostume(sprite.NewCostume(t2, 'x'))
	t.AddCostume(sprite.NewCostume(t3, 'x'))
	return t
}

func NewS() *Tetromino {
	s := NewTetromino()
	s.AddCostume(sprite.NewCostume(s0, 'x'))
	s.AddCostume(sprite.NewCostume(s1, 'x'))
	return s
}

func NewZ() *Tetromino {
	z := NewTetromino()
	z.AddCostume(sprite.NewCostume(z0, 'x'))
	z.AddCostume(sprite.NewCostume(z1, 'x'))
	return z
}

func NewI() *Tetromino {
	i := NewTetromino()
	i.AddCostume(sprite.NewCostume(i0, 'x'))
	i.AddCostume(sprite.NewCostume(i1, 'x'))
	return i
}

func NewSq() *Tetromino {
	sq := NewTetromino()
	sq.AddCostume(sprite.NewCostume(sq0, 'x'))
	return sq
}
