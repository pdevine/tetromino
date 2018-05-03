package main

import (
	"fmt"

	sprite "github.com/pdevine/go-asciisprite"
)

type LinesText struct {
	sprite.BaseSprite
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

func NewWell() sprite.BaseBackground {
	bg := sprite.BaseBackground{
		X:     20,
		Y:     10,
		Tiled: false,
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
