package main

// Font Engine vN25
// Adapted from another of my projects, which adapted it from this project earlier on

import (
	"image"
	"image/color"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Font struct {
	Atlas      *ebiten.Image
	GridWidth  int
	CharSize   [2]int // Width and height of each character cell (excluding 1px padding)
	CharSet    string // String containing all supported characters in order matching atlas
	CharWidths []int  // Width of each character (indices match CharSet)
	NewlinePad int    // Extra vertical padding between lines
	LetterPad  int    // Extra horizontal padding between characters

	GlyphCache map[rune]*ebiten.Image
}

// loadGlyph extracts a single character image from the font atlas
func (f *Font) loadGlyph(char rune) (subImage *ebiten.Image) {

	if cached, ok := f.GlyphCache[char]; ok {
		return cached
	}

	// Find character position in atlas grid
	byteIndex := strings.IndexRune(f.CharSet, char)
	if byteIndex == -1 {
		if strings.ContainsRune(f.CharSet, '�') {
			return f.loadGlyph('�')
		}

		// Debug: print the rune values
		log.Printf("Looking for rune: %q (value: %d)", char, char)
		log.Printf("Charset runes:")
		for i, r := range f.CharSet {
			log.Printf("  [%d] %q (value: %d)", i, r, r)
		}
		log.Panicf("Couldn't find %q in Charset!", char)
	}

	index := utf8.RuneCountInString(f.CharSet[:byteIndex])

	gridY := (index / f.GridWidth) * (f.CharSize[1] + 1) // Add 1px padding
	gridX := (index % f.GridWidth) * (f.CharSize[0] + 1) // Add 1px padding

	// Extract character image
	subImage, ok := f.Atlas.SubImage(image.Rectangle{
		Min: image.Point{gridX, gridY},
		Max: image.Point{gridX + f.getCharWidthIdx(index), gridY + f.CharSize[1]},
	}).(*ebiten.Image)
	if !ok {
		log.Fatal("Failed to extract character from font atlas")
	}
	f.GlyphCache[char] = subImage
	return
}

func (f *Font) getCharWidthRune(char rune) (width int) {
	return f.getCharWidthIdx(strings.IndexRune(f.CharSet, char))
}

func (f *Font) getCharWidthIdx(idx int) (width int) {
	// get width of char
	if len(f.CharWidths) == 1 {
		// monospace
		width = f.CharWidths[0]
	} else {
		width = f.CharWidths[idx]
	}
	return
}

// renderString draws text to the screen with specified position and color
func (f *Font) renderString(screen *ebiten.Image, text string, x, y int, color color.Color) {
	cursorX := x
	cursorY := y

	for _, char := range text {
		if char == rune(0x00) {
			continue
		}

		if char == '\n' {
			// Handle newlines
			cursorX = x
			cursorY += f.CharSize[1] + f.NewlinePad
			continue
		}

		// Draw character
		glyph := f.loadGlyph(char)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(cursorX), float64(cursorY))
		op.ColorM.ScaleWithColor(color)
		screen.DrawImage(glyph, op)

		// Advance cursor

		// update x
		cursorX += f.getCharWidthRune(char) + f.LetterPad
	}
}

// newFont creates a new Font from an atlas image file
func newFont(atlasPath string, gridWidth int, charSize [2]int, charSet string, charWidths []int) (font *Font) {
	// Load font atlas image
	atlasImage, _, err := ebitenutil.NewImageFromFile(atlasPath)
	if err != nil {
		log.Fatal("Failed to load font atlas:", err)
	}

	font = &Font{
		Atlas:      atlasImage,
		GridWidth:  gridWidth,
		CharSize:   charSize, // Most chars are 5x7, some extend below baseline to 11px
		CharSet:    charSet,
		CharWidths: charWidths,
		LetterPad:  1,
		NewlinePad: 5,
		GlyphCache: make(map[rune]*ebiten.Image),
	}

	return
}

// FontData stores metadata about a font for it to be easily manipulated and stored while using less memory.
type FontData struct {
	AtlasPath  string
	GridWidth  int
	CharSize   [2]int
	CharSet    string
	CharWidths []int
}

// binds FontData into the newFont() func
func (f FontData) newFont() (font *Font) {
	return newFont(
		f.AtlasPath,
		f.GridWidth,
		f.CharSize,
		f.CharSet,
		f.CharWidths,
	)
}
