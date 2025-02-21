package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// texture atlas
var groundTextureAtlas, _, _ = ebitenutil.NewImageFromFile("assets/block_atlas.png")

// default voxel dictionary/lookup table
var defaultVoxelDictionary = VoxelDictionary{
	Voxels: []Voxel{
		{Name: "Air", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{96, 96}, Max: image.Point{128, 128}}},
		{Name: "Grass", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{32, 0}, Max: image.Point{64, 32}}},
		{Name: "Water", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{64, 0}, Max: image.Point{96, 32}}},
		{Name: "Sand", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{32, 32}}},
		{Name: "Stone", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{96, 0}, Max: image.Point{128, 32}}},
		{Name: "Dirt", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{0, 32}, Max: image.Point{32, 64}}},
		{Name: "Wood", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{32, 32}, Max: image.Point{64, 64}}},
		{Name: "Leaves", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{64, 32}, Max: image.Point{96, 64}}},
		{Name: "Flower", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{96, 32}, Max: image.Point{128, 64}}},
		{Name: "Tall_Grass", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{0, 64}, Max: image.Point{32, 96}}},
		{Name: "Cobblestone", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{32, 64}, Max: image.Point{64, 96}}},
		{Name: "Snowy_Grass", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{64, 64}, Max: image.Point{96, 96}}},
		{Name: "Snowy_Leaves", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{96, 64}, Max: image.Point{128, 96}}},
		{Name: "Snowy_Tall_Grass", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{0, 96}, Max: image.Point{32, 128}}},
		{Name: "Snowy_Flower", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{32, 96}, Max: image.Point{64, 128}}},
	},
	Transparent:          []string{"Air", "Water", "Flower"},
	Opaque:               []string{"Grass", "Sand", "Stone", "Dirt", "Wood", "Leaves"},
	TransparentNoCulling: []string{"Flower", "Tall_Grass"},
}

// voxel to be used when an error occurs
var errorVoxelDictionary = VoxelDictionary{
	Voxels: []Voxel{
		{Name: "Error", Atlas: groundTextureAtlas, TextureRect: image.Rectangle{Min: image.Point{64, 96}, Max: image.Point{96, 128}}},
	},
}
