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
var NextTetromino *Tetromino

func getRandTetromino(src rand.Source) *Tetromino {
	r := rand.New(src)
	i := r.Intn(7)

	var t *Tetromino

	switch {
	case i == 0:
		t = NewL()
	case i == 1:
		t = NewJ()
	case i == 2:
		t = NewT()
	case i == 3:
		t = NewS()
	case i == 4:
		t = NewZ()
	case i == 5:
		t = NewI()
	case i == 6:
		t = NewSq()
	}

	t.X = 3
	t.Y = 10
	t.Stopped = true

	return t
}

func main() {
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

	activeTetromino = getRandTetromino(src)
	activeTetromino.Stopped = false
	activeTetromino.X = 20
	activeTetromino.Y = 7

	nextTetromino = getRandTetromino(src)

	allSprites.Sprites = append(allSprites.Sprites, activeTetromino)
	allSprites.Sprites = append(allSprites.Sprites, nextTetromino)

	bg := NewWell()

mainloop:
	for {
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
					activeTetromino.Timer = 20
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
				nextTetromino = getRandTetromino(src)
				allSprites.Sprites = append(allSprites.Sprites, nextTetromino)
			}
			allSprites.Render()
			time.Sleep(50 * time.Millisecond)
		}
	}
}
