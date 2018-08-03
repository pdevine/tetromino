package main

import (
	"fmt"

	sprite "github.com/pdevine/go-asciisprite"
)

const wellLine = `[][][][][][][][][][]`

// are (time to wait) before next tetromino
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

type ScoreText struct {
	sprite.BaseSprite
	Val int
}

func NewScoreText() *ScoreText {
	c := sprite.NewCostume(fmt.Sprintf("SCORE\n%07d", 0), '!')

	s := &ScoreText{BaseSprite: sprite.BaseSprite{
		X:       47,
		Y:       4,
		Visible: true,
	},
		Val: 0,
	}
	s.AddCostume(c)
	return s
}

func (s *ScoreText) AddVal(v int) {
	s.Val += v
	s.Costumes[s.CurrentCostume].ChangeCostume(fmt.Sprintf("SCORE\n%07d", s.Val), '!')
}

type Well struct {
	sprite.BaseBackground
	Timer       int
	TimeOut     int
	FilledLines int
	Paused      bool
	PauseSprite *sprite.BaseSprite
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
	if gamemode == gameover || gamemode == cathedral {
		s.Timer++
		if gamemode == gameover {
			if s.Timer >= s.TimeOut && s.FilledLines <= 20 {
				c := sprite.NewCostume(wellLine, '~')
				l := sprite.NewBaseSprite(s.X+2, s.Y+s.FilledLines, c)
				allSprites.Sprites = append(allSprites.Sprites, l)
				s.Timer = 10
				s.FilledLines++
			} else if s.FilledLines >= 20 {
				if scoreText.Val < 30000 {
					return
				}
				for _, ts := range allSprites.Sprites {
					allSprites.Remove(ts)
				}
				gamemode = cathedral
				c := NewCathedral(scoreText.Val)
				l := NewLaunchPad()
				r := NewRocket(scoreText.Val)

				allSprites.Sprites = append(allSprites.Sprites, c)
				allSprites.Sprites = append(allSprites.Sprites, l)
				allSprites.Sprites = append(allSprites.Sprites, r)
			}
		}
		return
	} else if gamemode == paused && !s.Paused {
		// fill the well when we're paused
		var filledWell string
		for n := 0; n <= 20; n++ {
			filledWell += wellLine + "\n"
		}
		c := sprite.NewCostume(filledWell, '~')
		s.PauseSprite = sprite.NewBaseSprite(s.X+2, s.Y, c)
		allSprites.Sprites = append(allSprites.Sprites, s.PauseSprite)
		s.Paused = true
		return
	}

	// remove the bricks
	if gamemode == play && s.Paused {
		allSprites.Remove(s.PauseSprite)
		s.Paused = false
	}

	if activeTetromino.Stopped && !activeTetromino.Dead {
		s.TimeOut = areToHeight[activeTetromino.BottomEdgeHeight()]
		s.Timer = 0
		activeTetromino.Dead = true
		CheckRows()
		linesText.UpdateLines(TotalLines)
	} else if activeTetromino.Dead {
		if s.Timer >= s.TimeOut {
			nextScore = 0
			activeTetromino = nextTetromino
			activeTetromino.PlaceInWell()
			nextTetromino = getRandTetromino(src, background)
			nextTetromino.X = 45
			allSprites.Sprites = append(allSprites.Sprites, nextTetromino)
		}
		s.Timer++
	}
	Vaccuum()
}
