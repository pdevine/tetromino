package main

import sprite "github.com/pdevine/go-asciisprite"

type Title struct {
	sprite.BaseSprite
	state   uint
	VX      int
	VY      int
	TargetY int
}

const (
	dropin = iota
	resting
)

const logo = `[][][] [][] [][][] [][]     []   []    [] [] []   []   []
  []   []     []   []  [] []  [] [][][][] [] [][] [] []  []
  []   [][]   []   [][]   []  [] [] [] [] [] [] [][] []  []
  []   []     []   []  [] []  [] []    [] [] []   [] []  []
  []   [][]   []   []  []   []   []    [] [] []   []   []
`

const elek_costume1 = "e l e k t r o n i k a"
const elek_costume2 = "Электроника"

func NewTitle() *Title {
	s := &Title{BaseSprite: sprite.BaseSprite{
		X:       17,
		Y:       -10,
		Visible: true,
	},
		VY:      2,
		TargetY: 5,
		state:   dropin,
	}

	s.AddCostume(sprite.NewCostume(logo, '%'))

	return s
}

func (s *Title) Update() {
	switch {
	case s.state == dropin:
		if s.Y < s.TargetY {
			s.Y += s.VY
		} else {
			s.state = resting
		}
	}
}

type TitleString struct {
	sprite.BaseSprite
	Timer   int
	TimeOut int
}

func NewTitleString() *TitleString {
	s := &TitleString{BaseSprite: sprite.BaseSprite{
		X:       38,
		Y:       12,
		Visible: true,
	},
		TimeOut: 100,
	}

	s.AddCostume(sprite.NewCostume(elek_costume1, '!'))
	s.AddCostume(sprite.NewCostume(elek_costume2, '!'))

	return s
}

func (s *TitleString) Update() {
	s.Timer++
	if s.Timer >= s.TimeOut {
		s.NextCostume()
		s.Timer = 0
	}
}
