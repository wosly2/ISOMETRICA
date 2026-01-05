package main

import (
	"log"
	"math/rand/v2"
	"net/http"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"

	_ "net/http/pprof"
)

// represents the current game state
type GameState int

const (
	GAMESTATE_TITLE GameState = iota
	GAMESTATE_MENU
	GAMESTATE_GAME
)

type Game struct {
	GameState          GameState // gamestate enum
	HasInitiatedUpdate bool      // init update flag
	HasInitiatedDraw   bool      // init draw flag
	DebugMode          bool      // are we debugging?

	World        World  // in-game world position
	Player       Player // player context
	CurrentChunk [2]int // global chunk location of player

	Camera              [2]float32 // camera location
	Direction           [4]int     // rotation factor
	DepthShift          float32    // depthshift coefficient
	UsingDepthShift     bool       // using depthshift
	DepthShiftDirection bool       // controls the direction of the depthshift

	Frames      int       // frames per second? not sure why this is different from ActualFPS but use ActualFPS
	SecondTimer time.Time // seconds per frame
	ActualFPS   float32   // calculated FPS

	Framebuffer *ebiten.Image // image destination
	ScreenX     int           // width of the screen
	ScreenY     int           // height of the screen
	Font        *Font         // global font

	ChunkSize  int // size of the chunk, for generation
	ChunkDepth int // depth of the chunk, for generation
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
		UsingDepthShift:    true,
	}

	// init ebiten
	initialWidth, initialHeight := 1280, 720
	ebiten.SetTPS(60)
	// SetVsyncEnabled(false) is only for debug purposes
	// It is greatly recommended to use SetVsyncEnabled(true) in production
	ebiten.SetVsyncEnabled(!game.DebugMode)
	ebiten.SetWindowSize(initialWidth, initialHeight)
	ebiten.SetWindowTitle("ISOMETRICA Infdev")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// init render
	initRender()

	// init player, needs to happen before world is loaded
	game.Player = Player{
		Position: Vec3{0, 0, float32(game.World.SurfaceFeaturesBeginAt) + 10},
		Velocity: Vec3{0, 0, 0},
		Drag:     Vec3{.9, .9, .9},
		Texture:  "Default",
	}

	// check if there is a save file at the save path
	savePath := "save/demo"
	err := game.LoadGame(savePath)
	if err != nil {
		log.Printf("Failed to load world: %v", err)
		log.Println("Creating new world ...")

		// create new world

		game.World = World{}
		game.World.Initialize(rand.Int64())
		game.World.SavePath = filepath.Join("save", "demo")
		game.MakeEmptySave()
	} else {
		log.Println("Loading world from file.")
	}

	// run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
