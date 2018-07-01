# tetromino
Tetromino Elektronika -- an ASCII Falling Tetromino Game
by Patrick Devine, 2018

TL;DR - A falling tetromino block game based upon NES Tetris in the style of the original 1984 Tetris by Alexey Pajitnov

Try it out with the command:

```
  $ docker run -it --rm pdevine/tetromino
```

This game was written because of an offhand remark my boss made to me in 1994.  I had just finished up two years of college
in Vancouver, British Columbia when I had taken a job at Radical Entertainment to write Super Nintendo games.  I wasn't
a particularly great programmer at the time (which is arguably still true), however, I found writing games in 65816 Assembly
(the language of choice for the Super Nintendo) extremely daunting.  There was sparse documentation, few code samples, and
the assembly that I had learned in college was fairly different than the code I was writing at Radical.  That coupled with
a slow build process on an Intel 486DX2-66 computer, and an even slower data transfer of the object code across a parallel
port onto our homebuilt development kits, meant a lot of 80 hour work weeks.

I remember mentioning to my boss, Jack Rebbetoy, about how difficult I was finding things, and he made a remark about how
we had engineers on-staff who could knock out a copy of Tetris in a weekend.  This really blew my mind.  We did have a lot of
really talented engineers at Radical, but at the time I was having a hard time even getting a sprite to render correctly on
the screen.

Fast forward 24 years, and I'm now working as a Product Manager at Docker.  I still love writing games, but haven't done it
professionally since that job at Radical.  The problem with using Docker to run and distribute games though, is that there
isn't an easy to use video subsystem.  You can use it with X11 and using a remote display, or just passing in the X socket,
but that's not exactly easy to do.  So given the limitation to really only use ASCII text, can you still create meaningful
games and applications?

I ended up putting together a very rudimentary ASCII sprite library (github.com/pdevine/go-asciisprite), written in Go.  Go
was chosen since I wanted the compiled binary to be really small.  I think the entire Tetromino docker image weighs in at
about 2.8MB, which is shameful by 1994 standards (two whole floppy disks!), but pretty tiny in 2018.  After putting together
the sprite library, the choice for the first game to write seemed pretty natural;  Could I actually write Tetris in a
weekend?

The answer is a pretty definitive "yes".  Even though I've spent more than a weekend writing Tetromino 
onika, the core of the game was written in about that amount of time, mostly spent riding Caltrain up and down the peninsula
in the SF Bay Area.  I tried to keep things as close as possible in terms of timings with the original Classic NES Tetris game
(sorry, no hard drops!), and tried to make the look-and-feel similar to the original Tetris from 1984 on the old Soviet based
Electronika 60.

So with that, hopefully you enjoy the game, and maybe even get inspired to create something with go-asciisprite.
