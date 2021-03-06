package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

	"github.com/nicholasblaskey/Conways-Game-Of-Life/conways"
	"github.com/nicholasblaskey/Conways-Game-Of-Life/glfwBoilerplate"
)

func init() {
	runtime.LockOSThread()
}

func makeBuffers(board []float32, numX, numY int) (uint32, uint32) {
	xOffset := 1.0 / float32(numX)
	yOffset := 1.0 / float32(numY)
	Vertices := []float32{
		// positions
		-xOffset, yOffset,
		xOffset, -yOffset,
		-xOffset, -yOffset,

		-xOffset, yOffset,
		xOffset, -yOffset,
		xOffset, yOffset,
	}
	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(Vertices)*4,
		gl.Ptr(Vertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	return VAO, VBO
}

func main() {
	numX := 1000
	numY := 1000
	title := "One color method"
	fmt.Println("Starting")

	window := glfwBoilerplate.InitGLFW(title,
		800, 600, false)
	defer glfw.Terminate()

	ourShader := shader.MakeShaders("oneColor.vs", "oneColor.fs")
	translations := conways.GetPositions(numX, numY)
	board := conways.CreateBoard(192921, numX, numY)

	VAO, VBO := makeBuffers(board, numX, numY)
	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteVertexArrays(1, &VBO)

	lastTime := 0.0
	numFrames := 0.0
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	for !window.ShouldClose() {
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)

		// Draw black tiles in with the background
		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		// Update board and VBO
		conways.UpdateBoard(board, numX, numY)

		// Render cubes
		ourShader.Use()
		for i := 0; i < len(translations); i++ {
			if board[i] == 1.0 {
				ourShader.SetVec2("aOffset", translations[i])

				gl.BindVertexArray(VAO)
				gl.DrawArrays(gl.TRIANGLES, 0, 6)
				gl.BindVertexArray(0)
			}
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
