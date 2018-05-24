package main

import (
	"math/rand"
	"time"

	sprite "github.com/pdevine/go-asciisprite"

	tm "github.com/nsf/termbox-go"
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
	play
	gameover
	cathedral
)

//var gamemode = gameover
var gamemode = play

func main() {
	// XXX - hack to make this work inside of a Docker container
	time.Sleep(500 * time.Millisecond)

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

	var seed int64
	seed = time.Now().Unix()

	src = rand.NewSource(seed)
	background = NewWell()
	background.TimeOut = 20
	NewStats()
	linesText = NewLinesText(0)
	linesText.Y = background.Y - 2
	linesText.X = background.X + 8
	CurrentLevel = NewLevelText(0)
	scoreText = NewScoreText()
	//scoreText.AddVal(140000)

	activeTetromino = getRandTetromino(src, background)
	activeTetromino.PlaceInWell()

	nextTetromino = getRandTetromino(src, background)
	nextTetromino.X = 45

	allSprites.Sprites = append(allSprites.Sprites, linesText)
	allSprites.Sprites = append(allSprites.Sprites, scoreText)
	allSprites.Sprites = append(allSprites.Sprites, CurrentLevel)
	allSprites.Sprites = append(allSprites.Sprites, activeTetromino)
	allSprites.Sprites = append(allSprites.Sprites, nextTetromino)

mainloop:
	for {
		start := time.Now()
		tm.Clear(tm.ColorDefault, tm.ColorDefault)

		select {
		case ev := <-event_queue:
			if ev.Type == tm.EventKey {
				if ev.Key == tm.KeyEsc {
					break mainloop
				} else if ev.Key == tm.KeySpace {
					activeTetromino.RotateClockwise()
				} else if ev.Key == tm.KeyArrowLeft {
					activeTetromino.MoveLeft()
				} else if ev.Key == tm.KeyArrowRight {
					activeTetromino.MoveRight()
				} else if ev.Key == tm.KeyArrowDown {
					activeTetromino.Timer += levelFPG[CurrentLevel.Val] / 2
					nextScore += 4
				}
			} else if ev.Type == tm.EventResize {
				Width = ev.Width
				Height = ev.Height
			}
		default:
			allSprites.Update()
			background.Update()
			if gamemode != cathedral {
				background.Render()
			}
			allSprites.Render()
			elapsed := time.Since(start)
			time.Sleep(time.Second/60 - elapsed)
		}
	}
}
