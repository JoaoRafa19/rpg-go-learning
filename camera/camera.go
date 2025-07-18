package camera

import "math"

type Camera struct {
	X, Y float64
}

func NewCamera(x, y float64) *Camera {
	return &Camera{x, y}
}

func (c *Camera) FollowTarget(targetX, targetY float64, screenW, screenH float64) {
	c.X = -targetX + screenW/2
	c.Y = -targetY + screenH/2
}

func (c *Camera) Constrain(tileMapWidthPixels, tileMapHeightPixels, screenWidth, screenHeight float64) {
	c.X = math.Min(c.X, 0.0)
	c.Y = math.Min(c.Y, 0.0)

	c.X = math.Max(c.X, screenWidth-tileMapWidthPixels)
	c.Y = math.Max(c.Y, screenHeight-tileMapHeightPixels)
}
