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
	Camera              [2]float32
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
	Player              Player
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (game *Game) Update() error {
	if !game.HasInitiatedUpdate {
		game.HasInitiatedUpdate = true
	}

	// INPUT

	// player movement
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		game.Player.Velocity[1] += .1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		game.Player.Velocity[1] -= .1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		game.Player.Velocity[0] += .1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		game.Player.Velocity[0] -= .1
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		game.Player.Velocity[2] -= .1
	}
	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		game.Player.Velocity[2] += .1
	}

	// toggle depth shift
	if ebiten.IsKeyPressed(ebiten.KeyBackslash) {
		game.UsingDepthShift = true
	} else {
		game.UsingDepthShift = false
		game.DepthShift = 0
	}

	// reload world
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		game.World.Initalize(rand.Int63n(1000000))
		// regenerate the chunk and its neighbors
		for x := -1; x < 2; x++ {
			for y := -1; y < 2; y++ {
				game.World.generateChunk([2]int{x + game.CurrentChunk[0], y + game.CurrentChunk[1]}, game.World.ChunkSize, game.World.ChunkSize, game.World.ChunkDepth, defaultVoxelDictionary)
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

	// update player
	game.Player.Update(game.World)

	// get current chunk based on player position
	game.CurrentChunk[0] = int(game.Player.Position[0] / float32(game.World.ChunkSize))
	game.CurrentChunk[1] = int(game.Player.Position[1] / float32(game.World.ChunkSize))

	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	// initiate
	if !game.HasInitiatedDraw {
		game.HasInitiatedDraw = true
		game.Camera[0] = float32(screen.Bounds().Dx() / 2)
		game.Camera[1] = float32(screen.Bounds().Dy() / 2)
	}

	// fill background
	screen.Fill(color.RGBA{0, 0, 88, 255})

	// render the chunks
	var blocksRendered int
	for x := game.CurrentChunk[0] - 1; x < game.CurrentChunk[0]+2; x++ {
		for y := game.CurrentChunk[1] - 1; y < game.CurrentChunk[1]+2; y++ {
			chunk, exists := game.World.Chunks[[2]int{x, y}]
			if exists {
				// get screen position of the top voxel in the chunk
				screenX, screenY := getScreenPosition(
					(x)*game.World.ChunkSize,
					(y)*game.World.ChunkSize,
					0,
					game.Camera[0],
					game.Camera[1],
					game.DepthShift,
				)
				// render
				blocksRendered += chunk.Render(screen, float32(screenX), float32(screenY), game.DepthShift)
			} else {
				// create chunk
				log.Printf("Creating chunk at %d, %d ...", x, y)
				game.World.generateChunk([2]int{x, y}, game.World.ChunkSize, game.World.ChunkSize, game.World.ChunkDepth, defaultVoxelDictionary)
				log.Printf("Done!\n")
			}
		}
	}

	// render the player
	game.Player.Render(screen, game.DepthShift)

	// player screen position
	screenX, screenY := game.Player.getScreenPosition(game.DepthShift)

	// change camera position to have player in the center
	game.Camera[0] = screenX
	game.Camera[1] = screenY

	// gui/text
	drawString(screen, "Isomicraft Infdev", 0, 10)
	drawString(screen, fmt.Sprintf("Player Position: %f, %f, %f", game.Player.Position[0], game.Player.Position[1], game.Player.Position[2]), 0, 22)
	drawString(screen, fmt.Sprintf("Focused Chunk: %d, %d", game.CurrentChunk[0], game.CurrentChunk[1]), 0, 34)
	drawString(screen, fmt.Sprintf("FPS: %f TPS: %f", game.ActualFPS, ebiten.ActualTPS()), 0, 46)
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

	// init render
	initRender()

	// init world
	log.Println("World initializing ...")
	game.World = World{
		Chunks:     make(map[[2]int]Chunk),
		Seed:       4311080085,
		ChunkSize:  32,
		ChunkDepth: 64,
	}
	game.World.Initalize(4311080085)
	log.Println("Done!")

	// init player
	game.Player = Player{
		Position: [3]float32{0, 0, 0},
		Velocity: [3]float32{0, 0, 0},
		Drag:     [3]float32{.9, .9, .9},
		Texture:  "Default",
	}

	// run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
