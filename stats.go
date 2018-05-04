package main

import (
	"fmt"

	sprite "github.com/pdevine/go-asciisprite"
)

type StatSprite struct {
	sprite.BaseSprite
	Val int
}

var stats map[string]*StatSprite

func NewStatSprite(x, y int, val int) *StatSprite {

	s := &StatSprite{BaseSprite: sprite.BaseSprite{
		X:       x,
		Y:       y,
		Visible: true,
	},
		Val: val,
	}
	c := sprite.NewCostume(fmt.Sprintf("%03d", val), '!')
	s.AddCostume(c)

	return s
}

func (s *StatSprite) SetVal(val int) {
	s.Costumes[s.CurrentCostume].ChangeCostume(fmt.Sprintf("%03d", val), '!')
	s.Val = val
}

func (s *StatSprite) IncVal() {
	s.Val++
	s.Costumes[s.CurrentCostume].ChangeCostume(fmt.Sprintf("%03d", s.Val), '!')
}

func NewStats() {
	var t *Tetromino

	stats = make(map[string]*StatSprite)

	for cnt, shape := range []string{"t", "j", "z", "sq", "s", "l", "i"} {
		t = NewTetromino(shape)
		t.X = 5
		t.Y = 4*cnt + 8
		t.Stopped = true

		stats[shape] = NewStatSprite(13, 4*cnt+9, 0)

		allSprites.Sprites = append(allSprites.Sprites, t)
		allSprites.Sprites = append(allSprites.Sprites, stats[shape])
	}
}
