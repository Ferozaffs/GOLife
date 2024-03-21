package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Coord struct {
	x int
	y int
}

const (
	screenWidth   = 1280
	screenHeight  = 640
	scaleFactor   = 4
	surfaceWidth  = screenWidth / scaleFactor
	surfaceHeight = screenHeight / scaleFactor
	aliveColor    = 0x00
	deadColor     = 0x88
)

var offsets = []Coord{
	{x: -1, y: -1},
	{x: 0, y: -1},
	{x: 1, y: -1},
	{x: -1, y: 0},
	{x: 1, y: 0},
	{x: -1, y: 1},
	{x: 0, y: 1},
	{x: 1, y: 1},
}

var paused = true
var spacePressed = false

type Game struct {
	image        *ebiten.Image
	pixels       []byte
	stagedPixels []byte
}

func NewGame() *Game {
	g := &Game{
		image:        ebiten.NewImage(surfaceWidth, surfaceHeight),
		pixels:       make([]byte, surfaceWidth*surfaceHeight*4),
		stagedPixels: make([]byte, surfaceWidth*surfaceHeight*4),
	}

	for i := range g.stagedPixels {
		g.stagedPixels[i] = deadColor
	}

	copy(g.pixels, g.stagedPixels)
	g.image.WritePixels(g.pixels)
	return g
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if !spacePressed {
			paused = !paused
			spacePressed = true
		}
	} else {
		spacePressed = false
	}

	if !paused {
		g.CheckPixels()
	}

	g.PaintPixel()

	copy(g.pixels, g.stagedPixels)
	g.image.WritePixels(g.pixels)

	return nil
}

func Clamp(min, value, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}

	return value
}

func (g *Game) PaintPixel() {
	mx, my := ebiten.CursorPosition()
	x := Clamp(0, (mx / scaleFactor), surfaceWidth-1)
	y := Clamp(0, (my / scaleFactor), surfaceHeight-1)
	px := (y*surfaceWidth + x) * 4

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.stagedPixels[px] = aliveColor
		g.stagedPixels[px+1] = aliveColor
		g.stagedPixels[px+2] = aliveColor
		g.stagedPixels[px+3] = 0xFF
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.stagedPixels[px] = deadColor
		g.stagedPixels[px+1] = deadColor
		g.stagedPixels[px+2] = deadColor
		g.stagedPixels[px+3] = 0xFF
	}
}

func (g *Game) CheckPixels() {
	for i := 0; i < surfaceWidth; i++ {
		for j := 0; j < surfaceHeight; j++ {
			px := (j*surfaceWidth + i) * 4
			alive := g.pixels[px] == aliveColor
			neighbors := 0

			for _, c := range offsets {
				nx := i + c.x
				if nx < 0 {
					nx += surfaceWidth
				} else if nx >= surfaceWidth {
					nx -= surfaceWidth
				}

				ny := j + c.y
				if ny < 0 {
					ny += surfaceHeight
				} else if ny >= surfaceHeight {
					ny -= surfaceHeight
				}

				if g.pixels[(ny*surfaceWidth+nx)*4] == aliveColor {
					neighbors++
				}
			}

			if alive {
				if neighbors < 2 || neighbors > 3 {
					alive = false
				}
			} else {
				if neighbors == 3 {
					alive = true
				}
			}

			if alive {
				g.stagedPixels[px] = aliveColor
				g.stagedPixels[px+1] = aliveColor
				g.stagedPixels[px+2] = aliveColor
				g.stagedPixels[px+3] = 0xFF
			} else {
				g.stagedPixels[px] = deadColor
				g.stagedPixels[px+1] = deadColor
				g.stagedPixels[px+2] = deadColor
				g.stagedPixels[px+3] = 0xFF
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	scaleFactor := screenWidth / surfaceWidth

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(scaleFactor), float64(scaleFactor))
	screen.DrawImage(g.image, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("GOLife")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
