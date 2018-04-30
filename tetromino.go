package main

import (
	"fmt"
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
		Xoffset: 10 + 2,
		Yoffset: 10,
		Timer:   0,
		TimeOut: 20,
		Stopped: false,
	}
	return t
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
	m := make(map[int]int)

	tm.SetCursor(0, 0)

	furthestLeft := c.Blocks[0].X
	for _, b := range c.Blocks {
		if _, ok := m[b.Y]; ok == false {
			m[b.Y] = b.X
		}
		m[b.Y] = min(m[b.Y], b.X)
		furthestLeft = min(furthestLeft, b.X)
	}

	if s.X+furthestLeft-s.Xoffset <= 0 {
		return
	}

	blocks := make(map[int]int)
	cBlocks := make(map[int]int)
	for y, x := range m {
		for _, b := range allBlocks {
			if y+s.Y == b.Y {
				if b.X+b.Width == s.X+x {
					fmt.Printf("nope")
					return
				}
				blocks[b.Y] = b.X + b.Width
			}
		}
		cBlocks[y+s.Y] = x + s.X
	}
	//fmt.Printf("blocks: %q  cBlocks: %q  --", blocks, cBlocks)
	s.X = s.X - 2
}

func findRightEdge(s *Tetromino) int {
	c := s.Costumes[s.CurrentCostume]
	x := 0
	for _, b := range c.Blocks {
		if b.X > x {
			x = b.X
		}
		// 4 blocks max * block size
		if x == 4*2 {
			break
		}
	}
	return x
}

func findBottomEdge(s *Tetromino) {
	c := s.Costumes[s.CurrentCostume]
	m := make(map[int]int)

	// iterate through each of the blocks and find the
	// lowest one (co-ordinates are flipped)
	for _, b := range c.Blocks {
		m[b.X] = max(m[b.X], b.Y)
	}
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
	// XXX - this is ugly, but saves us having two data represenations
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
	fmt.Printf("adding %d blocks ", cnt)
	/* remove the piece */
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

	//for y, xrow := range rows {
	for _, y := range rowKeys {
		if len(rows[y]) == 10 {
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

	// debug
	morerows := make(map[int][]int)
	for _, b := range allBlocks {
		morerows[b.Y] = append(morerows[b.Y], b.X)
	}
	printRows(morerows)
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
	s.X = s.X + 2
	if s.X+findRightEdge(s) > 10*2+s.Yoffset+2 {
		s.X = s.X - 2
	}
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
