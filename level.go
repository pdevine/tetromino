package main

import (
	"strconv"

	tm "github.com/gdamore/tcell/termbox"
	sprite "github.com/pdevine/go-asciisprite"
)

type LevelSelector struct {
	sprite.BaseSprite
	levels       []Level
	currentLevel uint
}

type Level struct {
	blocks []*sprite.Block
}

func (l Level) Select() {
	for _, b := range l.blocks {
		b.Fg |= tm.AttrReverse
		b.Bg |= tm.AttrReverse
	}
}

func (l Level) UnSelect() {
	for _, b := range l.blocks {
		b.Fg = 0
		b.Bg = 0
	}
}

// +-----------------------------+
// |  0  |  1  |  2  |  3  |  4  |
// +-----------------------------+

func addLine(y int, c *sprite.Costume) {
	b := &sprite.Block{X: 0, Y: y, Char: '+'}
	c.Blocks = append(c.Blocks, b)

	for cnt := 0; cnt < 29; cnt++ {
		b := &sprite.Block{X: cnt + 1, Y: y, Char: '-'}
		c.Blocks = append(c.Blocks, b)
	}

	b = &sprite.Block{X: 30, Y: y, Char: '+'}
	c.Blocks = append(c.Blocks, b)
}

func NewSelector() *LevelSelector {
	s := &LevelSelector{BaseSprite: sprite.BaseSprite{
		X:       33,
		Y:       17,
		Visible: true,
	},
	}

	for n := 0; n < 2; n++ {
		c := sprite.NewCostume("", '*')

		addLine(0, &c)

		var l Level
		// draw nums
		for cnt := 0; cnt <= 9; cnt++ {
			if n == 0 {
				var b *sprite.Block
				if cnt < 5 {
					b = &sprite.Block{X: cnt*6 + 3, Y: 1, Char: rune(strconv.Itoa(cnt)[0])}
				} else {
					b = &sprite.Block{X: (cnt-5)*6 + 3, Y: 3, Char: rune(strconv.Itoa(cnt)[0])}
				}
				l = Level{blocks: []*sprite.Block{b}}
				c.Blocks = append(c.Blocks, b)
			} else {
				var b0 *sprite.Block
				var b1 *sprite.Block
				if cnt < 5 {
					b0 = &sprite.Block{X: cnt*6 + 3, Y: 1, Char: '1'}
					b1 = &sprite.Block{X: cnt*6 + 4, Y: 1, Char: rune(strconv.Itoa(cnt)[0])}
				} else {
					b0 = &sprite.Block{X: (cnt-5)*6 + 3, Y: 3, Char: '1'}
					b1 = &sprite.Block{X: (cnt-5)*6 + 4, Y: 3, Char: rune(strconv.Itoa(cnt)[0])}
				}
				l = Level{blocks: []*sprite.Block{b0, b1}}
				c.Blocks = append(c.Blocks, b0)
				c.Blocks = append(c.Blocks, b1)
			}
			s.levels = append(s.levels, l)
		}

		s.SetLevel(0)

		// draw bars
		for cnt := 0; cnt <= 5; cnt++ {
			b := &sprite.Block{X: cnt * 6, Y: 1, Char: '|'}
			c.Blocks = append(c.Blocks, b)
			b = &sprite.Block{X: cnt * 6, Y: 3, Char: '|'}
			c.Blocks = append(c.Blocks, b)
		}

		addLine(2, &c)
		addLine(4, &c)

		s.AddCostume(c)
	}

	return s
}

func (s *LevelSelector) SetLevel(l uint) {
	if l >= 0 && l < uint(len(s.levels)) {
		s.levels[s.currentLevel].UnSelect()
		s.currentLevel = l
		s.levels[s.currentLevel].Select()
	} else {
		// XXX - warn?
	}
}

func (s *LevelSelector) MoveRight() {
	if s.currentLevel == 9 {
		s.SetLevel(0)
	} else if s.currentLevel == 19 {
		s.SetLevel(10)
	} else {
		s.SetLevel(s.currentLevel + 1)
	}
}

func (s *LevelSelector) MoveLeft() {
	if s.currentLevel == 0 {
		s.SetLevel(9)
	} else if s.currentLevel == 10 {
		s.SetLevel(19)
	} else {
		s.SetLevel(s.currentLevel - 1)
	}
}

func (s *LevelSelector) MoveUp() {
	if (s.currentLevel >= 5 && s.currentLevel < 10) || s.currentLevel >= 15 {
		s.SetLevel(s.currentLevel - 5)
	}
}

func (s *LevelSelector) MoveDown() {
	if s.currentLevel < 5 || (s.currentLevel >= 10 && s.currentLevel < 15) {
		s.SetLevel(s.currentLevel + 5)
	}
}

func (s *LevelSelector) GetVal() int {
	return int(s.currentLevel)
}
