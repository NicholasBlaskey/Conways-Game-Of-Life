package main

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

	"github.com/nicholasblaskey/Conways-Game-Of-Life/conways"
	"github.com/nicholasblaskey/Conways-Game-Of-Life/glfwBoilerplate"
)

func init() {
	runtime.LockOSThread()
}

func makeBuffers(board []float32,
	numX, numY int) ([]mgl.Vec2, uint32, uint32, uint32) {

	translations := conways.GetPositions(numX, numY)

	// Store these positions in a buffer
	sizeOfVec2 := 4 * 2
	var instanceVBO uint32
	gl.GenBuffers(1, &instanceVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, instanceVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(translations)*sizeOfVec2,
		unsafe.Pointer(&translations[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// Set up vertex data and buffers and config vertex attribs
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
	var quadVBO, quadVAO uint32
	gl.GenVertexArrays(1, &quadVAO)
	gl.GenBuffers(1, &quadVBO)
	gl.BindVertexArray(quadVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, quadVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(Vertices)*4,
		gl.Ptr(Vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	// Also set instance data
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, instanceVBO)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, int32(sizeOfVec2),
		gl.PtrOffset(0))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.VertexAttribDivisor(1, 1)

	return translations, quadVAO, quadVBO, instanceVBO
}

func main() {
	numX := 1000
	numY := 1000
	title := "One Color Instancing method"
	fmt.Println("Starting")

	window := glfwBoilerplate.InitGLFW(title,
		800, 600, false)
	defer glfw.Terminate()

	ourShader := shader.MakeShaders("oneColorInstancing.vs",
		"oneColorInstancing.fs")

	board := conways.CreateBoard(192921, numX, numY)
	translations, quadVAO, quadVBO, instanceVBO := makeBuffers(board, numX, numY)
	defer gl.DeleteVertexArrays(1, &quadVAO)
	defer gl.DeleteVertexArrays(1, &quadVBO)
	defer gl.DeleteVertexArrays(1, &instanceVBO)

	lastTime := 0.0
	numFrames := 0.0
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	for !window.ShouldClose() {
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)

		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		// Update board and VBO
		conways.UpdateBoard(board, numX, numY)
		positions := []mgl.Vec2{}
		for i := 0; i < len(board); i++ {
			if board[i] != 0 {
				positions = append(positions, translations[i])
			}
		}
		gl.BindBuffer(gl.ARRAY_BUFFER, instanceVBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(positions)*4*2,
			unsafe.Pointer(&positions[0]), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		// Render cubes
		ourShader.Use()
		gl.BindVertexArray(quadVAO)
		gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, int32(len(positions)))
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
