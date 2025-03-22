package figures

import (
	"github.com/go-gl/gl/v2.1/gl"
)

type Cube struct {
	size float32
}

func NewCube(size float32) *Cube {
	return &Cube{size: size}
}

func (c *Cube) Draw() {
	s := c.size / 2

	gl.Begin(gl.QUADS)

	// Лицевая грань
	gl.Normal3f(0, 0, 1)
	gl.Color3f(1.0, 0.0, 0.0)
	gl.Vertex3f(-s, -s, s)
	gl.Vertex3f(s, -s, s)
	gl.Vertex3f(s, s, s)
	gl.Vertex3f(-s, s, s)

	// Задняя грань
	gl.Normal3f(0, 0, -1)
	gl.Color3f(0.0, 1.0, 0.0)
	gl.Vertex3f(-s, -s, -s)
	gl.Vertex3f(-s, s, -s)
	gl.Vertex3f(s, s, -s)
	gl.Vertex3f(s, -s, -s)

	// Верхняя грань
	gl.Normal3f(0, 1, 0)
	gl.Color3f(0.0, 0.0, 1.0)
	gl.Vertex3f(-s, s, -s)
	gl.Vertex3f(-s, s, s)
	gl.Vertex3f(s, s, s)
	gl.Vertex3f(s, s, -s)

	// Нижняя грань
	gl.Normal3f(0, -1, 0)
	gl.Color3f(1.0, 1.0, 0.0)
	gl.Vertex3f(-s, -s, -s)
	gl.Vertex3f(s, -s, -s)
	gl.Vertex3f(s, -s, s)
	gl.Vertex3f(-s, -s, s)

	// Правая грань
	gl.Normal3f(1, 0, 0)
	gl.Color3f(1.0, 0.0, 1.0)
	gl.Vertex3f(s, -s, -s)
	gl.Vertex3f(s, s, -s)
	gl.Vertex3f(s, s, s)
	gl.Vertex3f(s, -s, s)

	// Левая грань
	gl.Normal3f(-1, 0, 0)
	gl.Color3f(0.0, 1.0, 1.0)
	gl.Vertex3f(-s, -s, -s)
	gl.Vertex3f(-s, -s, s)
	gl.Vertex3f(-s, s, s)
	gl.Vertex3f(-s, s, -s)

	gl.End()
}
