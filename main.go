package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (game *Game) Update() error {
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	demoChunk.Render(screen, screen.Bounds().Dx()/2, screen.Bounds().Dy()/2)
}

func main() {
	// init ebiten
	game := &Game{}
	initialWidth, initialHeight := 640, 480
	ebiten.SetTPS(60)
	ebiten.SetWindowSize(initialWidth, initialHeight)
	ebiten.SetWindowTitle("Isometric Voxel Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// other init
	initRender()

	// run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
