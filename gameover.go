package main

import sprite "github.com/pdevine/go-asciisprite"

const basil = `                           .
                           T
                          ( )
                          <==>
                           FJ
                           ==
                          J||F
                          F||J
                         /\/\/\
                         F++++J
                        J{}{}{}F         .
                     .  F{}{}{}J         T
          .          T J{}{}{}{}F        ;;
          T         /|\F \/ \/ \J  .   ,;;;;.
         /:\      .'/|\\:========F T ./;;;;;;\
       ./:/:/.   ///|||\\\""""""" /x\T\;;;;;;/
      //:/:/:/\  \\\\|////..[]...xXXXx.|====|
      \:/:/:/:T7 :.:.:.:.:||[]|/xXXXXXx\|||||
      ::.:.:.:A. ` + "`" + `;:;:;:;'=====\XXXXXXX/=====.
      ` + "`" + `;""::/xxx\.|,|,|,| ( )( )| | | |.=..=.|
       :. :` + "`" + `\xxx/(_)(_)(_) _  _ | | | |'-''-'|
       :T-'-.:"":|"""""""|/ \/ \|=====|======|
       .A."""||_|| ,. .. || || |/\/\/\/ | | ||
   :;:////\:::.'.| || || ||-||-|/\/\/\+|+| | |
  ;:;;\////::::,='======='============/\/\=====.
:;:::;""":::::;:|__..,__|===========/||\|\====|
:::::;|=:::;:;::|,;:::::         |========|   |
::l42::::::(}:::::;::::::________|========|___|__`

const pad = `
          ,========.
          |	       |
      ====|\      /|                =======
       \  | ,-.  / |======          XXX|XXX
        \ | |_| /  |------          XXX|XXX
         \|   \/   |                XXX|XXX
==========|   /\   |                XXX|XXX
----------|  / ,--.|                XXX|XXX
   ,======| /  |  ||                XXX|XXX
===  -----|/   |  ||                XXX|XXX
----' '=====================================.
      ||        ||       ||       ||       ||`

const rocket1 = `xx_
x/ \
x| |
/| |\
x/_\`
const rocket2 = ` _
/ \
| |
| |
| |
| |
===
/_\`

const rocket3 = `xx_
x/ \
x| |
x| |
/| |\
|| ||
|===|
|/_\|`

const rocket4 = `xxx__
xx/  \
xx|  |
x/|__|\
| /__\ |
| |--| |
|/|  |\|
/ |  | \
| |  | |
==/||\==
/_\  /_\`

const rocket5 = `xxxx__
xxx/o \
==========
x\______/`

type Cathedral struct {
	sprite.BaseSprite
	Timer   int
	TimeOut int
	state   uint
	spark   *Spark
}

func NewCathedral(score int) *Cathedral {
	s := &Cathedral{BaseSprite: sprite.BaseSprite{
		X:       46,
		Y:       10,
		Visible: true,
	},
		TimeOut: 20,
		state:   hovering,
	}

	s.AddCostume(sprite.NewCostume(basil, 'q'))

	if score >= 120000 {
		s.state = igniting
	}

	return s
}

func (s *Cathedral) Update() {
	s.Timer++
	if s.state == igniting {
		if s.Timer >= s.TimeOut {
			s.state = flying
			s.TimeOut = 6
			s.Timer = 0
			f := NewSpark(s.X+s.Width/2, s.Y+s.Height+1, true)
			s.spark = f
			allSprites.Sprites = append(allSprites.Sprites, f)
		}
	} else if s.state == flying {
		if s.Timer >= s.TimeOut {
			s.Y--
			s.spark.Y--
			s.Timer = 0
		}
	}
}

func NewLaunchPad() *sprite.BaseSprite {
	c := sprite.NewCostume(pad, 'q')
	s := sprite.NewBaseSprite(0, 26, c)
	return s
}

const (
	igniting = iota
	flying
	hovering
)

type Rocket struct {
	sprite.BaseSprite
	Timer   int
	TimeOut int
	state   uint
	sparks  []*Spark
}

func NewRocket(score int) *Rocket {
	var costume int
	var x, y int

	switch {
	case score >= 120000:
		costume = 4
		x = 23
		y = 29
	case score >= 100000:
		costume = 3
		x = 24
		y = 26
	case score >= 60000:
		costume = 2
		x = 25
		y = 28
	case score >= 40000:
		costume = 1
		x = 26
		y = 28
	case score >= 30000:
		costume = 0
		x = 25
		y = 28
	}

	s := &Rocket{BaseSprite: sprite.BaseSprite{
		X:       x,
		Y:       y,
		Visible: true,
	},
		TimeOut: 20,
	}

	for _, r := range []string{rocket1, rocket2, rocket3, rocket4, rocket5} {
		s.AddCostume(sprite.NewCostume(r, 'x'))
	}
	s.SetCostume(costume)
	if costume == 4 {
		s.state = hovering
	}

	return s
}

func (s *Rocket) Update() {
	if s.state == igniting {
		if s.Timer >= s.TimeOut {
			s.state = flying
			s.TimeOut = 6
			s.Timer = 0
			// XXX - clean this up
			switch {
			case s.CurrentCostume == 0:
				f := NewSpark(s.X+1, s.Y+s.Height+1, false)
				s.sparks = append(s.sparks, f)
				allSprites.Sprites = append(allSprites.Sprites, f)
			case s.CurrentCostume == 1:
				f := NewSpark(s.X, s.Y+s.Height+1, false)
				s.sparks = append(s.sparks, f)
				allSprites.Sprites = append(allSprites.Sprites, f)
			case s.CurrentCostume == 2:
				f := NewSpark(s.X+1, s.Y+s.Height+1, false)
				s.sparks = append(s.sparks, f)
				allSprites.Sprites = append(allSprites.Sprites, f)
			case s.CurrentCostume == 3:
				f1 := NewSpark(s.X, s.Y+s.Height+1, false)
				f2 := NewSpark(s.X+s.Width-3, s.Y+s.Height+1, false)
				s.sparks = append(s.sparks, f1)
				s.sparks = append(s.sparks, f2)
				allSprites.Sprites = append(allSprites.Sprites, f1)
				allSprites.Sprites = append(allSprites.Sprites, f2)
			}
			return
		}
	} else if s.state == flying {
		if s.Timer >= s.TimeOut {
			s.Y--
			s.Timer = 0
			for _, f := range s.sparks {
				f.Y--
			}
		}
	}
	s.Timer++
}

type Spark struct {
	sprite.BaseSprite
	Timer   int
	TimeOut int
}

const flame1 = `OoO
 o,
.`

const flame2 = `oOO
,o
 .`

const flame3 = `OOo
.o+
 ,
`

const bigFlame1 = `
OoOOooO
 oOOoO
 ooO,
 .,
  .`

const bigFlame2 = `
oO+oOO+
 ,ooOo
  ,Oo.
   ,+
  ,`

const bigFlame3 = `
OoOOoOO
 +OO,O
  .oO,
  .,
   .`

func NewSpark(x, y int, big bool) *Spark {
	s := &Spark{BaseSprite: sprite.BaseSprite{
		X:       x,
		Y:       y,
		Visible: true,
	},
		TimeOut: 3,
	}

	if !big {
		for _, c := range []string{flame1, flame2, flame3} {
			s.AddCostume(sprite.NewCostume(c, ';'))
		}
	} else {
		for _, c := range []string{bigFlame1, bigFlame2, bigFlame3} {
			s.AddCostume(sprite.NewCostume(c, ';'))
		}
	}

	return s
}

func (s *Spark) Update() {
	s.Timer++
	if s.Timer >= s.TimeOut {
		s.NextCostume()
		s.Timer = 0
	}
}
