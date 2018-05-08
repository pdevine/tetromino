package main

import (
	"fmt"

	sprite "github.com/pdevine/go-asciisprite"
)

var areToHeight = map[int]int{
	20: 10,
	19: 10,
	18: 12,
	17: 12,
	16: 12,
	15: 12,
	14: 14,
	13: 14,
	12: 14,
	11: 14,
	10: 16,
	9:  16,
	8:  16,
	7:  16,
	6:  18,
	5:  18,
	4:  18,
	3:  18,
	2:  18,
	1:  18,
}

type LinesText struct {
	sprite.BaseSprite
	Val int
}

func NewLinesText(linesNum int) *LinesText {
	l := fmt.Sprintf("LINES-%03d", linesNum)
	c := sprite.NewCostume(l, '!')

	s := &LinesText{BaseSprite: sprite.BaseSprite{
		Costumes: []sprite.Costume{c},
		Width:    c.Width,
		Height:   c.Height,
		Visible:  true,
	},
	}
	return s
}

func (s *LinesText) UpdateLines(linesNum int) {
	l := fmt.Sprintf("LINES-%03d", linesNum)
	s.Costumes[s.CurrentCostume].ChangeCostume(l, '!')
}

type LevelText struct {
	sprite.BaseSprite
	Val int
}

func NewLevelText(level int) *LevelText {
	c := sprite.NewCostume(fmt.Sprintf("LEVEL\n  %01d", level), '!')

	s := &LevelText{BaseSprite: sprite.BaseSprite{
		X:       47,
		Y:       20,
		Visible: true,
	},
		Val: level,
	}
	s.AddCostume(c)
	return s
}

func (s *LevelText) IncVal() {
	s.Val++
	s.Costumes[s.CurrentCostume].ChangeCostume(fmt.Sprintf("LEVEL\n  %01d", s.Val), '!')
}

type Well struct {
	sprite.BaseBackground
	Timer   int
	TimeOut int
}

func NewWell() *Well {
	bg := &Well{BaseBackground: sprite.BaseBackground{
		X:     20,
		Y:     10,
		Tiled: false,
	},
	}

	block_size := 2

	for cnt := 0; cnt < 21; cnt++ {
		b := sprite.Block{
			X:    0,
			Y:    cnt,
			Char: '<',
		}
		bg.Background = append(bg.Background, b)

		b = sprite.Block{
			X:    1,
			Y:    cnt,
			Char: '!',
		}
		bg.Background = append(bg.Background, b)

		b = sprite.Block{
			X:    10*block_size + 2,
			Y:    cnt,
			Char: '!',
		}
		bg.Background = append(bg.Background, b)

		b = sprite.Block{
			X:    10*block_size + 3,
			Y:    cnt,
			Char: '>',
		}
		bg.Background = append(bg.Background, b)
	}

	for cnt := 0; cnt < 10*block_size; cnt++ {
		b := sprite.Block{
			X:    2 + cnt,
			Y:    21,
			Char: '*',
		}
		bg.Background = append(bg.Background, b)
	}
	return bg
}

func (s *Well) Update() {
	if activeTetromino.Stopped == true {
		CheckRows()
		linesText.UpdateLines(TotalLines)
		activeTetromino = nextTetromino
		activeTetromino.PlaceInWell()
		nextTetromino = getRandTetromino(src, background)
		nextTetromino.X = 45
		allSprites.Sprites = append(allSprites.Sprites, nextTetromino)
	}
	Vaccuum()
}
