// Translated from
// https://github.com/JoeyDeVries/LearnOpenGL/blob/master/src/2.lighting/1.colors/colors.cpp

package main

import(
	"runtime"
	"unsafe"
	"fmt"
//	"time"
	
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	
	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

	"github.com/nicholasblaskey/Conways-Game-Of-Life/glfwBoilerplate"
)

func init() {
	runtime.LockOSThread()
}

func makeBuffers(board []float32,
	num_x, num_y int) (uint32, uint32, uint32, uint32) {
	// Generate list of positions for our squares
	translations := []mgl32.Vec2{}
	xOffset := 1.0 / float32(num_x)
	yOffset := 1.0 / float32(num_y)
	for y := -num_y; y < num_y; y += 2 {
		for x := -num_x; x < num_x; x += 2 {
			translations = append(translations,
				mgl32.Vec2{float32(x) / float32(num_x) + xOffset,
					float32(y) / float32(num_y) + yOffset})
		}
	}

	// Store these positions in a buffer
	sizeOfVec2 := 4 * 2
	var instanceVBO uint32
	gl.GenBuffers(1, &instanceVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, instanceVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(translations) * sizeOfVec2,
		unsafe.Pointer(&translations[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// Bind the board to the VBO
	var colorVBO uint32
	gl.GenBuffers(1, &colorVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(board) * 4,
		unsafe.Pointer(&board[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	
	
	// Set up vertex data and buffers and config vertex attribs
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
	// Set the color buffer
	gl.EnableVertexAttribArray(2)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorVBO)
	gl.VertexAttribPointer(2, 1, gl.FLOAT, false, 4,
		gl.PtrOffset(0))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.VertexAttribDivisor(2, 1)
	
	
	return quadVAO, quadVBO, colorVBO, instanceVBO
}

func main() {
	num_x := 5
	num_y := 5
	title := "Instancing method"
	
	window := glfwBoilerplate.InitGLFW(title,
		800, 600, false)
	defer glfw.Terminate()
	
	ourShader := shader.MakeShaders("instancing.vs", "instancing.fs")

	board := []float32{}
	for i := 0; i < num_x * num_y; i++ {
		board = append(board, 1)
	}
	board[0] = 0
		
	quadVAO, quadVBO, colorVBO, instanceVBO := makeBuffers(board, num_x, num_y)
	defer gl.DeleteVertexArrays(1, &quadVAO)
	defer gl.DeleteVertexArrays(1, &quadVBO)
	defer gl.DeleteVertexArrays(1, &colorVBO)
	defer gl.DeleteVertexArrays(1, &instanceVBO)
	
	// Program loop
	lastTime := 0.0
	numFrames := 0.0

	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	for !window.ShouldClose() {		
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)	
		
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		
		if board[0] == 0 {
			board[0] = 1
		} else {
			board[0] = 0
		}
		fmt.Println(board[0])
		gl.BindBuffer(gl.ARRAY_BUFFER, colorVBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(board) * 4,
			unsafe.Pointer(&board[0]), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		
		// Draw
		ourShader.Use()
		gl.BindVertexArray(quadVAO)
		gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, int32(num_x * num_y))
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

