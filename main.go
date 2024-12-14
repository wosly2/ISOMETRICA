package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	World               World
	Camera              [2]int
	CameraOffset        [2]int
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
	CameraScale         float64
	CameraScaleTo       float64
	BlocksOffscreen     *ebiten.Image
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (game *Game) Update() error {
	if !game.HasInitiatedUpdate {
		game.HasInitiatedUpdate = true
	}

	// INPUT

	// IMPORTANT: this following code (the camera zoom code) is from the Ebitengine example, which uses Apache 2.0 license
	// I need to refactor the algorithm to be unique so I can so I don't have to add their license to my project
	// and i want to understand it more
	// https://github.com/hajimehoshi/ebiten/blob/main/examples/isometric/game.go

	// ebitengine example code
	// zoom camera
	var scrollY float64
	if ebiten.IsKeyPressed(ebiten.KeyC) || ebiten.IsKeyPressed(ebiten.KeyPageDown) {
		scrollY = -0.25
	} else if ebiten.IsKeyPressed(ebiten.KeyE) || ebiten.IsKeyPressed(ebiten.KeyPageUp) {
		scrollY = .25
	} else {
		_, scrollY = ebiten.Wheel()
		if scrollY < -1 {
			scrollY = -1
		} else if scrollY > 1 {
			scrollY = 1
		}
	}
	game.CameraScaleTo += scrollY * (game.CameraScaleTo / 7)

	// ebitengine example code
	// Clamp target zoom level.
	if game.CameraScaleTo < 0.8 {
		game.CameraScaleTo = 0.8
	} else if game.CameraScaleTo > 5 {
		game.CameraScaleTo = 5
	}

	// ebitengine example code
	// Smooth zoom transition.
	factor := math.Pow(10, 10)
	div := 20.0
	if game.CameraScaleTo > game.CameraScale {
		game.CameraScale += math.Round(((game.CameraScaleTo-game.CameraScale)/div)*factor) / factor
	} else if game.CameraScaleTo < game.CameraScale {
		game.CameraScale -= math.Round(((game.CameraScale-game.CameraScaleTo)/div)*factor) / factor
	}

	// end of ebitengine example code
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
		cameraMoveSensitivity = int(math.Round(8 / game.CameraScale))
	} else {
		cameraMoveSensitivity = int(math.Round(4 / game.CameraScale))
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
	game.BlocksOffscreen = ebiten.NewImage(int(math.Round(float64(screen.Bounds().Dx())/game.CameraScale))+1, int(math.Round(float64(screen.Bounds().Dy())/game.CameraScale))+1)
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			chunk, exists := game.World.Chunks[[2]int{x + game.CurrentChunk[0], y + game.CurrentChunk[1]}]
			if exists {
				// get screen position of the top voxel in the chunk
				screenX, screenY := getScreenPosition(
					x*game.ChunkSize,
					y*game.ChunkSize,
					0,
					game.Camera[0]+game.CameraOffset[0],
					game.Camera[1]+game.CameraOffset[1],
					game.DepthShift,
				)
				// render
				blocksRendered += chunk.Render(game.BlocksOffscreen, screenX, screenY, game.DepthShift)
			} else {
				// create chunk
				log.Printf("Creating chunk at %d, %d ...", x+game.CurrentChunk[0], y+game.CurrentChunk[1])
				game.World.generateChunk([2]int{x + game.CurrentChunk[0], y + game.CurrentChunk[1]}, game.ChunkSize, game.ChunkSize, game.ChunkDepth, defaultVoxelDictionary)
				log.Printf("Done!\n")
			}
		}
	}

	// scale the rendered image
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(game.CameraScale, game.CameraScale)
	//op.GeoM.Translate(float64(screen.Bounds().Dx()/2-game.BlocksOffscreen.Bounds().Dx()/2), float64(screen.Bounds().Dy()/2-game.BlocksOffscreen.Bounds().Dy()/2))
	screen.DrawImage(game.BlocksOffscreen, op)

	// get the camera offset to center it back back what it was before we zoom/scaled it
	//game.CameraOffset[0] = int(math.Round(float64(game.Camera[0]) / game.CameraScale))
	//game.CameraOffset[1] = int(math.Round(float64(game.Camera[1]) / game.CameraScale))

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
	game.CameraScale = 1
	game.CameraScaleTo = 1
	log.Println("Done!")

	// run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
