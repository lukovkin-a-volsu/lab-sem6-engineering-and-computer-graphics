package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	glut "github.com/vbsw/freeglut"
	"lab/figures"
	"math"
)

const (
	WindowWidth      = 640
	WindowHeight     = 480
	WindowTitle      = "Laboratory Report 3"
	RotationStep     = 5.0
	MouseSensitivity = 0.5
	MovementSpeed    = 0.1
)

var (
	angleX, angleY, angleZ float32
	isDragging             bool
	xPos                   float32 = 0.0
	yPos                   float32 = 0.0
	zPos                   float32 = -5.0
	lastMouseX             int
	lastMouseY             int
	fov                    float32 = 45.0
	wPressed               bool
	aPressed               bool
	sPressed               bool
	dPressed               bool
	brick                  *figures.Brick
	cube                   *figures.Cube
)

func display() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity()

	gl.Translatef(xPos, yPos, zPos)
	gl.Rotatef(angleX, 1.0, 0.0, 0.0)
	gl.Rotatef(angleY, 0.0, 1.0, 0.0)
	gl.Rotatef(angleZ, 0.0, 0.0, 1.0)

	gl.PushMatrix()
	gl.Translatef(-2.5, 0, 0)
	cube.Draw()
	gl.PopMatrix()

	gl.PushMatrix()
	gl.Translatef(2.5, 0, 0)
	brick.Draw()
	gl.PopMatrix()

	glut.SwapBuffers()
}

func initGL() {
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.NORMALIZE)
	gl.Enable(gl.LIGHT1)
	gl.Enable(gl.COLOR_MATERIAL)
	gl.ColorMaterial(gl.FRONT_AND_BACK, gl.AMBIENT_AND_DIFFUSE)

	lightAmbient := []float32{0.2, 0.2, 0.2, 1.0}
	lightDiffuse := []float32{0.2, 0.2, 0.2, 1.0}
	lightSpecular := []float32{1.0, 1.0, 1.0, 1.0}
	gl.Lightfv(gl.LIGHT1, gl.AMBIENT, &lightAmbient[0])
	gl.Lightfv(gl.LIGHT1, gl.DIFFUSE, &lightDiffuse[0])
	gl.Lightfv(gl.LIGHT1, gl.SPECULAR, &lightSpecular[0])

	// Настройки материала
	materialSpecular := []float32{0.5, 0.5, 0.5, 1.0}
	materialShininess := []float32{50.0}
	gl.Materialfv(gl.FRONT, gl.SPECULAR, &materialSpecular[0])
	gl.Materialfv(gl.FRONT, gl.SHININESS, &materialShininess[0])

	gl.ClearColor(0.25, 0.25, 0.35, 1.0)
}

func reshape(width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	aspect := float32(width) / float32(height)
	perspective(fov, aspect, 0.1, 100.0)
	gl.MatrixMode(gl.MODELVIEW)
}

func specialKeyboard(key int, x, y int) {
	switch key {
	case glut.KEY_LEFT:
		angleY -= RotationStep
	case glut.KEY_RIGHT:
		angleY += RotationStep
	case glut.KEY_UP:
		angleX -= RotationStep
	case glut.KEY_DOWN:
		angleX += RotationStep
	}
	glut.PostRedisplay()
}

func mouse(button, state, x, y int) {
	if button == glut.LEFT_BUTTON {
		if state == glut.DOWN {
			isDragging = true
			lastMouseX = x
			lastMouseY = y
		} else {
			isDragging = false
		}
	}
}

func motion(x, y int) {
	if isDragging {
		deltaX := x - lastMouseX
		deltaY := y - lastMouseY
		angleY += float32(deltaX) * MouseSensitivity
		angleX += float32(deltaY) * MouseSensitivity
		lastMouseX = x
		lastMouseY = y
		glut.PostRedisplay()
	}
}

func perspective(fov, aspect, near, far float32) {
	fovRad := fov * (math.Pi / 180.0)
	tanHalfFov := float32(math.Tan(float64(fovRad / 2.0)))

	gl.Frustum(
		float64(-near*tanHalfFov*aspect),
		float64(near*tanHalfFov*aspect),
		float64(-near*tanHalfFov),
		float64(near*tanHalfFov),
		float64(near),
		float64(far),
	)
}

func keyboard(key uint8, x, y int) {
	switch key {
	case 'w', 'W':
		wPressed = true
	case 's', 'S':
		sPressed = true
	case 'a', 'A':
		aPressed = true
	case 'd', 'D':
		dPressed = true
	}
}

func keyboardUp(key uint8, x, y int) {
	switch key {
	case 'w', 'W':
		wPressed = false
	case 's', 'S':
		sPressed = false
	case 'a', 'A':
		aPressed = false
	case 'd', 'D':
		dPressed = false
	}
}

func update() {
	if wPressed {
		zPos += MovementSpeed
	}
	if sPressed {
		zPos -= MovementSpeed
	}
	if aPressed {
		xPos += MovementSpeed
	}
	if dPressed {
		xPos -= MovementSpeed
	}
	glut.PostRedisplay()
}

func init() {
	glut.Init()
}

func timer(value int) {
	update()
	glut.TimerFunc(16, timer, 0)
}

func main() {
	glut.InitDisplayMode(glut.DOUBLE | glut.RGBA | glut.DEPTH)
	glut.InitWindowSize(WindowWidth, WindowHeight)
	glut.CreateWindow(WindowTitle)

	if err := gl.Init(); err != nil {
		if err := gl.GetError(); err != gl.NO_ERROR {
			fmt.Println("OpenGL error:", err)
		}
		panic(err)
	}

	initGL()

	glut.ReshapeFunc(reshape)
	glut.DisplayFunc(display)
	glut.SpecialFunc(specialKeyboard)
	glut.MotionFunc(motion)
	glut.MouseFunc(mouse)

	glut.KeyboardFunc(keyboard)
	glut.KeyboardUpFunc(keyboardUp)
	glut.TimerFunc(0, timer, 0)

	cube = figures.NewCube(2.0)
	brick = figures.NewBrick(0.5)

	glut.MainLoop()
}

// 6
// gl.BindBuffer()
// gl.BufferData()
// gl.glUnmapBuffer()
// gl.DrawArrays()
// gl.DrawElements()
// gl.Vertex()
// gl.Color()
// glu.Sphere
// glu.Cylinder

// 7
// gl.MatrixMode: Выбирает матрицу для манипуляций (модель, проекция, текстура).
// gl.LoadIdentity: Сбрасывает матрицу в единичную.
// gl.LoadMatrixf: Загружает пользовательскую матрицу.
// gl.MultMatrixf: Умножает текущую матрицу на заданную.
// gl.Translatef: Перемещает матрицу.
// gl.Rotatef: Вращает матрицу.
// gl.Scalef: Масштабирует объект.

// 8
// Функции для работы с проекциями и областью вывода:
// gl.Viewport: Устанавливает область вывода на экране.
// gl.Ortho: Устанавливает ортогональную проекцию.
// gluPerspective: Устанавливает перспективную проекцию (используется библиотека GLU)
// gl.Frustum: Устанавливает обрезающую проекцию.
// gluLookAt: Устанавливает точку зрения камеры (используется библиотека GLU)
