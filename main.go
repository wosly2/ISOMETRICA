package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	World               World
	Camera              [2]int
	HasInitiatedUpdate  bool
	HasInitiatedDraw    bool
	Frames              int
	SecondTimer         time.Time
	ActualFPS           float32
	DepthShift          float32
	UsingDepthShift     bool
	DepthShiftDirection bool
	CurrentChunk        [2]int
	ChunkSize           int
	ChunkDepth          int
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (game *Game) Update() error {
	if !game.HasInitiatedUpdate {
		game.HasInitiatedUpdate = true
	}

	// INPUT

	// camera movement
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		game.Camera[1] += cameraMoveSensitivity
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		game.Camera[1] -= cameraMoveSensitivity
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		game.Camera[0] += cameraMoveSensitivity
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		game.Camera[0] -= cameraMoveSensitivity
	}

	// toggle depth shift
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		game.UsingDepthShift = true
	} else {
		game.UsingDepthShift = false
		game.DepthShift = 0
	}

	// make camera move faster
	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		cameraMoveSensitivity = 8
	} else {
		cameraMoveSensitivity = 4
	}

	// reload world
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		game.World.Initalize(rand.Int63n(1000000))
		// regenerate the chunk and its neighbors
		for x := -1; x < 2; x++ {
			for y := -1; y < 2; y++ {
				game.World.generateChunk([2]int{x + game.CurrentChunk[0], y + game.CurrentChunk[1]}, game.ChunkSize, game.ChunkSize, game.ChunkDepth, defaultVoxelDictionary)
			}
		}
	}

	// END OF INPUT

	// depth shift
	if game.UsingDepthShift {
		if game.DepthShift > 1 {
			game.DepthShiftDirection = false
		} else if game.DepthShift < -1 {
			game.DepthShiftDirection = true
		}

		if game.DepthShiftDirection {
			game.DepthShift += .02
		} else {
			game.DepthShift -= .02
		}
	}

	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	// initiate
	if !game.HasInitiatedDraw {
		game.HasInitiatedDraw = true
		game.Camera[0] = screen.Bounds().Dx() / 2
		game.Camera[1] = screen.Bounds().Dy() / 2
	}

	// fill background
	screen.Fill(color.RGBA{0, 0, 88, 255})

	// render the chunks
	var blocksRendered int
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			chunk, exists := game.World.Chunks[[2]int{x + game.CurrentChunk[0], y + game.CurrentChunk[1]}]
			if exists {
				// get screen position of the top voxel in the chunk
				screenX, screenY := getScreenPosition(
					x*game.ChunkSize,
					y*game.ChunkSize,
					0,
					game.Camera[0],
					game.Camera[1],
					game.DepthShift,
				)
				// render
				blocksRendered += chunk.Render(screen, screenX, screenY, game.DepthShift)
			} else {
				// create chunk
				log.Printf("Creating chunk at %d, %d ...", x+game.CurrentChunk[0], y+game.CurrentChunk[1])
				game.World.generateChunk([2]int{x + game.CurrentChunk[0], y + game.CurrentChunk[1]}, game.ChunkSize, game.ChunkSize, game.ChunkDepth, defaultVoxelDictionary)
				log.Printf("Done!\n")
			}
		}
	}

	// gui/text
	drawString(screen, "Isomicraft Infdev", 0, 10)
	drawString(screen, fmt.Sprintf("Camera: %d, %d", game.Camera[0], game.Camera[1]), 0, 22)
	drawString(screen, fmt.Sprintf("FPS: %f", game.ActualFPS), 0, 34)
	drawString(screen, fmt.Sprintf("TPS: %f", ebiten.ActualTPS()), 0, 46)
	drawString(screen, fmt.Sprintf("Blocks Rendered: %d", blocksRendered), 0, 58)

	// frame control
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
	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowSize(initialWidth, initialHeight)
	ebiten.SetWindowTitle("ISOMICRAFT Infdev")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// init game
	game := &Game{
		HasInitiatedDraw:   false,
		HasInitiatedUpdate: false,
		CurrentChunk:       [2]int{0, 0},
		ChunkSize:          32,
		ChunkDepth:         64,
	}
	initRender()
	log.Println("World initializing ...")
	game.World = World{
		Chunks: make(map[[2]int]Chunk),
		Seed:   4311080085,
	}
	game.World.Initalize(4311080085)
	log.Println("Done!")

	// run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
