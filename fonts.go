package main

import "path/filepath"

var FontLibrary map[string]FontData = map[string]FontData{
	"Isometrica": {
		AtlasPath:  filepath.Join("assets", "font_atlas.png"),
		GridWidth:  10,
		CharSize:   [2]int{5, 11},
		CharSet:    " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]\\^_`abcdefghijklmnopqrstuvwxyz{}|~",
		CharWidths: []int{3, 1, 3, 5, 5, 5, 5, 1, 2, 2, 3, 3, 1, 3, 1, 5, 5, 3, 5, 5, 5, 5, 5, 5, 5, 5, 1, 1, 3, 3, 3, 4, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 4, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 2, 2, 5, 3, 3, 2, 5, 5, 4, 5, 5, 5, 5, 5, 1, 4, 4, 3, 5, 4, 4, 5, 5, 4, 4, 4, 5, 3, 5, 3, 4, 4, 3, 3, 1, 4},
	},
}
