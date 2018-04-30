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

	if i == 0 {
		t = NewL()
	} else if i == 1 {
		t = NewJ()
	} else if i == 2 {
		t = NewT()
	} else if i == 3 {
		t = NewS()
	} else if i == 4 {
		t = NewZ()
	} else if i == 5 {
		t = NewI()
	} else if i == 6 {
		t = NewSq()
	}

	t.X = 20
	t.Y = 7

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

	var t *Tetromino
	var seed int64
	seed = time.Now().Unix()

	src := rand.NewSource(seed)

	t = getRandTetromino(src)

	allSprites.Sprites = append(allSprites.Sprites, t)

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
					t.RotateClockwise()
				} else if ev.Key == tm.KeyArrowLeft {
					t.MoveLeft()
				} else if ev.Key == tm.KeyArrowRight {
					t.MoveRight()
				} else if ev.Key == tm.KeyArrowDown {
					t.Timer = 20
				}
			} else if ev.Type == tm.EventResize {
				Width = ev.Width
				Height = ev.Height
			}
		default:
			bg.Render()
			allSprites.Update()
			if t.Stopped == true {
				CheckRows()
				t = getRandTetromino(src)
				allSprites.Sprites = append(allSprites.Sprites, t)
			}
			allSprites.Render()
			time.Sleep(50 * time.Millisecond)
		}
	}
}
