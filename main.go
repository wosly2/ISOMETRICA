package main

import (
	"log"
	"math/rand/v2"
	"net/http"
	"time"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"

	_ "net/http/pprof"
)

// represents the current game state
type GameState int

const (
	GAMESTATE_TITLE GameState = 2
	GAMESTATE_MENU  GameState = 1
	GAMESTATE_GAME  GameState = 0
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
	Framebuffer         *ebiten.Image
	DebugMode           bool
	ScreenX             int
	ScreenY             int
	GameState           GameState
}

// MapSize returns the size of a map in bytes.
func MapSize[K comparable, V any](m map[K]V) uintptr {
	// get the size of the map header
	headerSize := unsafe.Sizeof(m)

	// iterate through the map to calculate the total size of keys and values
	var keysSize, valuesSize uintptr
	for k, v := range m {
		keysSize += unsafe.Sizeof(k)
		valuesSize += unsafe.Sizeof(v)
	}

	// calculate the total size
	totalSize := headerSize + keysSize + valuesSize
	return totalSize
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, (640 * outsideHeight) / outsideWidth
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// init game
	game := &Game{
		HasInitiatedDraw:   false,
		HasInitiatedUpdate: false,
		CurrentChunk:       [2]int{0, 0},
		ChunkSize:          32,
		ChunkDepth:         64,
		GameState:          GAMESTATE_TITLE,
	}

	// init ebiten
	initialWidth, initialHeight := 640, 480
	ebiten.SetTPS(60)
	// SetVsyncEnabled(false) is only for debug purposes
	// It is greatly reccomended to use SetVsyncEnabled(true) in production
	ebiten.SetVsyncEnabled(!game.DebugMode)
	ebiten.SetWindowSize(initialWidth, initialHeight)
	ebiten.SetWindowTitle("ISOMETRICA Infdev")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// init render
	initRender()

	// check if there is a save file at the save path
	savePath := "save/demo"
	world, err := readWorld(savePath)
	if err != nil {
		log.Printf("Failed to load world: %v", err)
		// create new world
		log.Println("Creating new world ...")
		game.World = World{
			Chunks:     make(map[[2]int]Chunk),
			ChunkSize:  32,
			ChunkDepth: 64,
			SavePath:   "save/demo",
		}
		game.World.Initalize(rand.Int64())
	} else {
		log.Println("Loading world from file.")
		game.World = world
	}

	// init player
	game.Player = Player{
		Position: Vec3{0, 0, float32(game.World.SurfaceFeaturesBeginAt) + 10},
		Velocity: Vec3{0, 0, 0},
		Drag:     Vec3{.9, .9, .9},
		Texture:  "Default",
	}

	// run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
