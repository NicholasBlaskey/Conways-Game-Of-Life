// Translated from
// https://github.com/JoeyDeVries/LearnOpenGL/blob/master/src/2.lighting/1.colors/colors.cpp

package main

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

	"github.com/nicholasblaskey/Conways-Game-Of-Life/conways"
	"github.com/nicholasblaskey/Conways-Game-Of-Life/glfwBoilerplate"
)

func init() {
	runtime.LockOSThread()
}

func makeBuffers(board []float32, numX, numY int) (uint32, uint32, uint32) {
	translations := conways.GetPositions(numX, numY)

	// Bind the board to the VBO
	var colorVBO uint32
	gl.GenBuffers(1, &colorVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(board)*4,
		unsafe.Pointer(&board[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(translations)*4*2,
		unsafe.Pointer(&translations[0]), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorVBO)
	gl.VertexAttribPointer(1, 1, gl.FLOAT, false, 4,
		gl.PtrOffset(0))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return VAO, VBO, colorVBO
}

func main() {
	numX := 1000
	numY := 1000
	title := "Geometry Shader method"
	fmt.Println("Starting")

	window := glfwBoilerplate.InitGLFW(title,
		800, 600, false)
	defer glfw.Terminate()

	ourShader := shader.MakeGeomShaders("geoShader.vs",
		"geoShader.fs", "geoShader.gs")
	board := conways.CreateBoard(192921, numX, numY)

	VAO, VBO, colorVBO := makeBuffers(board, numX, numY)
	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteVertexArrays(1, &VBO)
	defer gl.DeleteVertexArrays(1, &colorVBO)

	lastTime := 0.0
	numFrames := 0.0

	ourShader.Use()
	ourShader.SetFloat("xOffset", 1.0/float32(numX))
	ourShader.SetFloat("yOffset", 1.0/float32(numY))

	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	for !window.ShouldClose() {
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		// Update board and VBO
		conways.UpdateBoard(board, numX, numY)

		gl.BindBuffer(gl.ARRAY_BUFFER, colorVBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(board)*4,
			unsafe.Pointer(&board[0]), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		// Render cubes
		ourShader.Use()
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.POINTS, 0, int32(len(board)))
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
