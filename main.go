package main

import (
	"math/rand"
	"time"

	sprite "github.com/pdevine/go-asciisprite"

	tm "github.com/nsf/termbox-go"
)

var allSprites sprite.SpriteGroup
var allBlocks []*TetrominoBlock
var Width int
var Height int
var CurrentLevel int
var TotalLines int

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

	var activeTetromino *Tetromino
	var nextTetromino *Tetromino
	var seed int64
	seed = time.Now().Unix()

	src := rand.NewSource(seed)
	bg := NewWell()

	activeTetromino = getRandTetromino(src, bg)
	activeTetromino.Stopped = false
	activeTetromino.X = 20
	activeTetromino.Y = 7
	activeTetromino.SetGravity(CurrentLevel)

	nextTetromino = getRandTetromino(src, bg)

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
					activeTetromino.Timer += levelFPG[CurrentLevel] / 2
				}
			} else if ev.Type == tm.EventResize {
				Width = ev.Width
				Height = ev.Height
			}
		default:
			bg.Render()
			allSprites.Update()
			if activeTetromino.Stopped == true {
				CheckRows()
				activeTetromino = nextTetromino
				activeTetromino.Stopped = false
				activeTetromino.X = 20
				activeTetromino.Y = 7
				activeTetromino.SetGravity(CurrentLevel)
				nextTetromino = getRandTetromino(src, bg)
				allSprites.Sprites = append(allSprites.Sprites, nextTetromino)
			}
			allSprites.Render()
			elapsed := time.Since(start)
			time.Sleep(time.Second/60 - elapsed)
		}
	}
}
