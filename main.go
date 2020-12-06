package main

import (
	"math/rand"
	"flag"
	"time"

	sprite "github.com/pdevine/go-asciisprite"

	tm "github.com/pdevine/go-asciisprite/termbox"
)

var allSprites sprite.SpriteGroup
var allBlocks []*TetrominoBlock
var background *Well
var Width int
var Height int
var CurrentLevel *LevelText
var TotalLines int
var activeTetromino *Tetromino
var nextScore int
var nextTetromino *Tetromino
var linesText *LinesText
var scoreText *ScoreText
var src rand.Source

const (
	title = iota
	levelselect
	play
	paused
	gameover
	cathedral
)

var gamemode = title

func main() {
	// XXX - hack to make this work inside of a Docker container
	time.Sleep(1000 * time.Millisecond)

	var seed int64

	debug := flag.Bool("debug", false, "")
	flag.Int64Var(&seed, "seed", -1, "")
	flag.Parse()

	err := tm.Init()
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	Width, Height = tm.Size()

	event_queue := make(chan tm.Event)
	go func() {
		for {
			event_queue <- tm.PollEvent()
		}
	}()

	if seed == -1 {
		seed = time.Now().Unix()
	}

	src = rand.NewSource(seed)
	background = NewWell()
	background.TimeOut = 20
	linesText = NewLinesText(0)
	linesText.Y = background.Y - 2
	linesText.X = background.X + 8
	scoreText = NewScoreText()
	helpText := NewHelpText()

	t := NewTitle()
	tstr := NewTitleString()
	allSprites.Sprites = append(allSprites.Sprites, t)
	allSprites.Sprites = append(allSprites.Sprites, tstr)

	selector := NewSelector()

mainloop:
	for {
		start := time.Now()
		tm.Clear(tm.ColorDefault, tm.ColorDefault)

		select {
		case ev := <-event_queue:
			if ev.Type == tm.EventKey {
				if ev.Key == tm.KeyEsc {
					break mainloop
				}
				if gamemode == title {
					gamemode = levelselect
					allSprites.Sprites = append(allSprites.Sprites, selector)
				} else if ev.Key == tm.KeyEnter {
					if gamemode == levelselect {
						allSprites.RemoveAll()
						NewStats()
						CurrentLevel = NewLevelText(selector.GetVal())
						activeTetromino = getRandTetromino(src, background)
						activeTetromino.PlaceInWell()
						nextTetromino = getRandTetromino(src, background)
						nextTetromino.X = 45
						allSprites.Sprites = append(allSprites.Sprites, helpText)
						allSprites.Sprites = append(allSprites.Sprites, linesText)
						allSprites.Sprites = append(allSprites.Sprites, scoreText)
						allSprites.Sprites = append(allSprites.Sprites, CurrentLevel)
						allSprites.Sprites = append(allSprites.Sprites, activeTetromino)
						allSprites.Sprites = append(allSprites.Sprites, nextTetromino)
						gamemode = play
					}
				} else if ev.Key == tm.KeySpace {
					if gamemode == levelselect {
						selector.NextCostume()
						if selector.currentLevel < 10 {
							selector.SetLevel(selector.currentLevel + 10)
						} else {
							selector.SetLevel(selector.currentLevel - 10)
						}
					} else if gamemode == play {
						activeTetromino.RotateClockwise()
					}
				} else if ev.Key == tm.KeyArrowLeft {
					if gamemode == levelselect {
						selector.MoveLeft()
					} else if gamemode == play {
						activeTetromino.MoveLeft()
					}
				} else if ev.Key == tm.KeyArrowRight {
					if gamemode == levelselect {
						selector.MoveRight()
					} else if gamemode == play {
						activeTetromino.MoveRight()
					}
				} else if ev.Key == tm.KeyArrowUp {
					if gamemode == levelselect {
						selector.MoveUp()
					} else if gamemode == play {
						activeTetromino.Drop()
					}
				} else if ev.Key == tm.KeyArrowDown {
					if gamemode == levelselect {
						selector.MoveDown()
					} else if gamemode == play {
						activeTetromino.Timer += levelFPG[CurrentLevel.Val] / 2
						nextScore += 4
					}
				} else if ev.Ch == 'p' || ev.Ch == 'P' {
					if gamemode == paused {
						gamemode = play
					} else if gamemode == play {
						gamemode = paused
					}
				} else if *debug && ev.Ch == '$' {
					scoreText.AddVal(10000)
				}
			} else if ev.Type == tm.EventResize {
				Width = ev.Width
				Height = ev.Height
			}
		default:
			allSprites.Update()
			if gamemode == play || gamemode == gameover || gamemode == paused {
				background.Update()
				background.Render()
			}
			allSprites.Render()
			elapsed := time.Since(start)
			time.Sleep(time.Second/60 - elapsed)
		}
	}
}
