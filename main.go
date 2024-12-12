package main

import (
	"image/color"
	"log"
	"time"

	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	World              World
	Camera             [2]int
	HasInitiatedUpdate bool
	HasInitiatedDraw   bool
	Frames             int
	SecondTimer        time.Time
	ActualFPS          float32
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (game *Game) Update() error {
	// get input
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		game.Camera[1] += cameraMoveSensitivity
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		game.Camera[1] -= cameraMoveSensitivity
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		game.Camera[0] += cameraMoveSensitivity
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		game.Camera[0] -= cameraMoveSensitivity
	}

	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	if !game.HasInitiatedDraw {
		game.HasInitiatedDraw = true
		game.Camera[0] = screen.Bounds().Dx() / 2
		game.Camera[1] = screen.Bounds().Dy() / 2
	}

	screen.Fill(color.RGBA{0, 0, 88, 255})

	game.World.GetChunk(0, 0).Render(screen, game.Camera[0], game.Camera[1])
	drawString(screen, "Isomicraft Indev", 0, 10)
	drawString(screen, fmt.Sprintf("Camera: %d, %d", game.Camera[0], game.Camera[1]), 0, 22)
	drawString(screen, fmt.Sprintf("FPS: %f", game.ActualFPS), 0, 34)
	drawString(screen, fmt.Sprintf("TPS: %f", ebiten.ActualTPS()), 0, 46)

	game.Frames++
	if time.Since(game.SecondTimer).Seconds() >= 1 {
		game.ActualFPS = float32(game.Frames)
		game.Frames = 0
		game.SecondTimer = time.Now()
	}

}

func main() {
	// init ebiten
	initialWidth, initialHeight := 640, 480
	ebiten.SetTPS(60)
	// SetVsyncEnabled(false) is only for debug purposes
	// It is greatly reccomended to use SetVsyncEnabled(true) in production
	ebiten.SetVsyncEnabled(false)
	ebiten.SetWindowSize(initialWidth, initialHeight)
	ebiten.SetWindowTitle("Isometric Voxel Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// init game
	game := &Game{
		HasInitiatedDraw:   false,
		HasInitiatedUpdate: false,
	}
	initRender()
	game.World = World{
		Chunks: make(map[[2]int]Chunk),
	}
	game.World.Chunks[[2]int{0, 0}] = generateChunk(16, 16, 16, defaultVoxelDictionary, 0)

	// run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
