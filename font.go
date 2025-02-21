package main

import (
	"image"
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Font struct {
	Atlas      *ebiten.Image
	GridWidth  int
	CharSize   [2]int // does NOT include the padding which should always be 1,1
	CharSet    string // just a string of characters but it it's easy to work with like a list
	CharWidths []int  // the width of each character. indices match CharSet
	NewlinePad int
	LetterPad  int
}

// load a glyph from a font
func (font *Font) loadGlyph(char string) (subImage *ebiten.Image) {
	// get the character position in the atlas
	index := strings.Index(font.CharSet, char)
	y := (index / font.GridWidth) * (font.CharSize[1] + 1) // +1 for padding
	x := (index % font.GridWidth) * (font.CharSize[0] + 1) // +1 for padding

	// generate the subImage
	subImage, ok := font.Atlas.SubImage(image.Rectangle{Min: image.Point{x, y}, Max: image.Point{x + font.CharWidths[index], y + font.CharSize[1]}}).(*ebiten.Image)
	if !ok {
		log.Fatal("Failed to load glyph")
	}
	return
}

// render a string
func (font *Font) renderString(screen *ebiten.Image, text string, x, y int, color color.Color) {
	var (
		xx = x
		yy = y
	)

	for _, char := range text {
		if char == '\n' {
			xx = x
			yy += font.CharSize[1] + font.NewlinePad
			continue
		}
		glyph := font.loadGlyph(string(char))
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(xx), float64(yy))
		op.ColorM.ScaleWithColor(color)
		screen.DrawImage(glyph, op)
		xx += font.CharWidths[strings.Index(font.CharSet, string(char))] + font.LetterPad
	}
}

// create a new *Font
func newFont(path string, gridWidth int, charSet string, charWidths []int) (font *Font) {
	// load atlas image
	atlasImage, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// create font
	font = &Font{
		GridWidth:  gridWidth,
		CharSize:   [2]int{5, 11}, // most characters fit withing 5x7, but some like 'g' and 'y' go below the baseline, y7
		Atlas:      atlasImage,
		CharSet:    charSet,
		CharWidths: charWidths,
		LetterPad:  1,
		NewlinePad: 5,
	}

	return
}

// load the isometrica font
var isometricaFont = newFont("assets/font_atlas.png", 10, " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]\\^_`abcdefghijklmnopqrstuvwxyz{}|~", []int{
	3, 1, 3, 5, 5, 5, 5, 1, 2, 2, 3, 3, 1, 3, 1, 5, 5, 3, 5, 5, 5, 5, 5, 5, 5, 5, 1, 1, 3, 3, 3, 4, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 4, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 2, 2, 5, 3, 3, 2, 5, 5, 4, 5, 5, 5, 5, 5, 1, 4, 4, 3, 5, 4, 4, 5, 5, 4, 4, 4, 5, 3, 5, 3, 4, 4, 3, 3, 1, 4,
})
