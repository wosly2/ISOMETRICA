package main

import "math"

type Vec2 struct {
	X, Y float32
}

type Vec2i struct {
	X, Y int
}

type Vec3 struct {
	X, Y, Z float32
}

func (v Vec3) Normalize() Vec3 {
	length := float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
	if length < 1e-6 {
		return Vec3{0, 0, 0}
	}
	return Vec3{
		X: v.X / length,
		Y: v.Y / length,
		Z: v.Z / length,
	}
}

// // rounds a number to the nearest multiple of `a`
// func roundNArepeat(n, a int) int {
// 	remainder := n % a
// 	if remainder < a/2 {
// 		return n - remainder
// 	}
// 	return n + (a - remainder)
// }

// // rounds a number to either 0 or `a` (no multiples)
// func roundNA(n, a int) int {
// 	if n >= a/2 {
// 		return a
// 	}
// 	return 0
// }

func pointsMakeCCWTurn(A, B, C Vec2i) bool {
	return (B.X-A.X)*(C.Y-A.Y) > (C.X-A.X)*(B.Y-A.Y)
}

func doLinesIntersect(A, B, C, D Vec2i) bool {
	return pointsMakeCCWTurn(A, C, D) != pointsMakeCCWTurn(B, C, D) && pointsMakeCCWTurn(A, B, C) != pointsMakeCCWTurn(A, B, D)
}

// checks if a point is in a rectangle defined by width and height
func isPointInBounds(x, y, width, height int) bool {
	return x >= 0 && y >= 0 && x < width && y < height
}

// checks if a line intersects a rectangle defined by width and height and starting point
func doesLineIntersectRectangle(x, y, x2, y2, rx, ry, rWidth, rHeight int) bool {
	// check if line points are in the rectangle
	if isPointInBounds(x, y, rWidth, rHeight) || isPointInBounds(x2, y2, rWidth, rHeight) {
		return true
	}

	// check intersections with the top and bottom edges
	if doLinesIntersect(Vec2i{x, y}, Vec2i{x2, y2}, Vec2i{rx, ry}, Vec2i{rx + rWidth, ry}) || doLinesIntersect(Vec2i{x, y}, Vec2i{x2, y2}, Vec2i{rx, ry + rHeight}, Vec2i{rx + rWidth, ry + rHeight}) {
		return true
	}
	// check intersections with the left and right edges
	if doLinesIntersect(Vec2i{x, y}, Vec2i{x2, y2}, Vec2i{rx, ry}, Vec2i{rx, ry + rHeight}) || doLinesIntersect(Vec2i{x, y}, Vec2i{x2, y2}, Vec2i{rx + rWidth, ry}, Vec2i{rx + rWidth, ry + rHeight}) {
		return true
	}

	return false
}

func absi(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}
