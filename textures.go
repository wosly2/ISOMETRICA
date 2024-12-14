package main

import (
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// texture atlas
var groundTextureAtlas, _, gErr = ebitenutil.NewImageFromFile("block_atlas.png")

// default voxel dictionary/lookup table
var defaultVoxelDictionary = VoxelDictionary{
	Voxels: []Voxel{
		{Name: "Air", Atlas: groundTextureAtlas, Texture: [4]int{96, 96, 128, 128}},
		{Name: "Grass", Atlas: groundTextureAtlas, Texture: [4]int{32, 0, 64, 32}},
		{Name: "Water", Atlas: groundTextureAtlas, Texture: [4]int{64, 0, 96, 32}},
		{Name: "Sand", Atlas: groundTextureAtlas, Texture: [4]int{0, 0, 32, 32}},
		{Name: "Stone", Atlas: groundTextureAtlas, Texture: [4]int{96, 0, 128, 32}},
		{Name: "Dirt", Atlas: groundTextureAtlas, Texture: [4]int{0, 32, 32, 64}},
		{Name: "Wood", Atlas: groundTextureAtlas, Texture: [4]int{32, 32, 64, 64}},
		{Name: "Leaves", Atlas: groundTextureAtlas, Texture: [4]int{64, 32, 96, 64}},
		{Name: "Flower", Atlas: groundTextureAtlas, Texture: [4]int{96, 32, 128, 64}},
		{Name: "Tall_Grass", Atlas: groundTextureAtlas, Texture: [4]int{0, 64, 32, 96}},
	},
	Transparent:          []string{"Air", "Water", "Flower"},
	Opaque:               []string{"Grass", "Sand", "Stone", "Dirt", "Wood", "Leaves"},
	TransparentNoCulling: []string{"Flower", "Tall_Grass"},
}

var errorVoxelDictionary = VoxelDictionary{
	Voxels: []Voxel{
		{Name: "Error", Atlas: groundTextureAtlas, Texture: [4]int{64, 96, 96, 128}},
	},
	Opaque: []string{"Error"},
}
