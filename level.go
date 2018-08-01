package main

import (
	"strconv"

	tm "github.com/gdamore/tcell/termbox"
	sprite "github.com/pdevine/go-asciisprite"
)

type LevelSelector struct {
	sprite.BaseSprite
	levels       []*sprite.Block
	currentLevel uint
}

// +-------------------+
// | 0 | 1 | 2 | 3 | 4 |
// +-------------------+

func addLine(y int, c *sprite.Costume) {
	b := &sprite.Block{X: 0, Y: y, Char: '+'}
	c.Blocks = append(c.Blocks, b)

	for cnt := 0; cnt < 19; cnt++ {
		b := &sprite.Block{X: cnt + 1, Y: y, Char: '-'}
		c.Blocks = append(c.Blocks, b)
	}

	b = &sprite.Block{X: 20, Y: y, Char: '+'}
	c.Blocks = append(c.Blocks, b)
}

func NewSelector() *LevelSelector {

	s := &LevelSelector{BaseSprite: sprite.BaseSprite{
		X:       38,
		Y:       17,
		Visible: true,
	},
	}

	c := sprite.NewCostume("", '*')

	addLine(0, &c)

	// draw nums
	for cnt := 0; cnt <= 9; cnt++ {
		var b *sprite.Block
		if cnt < 5 {
			b = &sprite.Block{X: cnt*4 + 2, Y: 1, Char: rune(strconv.Itoa(cnt)[0])}
		} else {
			b = &sprite.Block{X: (cnt-5)*4 + 2, Y: 3, Char: rune(strconv.Itoa(cnt)[0])}
		}
		c.Blocks = append(c.Blocks, b)
		s.levels = append(s.levels, b)
	}

	s.levels[0].Fg = s.levels[0].Fg | tm.AttrReverse
	s.levels[0].Bg = s.levels[0].Fg | tm.AttrReverse

	// draw bars
	for cnt := 0; cnt <= 5; cnt++ {
		b := &sprite.Block{X: cnt * 4, Y: 1, Char: '|'}
		c.Blocks = append(c.Blocks, b)
		b = &sprite.Block{X: cnt * 4, Y: 3, Char: '|'}
		c.Blocks = append(c.Blocks, b)
	}

	addLine(2, &c)
	addLine(4, &c)

	s.AddCostume(c)

	return s
}

func (s *LevelSelector) SetLevel(l uint) {
	if l >= 0 && l < uint(len(s.levels)) {
		s.levels[s.currentLevel].Fg = 0
		s.levels[s.currentLevel].Bg = 0
		s.currentLevel = l
		s.levels[s.currentLevel].Fg |= tm.AttrReverse
		s.levels[s.currentLevel].Bg |= tm.AttrReverse
	}
}

func (s *LevelSelector) MoveRight() {
	if s.currentLevel == uint(len(s.levels)-1) {
		s.SetLevel(0)
	} else {
		s.SetLevel(s.currentLevel + 1)
	}
}

func (s *LevelSelector) MoveLeft() {
	if s.currentLevel == 0 {
		s.SetLevel(uint(len(s.levels)) - 1)
	} else {
		s.SetLevel(s.currentLevel - 1)
	}
}

func (s *LevelSelector) MoveUp() {
	if s.currentLevel >= 5 {
		s.SetLevel(s.currentLevel - 5)
	}
}

func (s *LevelSelector) MoveDown() {
	if s.currentLevel < 5 {
		s.SetLevel(s.currentLevel + 5)
	}
}

func (s *LevelSelector) GetVal() int {
	return int(s.currentLevel)
}
