package goquarium

import (
	"github.com/ansoni/termination"
	"github.com/nsf/termbox-go"
	"os"
	"math/rand"
	"time"
)

type Goquarium struct {	
	term *termination.Termination
	surface []*termination.Entity
        ripples []*termination.Entity
        seaweeds []*termination.Entity
        fishes []*termination.Entity
	castle *termination.Entity
	bubbles *termination.Entity

}

func (goq *Goquarium) generateBubbles() {
	term := goq.term
	for {
		rand.Seed(int64(time.Now().Nanosecond()))
		randX := random(0, term.Width)
		randY := random(17, term.Height)
		
		bubble := term.NewEntity(termination.Position{randX, randY, 10})
		bubble.Shape = bubbleShape
		bubble.MovementCallback = termination.UpMovement
		bubble.DefaultColor = 'c'
		bubble.FramesPerSecond = 2
		bubble.DeathOnLastFrame = true
		time.Sleep(3 * time.Second)
	}
}

func (goq *Goquarium) deadFish(term *termination.Termination, entity *termination.Entity) {
	goq.addFish()
}

func (goq *Goquarium) deadWhale(term *termination.Termination, entity *termination.Entity) {
	goq.addWhale()
}

func (goq *Goquarium) addWhale() {
	direction := []string{"left","right"}[random(0, 2)]
	width := goq.term.Width
	position := termination.Position{-10, 0, 50}	
	if direction == "left" {
		position = termination.Position{width+10, 0, 50}
	}

	whale := goq.term.NewEntity(position)
	whale.Shape = whaleShape
	whale.DeathOnOffScreen = true
	whale.ColorMask = whaleMask
	whale.ShapePath = direction
	whale.DefaultColor = 'b'
	whale.DeathCallback = goq.deadWhale
	if direction == "left" {
		whale.MovementCallback = termination.LeftMovement
	} else {
		whale.MovementCallback = termination.RightMovement
	}
	//fish.Death = goq.deadFish
	whale.FramesPerSecond = 10
}


func (goq *Goquarium) addFish() {
	rand.Seed(int64(time.Now().Nanosecond()))
	fishSelection := random(0,len(fishShapes))
	fishShape := fishShapes[fishSelection]
	fishMask := fishMasks[fishSelection]

	direction := []string{"left","right"}[random(0, 2)]

	// we unfortunately have to iterate to get its height before insert...
	shapeData := []rune(fishShape["left"][0])
	shapeHeight := 1
	for _, char := range shapeData {
		if char == '\n' {
			shapeHeight += 1
		}
	}

	height := goq.term.Height-shapeHeight
	width := goq.term.Width
	randY := random(9, height)
	speed := random(5, 10)
	// TODO: lots of duplicate code here
	if direction == "left" {
		fish := goq.term.NewEntity(termination.Position{width+10, randY, 5})
		fish.Shape = fishShape
		fish.DeathOnOffScreen = true
		fish.ColorMask = fishMask
		fish.ShapePath = "left"
		fish.MovementCallback = termination.LeftMovement
		fish.DeathCallback = goq.deadFish
		fish.FramesPerSecond = speed
		goq.fishes = append(goq.fishes, fish)
	} else {
		fish := goq.term.NewEntity(termination.Position{-10, randY, 5})
		fish.DeathOnOffScreen = true
		fish.Shape = fishShape
		fish.DeathCallback = goq.deadFish
		fish.ColorMask = fishMask
		fish.ShapePath = "right"
		fish.MovementCallback = termination.RightMovement
		fish.FramesPerSecond = speed
		goq.fishes = append(goq.fishes, fish)
	}
}

func (goq *Goquarium) generateFishes() {
	screenSize := (goq.term.Height - 9) * goq.term.Width
	fishCount := int(screenSize / 200)

	// keep adding fish when we need
	for i:=0;i<fishCount;i++ {
		goq.addFish()
		time.Sleep(500 * time.Millisecond) // space out the fish a bit
	}
}

func (goq *Goquarium) setupEnvironment() {
	term := goq.term
	top_y := 5
	height := term.Height
	width := term.Width

	needed := width/4
	for i := 0;i < needed;i++ {
		surface := term.NewEntity(termination.Position{i*4, top_y, 10})
		surface.Shape = surfaceShape
		surface.ColorMask = waterMask
		surface.FramesPerSecond = 1
		goq.surface = append(goq.surface, surface)

		ripplePaths := []string{"a","b","c","d"}
		rand := random(0,4)
		ripples := term.NewEntity(termination.Position{i*4, top_y+1, 10})
		ripples.Shape = rippleShape
		ripples.ColorMask = waterMask
		ripples.ShapePath = ripplePaths[rand]
		ripples.FramesPerSecond = 1
		goq.ripples = append(goq.ripples, surface)
	}

	//castle! castle is 13 high... if we change it... gonna suck here
	goq.castle = term.NewEntity(termination.Position{width-31, height-13, -1})
	goq.castle.Shape = castleShape
	goq.castle.FramesPerSecond = 1
	goq.castle.ShapePath = "default"
	goq.castle.ColorMask = castleMask
	goq.castle.DefaultColor = 'W'

	//grass is a bit random
	seaweedCount := width/15
	for i := 0;i < seaweedCount;i++ {
	        rand.Seed(int64(time.Now().Nanosecond()))
		seaweedHeight := random(1,7)
		seaweedX := random(0, width)
		path := []string{"a","b"}
		h := 0
	        for j := seaweedHeight;j >= 0;j-- {
			x:=seaweedX
			y:=height - seaweedHeight
			if x > width {
				x=width
			}
			
		        seaweed := term.NewEntity(termination.Position{x, y+j, 10})
			seaweed.Shape = seaweedShape
			seaweed.FramesPerSecond = 2
			seaweed.ShapePath = path[h]
			seaweed.ColorMask = seaweedMask
			goq.seaweeds = append(goq.seaweeds, seaweed)
			if h == 0 {
				h = 1 
			} else {
				h = 0
			}
		}
	}

	
	
}

func random(min int, max int) int {
    return rand.Intn(max-min) + min
}

func Fish() {
	rand.Seed(int64(time.Now().Nanosecond()))

	goquarium := Goquarium{}
	term := termination.New()
	//term.Debug="./debug.out"
	goquarium.term = term
	term.FramesPerSecond = 10
	goquarium.setupEnvironment()

	go goquarium.generateBubbles()
	go goquarium.generateFishes()
	goquarium.addWhale()
	go term.Animate()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
				term.Close()
				os.Exit(0)
			}
		case termbox.EventMouse:
			//update_mouse(mouse, &ev)
		}
	}
}

