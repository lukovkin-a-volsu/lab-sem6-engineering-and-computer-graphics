package main

import (
	gl "github.com/go-gl/gl/v2.1/gl"
	glut "github.com/vbsw/freeglut"
	"math"
)

const (
	WindowWidth      = 640
	WindowHeight     = 480
	WindowTitle      = "Laboratory Report 3"
	RotationStep     = 5.0
	MouseSensitivity = 0.5
	InitialFOV       = 45.0
	PerspectiveNear  = 1.0
	PerspectiveFar   = 100.0
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
)

func drawCube() {
	gl.Begin(gl.QUADS)

	// Лицевая грань
	gl.Color3f(1.0, 0.0, 0.0)
	gl.Vertex3f(-1.0, -1.0, 1.0)
	gl.Vertex3f(1.0, -1.0, 1.0)
	gl.Vertex3f(1.0, 1.0, 1.0)
	gl.Vertex3f(-1.0, 1.0, 1.0)

	// Задняя грань
	gl.Color3f(0.0, 1.0, 0.0)
	gl.Vertex3f(-1.0, -1.0, -1.0)
	gl.Vertex3f(-1.0, 1.0, -1.0)
	gl.Vertex3f(1.0, 1.0, -1.0)
	gl.Vertex3f(1.0, -1.0, -1.0)

	// Верхняя грань
	gl.Color3f(0.0, 0.0, 1.0)
	gl.Vertex3f(-1.0, 1.0, -1.0)
	gl.Vertex3f(-1.0, 1.0, 1.0)
	gl.Vertex3f(1.0, 1.0, 1.0)
	gl.Vertex3f(1.0, 1.0, -1.0)

	// Нижняя грань
	gl.Color3f(1.0, 1.0, 0.0)
	gl.Vertex3f(-1.0, -1.0, -1.0)
	gl.Vertex3f(1.0, -1.0, -1.0)
	gl.Vertex3f(1.0, -1.0, 1.0)
	gl.Vertex3f(-1.0, -1.0, 1.0)

	// Правая грань
	gl.Color3f(1.0, 0.0, 1.0)
	gl.Vertex3f(1.0, -1.0, -1.0)
	gl.Vertex3f(1.0, 1.0, -1.0)
	gl.Vertex3f(1.0, 1.0, 1.0)
	gl.Vertex3f(1.0, -1.0, 1.0)

	// Левая грань
	gl.Color3f(0.0, 1.0, 1.0)
	gl.Vertex3f(-1.0, -1.0, -1.0)
	gl.Vertex3f(-1.0, -1.0, 1.0)
	gl.Vertex3f(-1.0, 1.0, 1.0)
	gl.Vertex3f(-1.0, 1.0, -1.0)
	gl.End()
}

func display() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity()

	// Применяем преобразования камеры
	gl.Translatef(xPos, yPos, zPos)
	gl.Rotatef(angleX, 1.0, 0.0, 0.0)
	gl.Rotatef(angleY, 0.0, 1.0, 0.0)
	gl.Rotatef(angleZ, 0.0, 0.0, 1.0)

	drawCube()
	glut.SwapBuffers()
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

func initGL() {
	gl.Enable(gl.DEPTH_TEST)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	aspect := float32(glut.Get(glut.WINDOW_WIDTH)) / float32(glut.Get(glut.WINDOW_HEIGHT))
	perspective(InitialFOV, aspect, PerspectiveNear, PerspectiveFar)
	gl.MatrixMode(gl.MODELVIEW)
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

// Добавляем обработчики клавиатуры
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
		xPos -= MovementSpeed
	}
	if dPressed {
		xPos += MovementSpeed
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

	// Критически важная инициализация
	if err := gl.Init(); err != nil {
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

	glut.MainLoop()
}
