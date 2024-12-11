package main

import (
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var groundTextureAtlas, _, err = ebitenutil.NewImageFromFile("block_atlas.png")

// default voxel dictionary/lookup table
var defaultVoxelDictionary = VoxelDictionary{
	Voxels: []Voxel{
		{Name: "Cake", Atlas: groundTextureAtlas, Texture: [4]int{0, 0, 32, 32}},
		{Name: "Grass", Atlas: groundTextureAtlas, Texture: [4]int{32, 0, 64, 32}},
		{Name: "Water", Atlas: groundTextureAtlas, Texture: [4]int{64, 0, 96, 32}},
	},
}

var demoChunk = MakeChunk([][][]VoxelPointer{
	// x
	{
		// y
		{
			// z
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 1},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
		},
		{
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 1},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
		},
		{
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 1},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
		},
	},
	{
		{
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 1},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
		},
		{
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 0},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 0},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 0},
		},
		{
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 1},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
		},
	},
	{
		{
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 1},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
		},
		{
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 1},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
		},
		{
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 1},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
			VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 2},
		},
	},
})
