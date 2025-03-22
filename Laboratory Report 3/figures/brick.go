package figures

import (
	"github.com/go-gl/gl/v2.1/gl"
	"math"
)

type Brick struct {
	vertices []float32
	indices  []int32
	scale    float32
}

func NewBrick(scale float32) *Brick {
	b := &Brick{
		scale: scale,
	}

	// Вершины шлакоблока
	rawVertices := []float32{
		// Нижняя грань
		0, 0, 0,
		0, 4, 0,
		4, 4, 0,
		4, 0, 0,

		// Внутренние отверстия
		1, 0, 1,
		1, 4, 1,
		3, 4, 1,
		3, 0, 1,

		1, 0, 3,
		1, 4, 3,
		3, 4, 3,
		3, 0, 3,

		1, 0, 4,
		1, 4, 4,
		3, 4, 4,
		3, 0, 4,

		1, 0, 6,
		1, 4, 6,
		3, 4, 6,
		3, 0, 6,

		// Задняя грань
		0, 0, 7,
		0, 4, 7,
		4, 4, 7,
		4, 0, 7,

		// Верхние грани
		1, 4, 0,
		1, 4, 7,
		3, 4, 7,
		3, 4, 0,

		1, 0, 0,
		1, 0, 7,
		3, 0, 7,
		3, 0, 0,
	}

	// Индексы треугольников
	rawIndices := []int32{
		0, 1, 2, 0, 2, 3,
		2, 22, 23, 2, 23, 3,
		0, 1, 21, 0, 21, 20,
		20, 21, 22, 20, 22, 23,
		4, 5, 6, 4, 6, 7,
		8, 9, 10, 8, 10, 11,
		4, 5, 9, 4, 9, 8,
		7, 6, 10, 7, 10, 11,
		12, 13, 14, 12, 14, 15,
		16, 17, 18, 16, 18, 19,
		12, 13, 17, 12, 17, 16,
		15, 14, 18, 15, 18, 19,
		1, 21, 24, 21, 24, 25,
		17, 18, 25, 18, 25, 26,
		9, 10, 13, 10, 13, 14,
		5, 6, 24, 6, 24, 27,
		2, 22, 26, 2, 26, 27,
		0, 20, 28, 20, 28, 29,
		16, 19, 29, 19, 29, 30,
		8, 11, 12, 11, 12, 15,
		4, 7, 28, 7, 28, 31,
		3, 23, 30, 3, 30, 31,
	}

	// Масштабируем координаты для OpenGL
	b.vertices = make([]float32, len(rawVertices))
	for i := 0; i < len(rawVertices); i += 3 {
		b.vertices[i] = (rawVertices[i] - 2) * b.scale
		b.vertices[i+1] = (rawVertices[i+1] - 2) * b.scale
		b.vertices[i+2] = (rawVertices[i+2] - 3.5) * b.scale
	}

	b.indices = rawIndices
	return b
}

func (b *Brick) Draw() {
	gl.Color3f(1.0, 1.0, 0.5)
	//gl.Color3f(rand.Float32(), rand.Float32(), rand.Float32())

	gl.Begin(gl.TRIANGLES)
	for _, idx := range b.indices {
		//gl.Color3f(rand.Float32(), rand.Float32(), rand.Float32())

		x := b.vertices[idx*3]
		y := b.vertices[idx*3+1]
		z := b.vertices[idx*3+2]

		// Расчет нормали
		length := float32(math.Sqrt(float64(x*x + y*y + z*z)))
		if length > 0 {
			gl.Normal3f(x/length, y/length, z/length)
		}

		gl.Vertex3f(x, y, z)
	}
	gl.End()
}
