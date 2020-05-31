// Translated from
// https://github.com/JoeyDeVries/LearnOpenGL/blob/master/src/2.lighting/1.colors/colors.cpp

package main

import(
	"runtime"
	"fmt"
	"unsafe"
//	"time"
	
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	
	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

	"github.com/nicholasblaskey/Conways-Game-Of-Life/glfwBoilerplate"	
	"github.com/nicholasblaskey/Conways-Game-Of-Life/conways"
)

func init() {
	runtime.LockOSThread()
}

func makeBuffers(board []float32,
	numX, numY int) (uint32, uint32, uint32) {

	translations := conways.GetPositions(numX, numY)
	
	// Store these positions in a buffer
	sizeOfVec2 := 4 * 2
	var instanceVBO uint32
	gl.GenBuffers(1, &instanceVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, instanceVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(translations) * sizeOfVec2,
		unsafe.Pointer(&translations[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	/*

	// Bind the board to the VBO
	var colorVBO uint32
	gl.GenBuffers(1, &colorVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(board) * 4,
		unsafe.Pointer(&board[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	*/
	
	// Set up vertex data and buffers and config vertex attribs
	xOffset := 1.0 / float32(numX)
	yOffset := 1.0 / float32(numY)
	Vertices := []float32{
		// positions     
		-xOffset,  yOffset, 
		xOffset, -yOffset,
		-xOffset, -yOffset, 

		-xOffset,  yOffset, 
		xOffset, -yOffset, 
		xOffset,  yOffset, 
	}
	var quadVBO, quadVAO uint32		
	gl.GenVertexArrays(1, &quadVAO)
	gl.GenBuffers(1, &quadVBO)
	gl.BindVertexArray(quadVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, quadVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(Vertices) * 4,
		gl.Ptr(Vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)	
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2 * 4, gl.PtrOffset(0))
	// Also set instance data
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, instanceVBO)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, int32(sizeOfVec2),
		gl.PtrOffset(0))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.VertexAttribDivisor(1, 1)	

	/*
	// Set the color buffer
	gl.EnableVertexAttribArray(2)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorVBO)
	gl.VertexAttribPointer(2, 1, gl.FLOAT, false, 4,
		gl.PtrOffset(0))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.VertexAttribDivisor(2, 1)
*/	
	
	return quadVAO, quadVBO, instanceVBO
}

func main() {
	numX := 1000
	numY := 1000
	title := "Instancing method"
	fmt.Println("Starting")
	
	window := glfwBoilerplate.InitGLFW(title,
		800, 600, false)
	defer glfw.Terminate()
	
	ourShader := shader.MakeShaders("oneColorInstancing.vs",
		"oneColorInstancing.fs")
	
	board := conways.CreateBoard(192921, numX, numY)

	quadVAO, quadVBO, instanceVBO := makeBuffers(board, numX, numY)
	defer gl.DeleteVertexArrays(1, &quadVAO)
	defer gl.DeleteVertexArrays(1, &quadVBO)
	defer gl.DeleteVertexArrays(1, &instanceVBO)
	
	lastTime := 0.0
	numFrames := 0.0
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	for !window.ShouldClose() {
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)	
		
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		//time.Sleep(1 * time.Millisecond)

		// Update board and VBO
		conways.UpdateBoard(board, numX, numY)
		gl.BindBuffer(gl.ARRAY_BUFFER, instanceVBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(board) * 4,
			unsafe.Pointer(&board[0]), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		// Render cubes
		ourShader.Use()
		gl.BindVertexArray(quadVAO)
		gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, int32(numX * numY))
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

