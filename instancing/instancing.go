// Translated from
// https://github.com/JoeyDeVries/LearnOpenGL/blob/master/src/2.lighting/1.colors/colors.cpp

package main

import(
	"runtime"
	//"log"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	
	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

	"github.com/nicholasblaskey/Conways-Game-Of-Life/glfwBoilerplate"
)

func init() {
	runtime.LockOSThread()
}

func makeBuffers(num_x, num_y int) (uint32, uint32, uint32) {
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
	
	// Set up vertex data and buffers and config vertex attribs
	Vertices := []float32{
		// positions     // colors
        -xOffset,  yOffset,  1.0, 0.0, 0.0,
		 xOffset, -yOffset,  0.0, 1.0, 0.0,
        -xOffset, -yOffset,  0.0, 0.0, 1.0,

        -xOffset,  yOffset,  1.0, 0.0, 0.0,
		 xOffset, -yOffset,  0.0, 1.0, 0.0,
		 xOffset,  yOffset,  0.0, 1.0, 1.0,
	}
	var quadVBO, quadVAO uint32		
	gl.GenVertexArrays(1, &quadVAO)
	gl.GenBuffers(1, &quadVBO)
	gl.BindVertexArray(quadVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, quadVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(Vertices) * 4,
		gl.Ptr(Vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)	
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 5 * 4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)	
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 5 * 4,
		gl.PtrOffset(sizeOfVec2))
	// Also set instance data
	gl.EnableVertexAttribArray(2)
	gl.BindBuffer(gl.ARRAY_BUFFER, instanceVBO)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 2 * 4, gl.PtrOffset(0))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.VertexAttribDivisor(2, 1)

	return quadVAO, quadVBO, instanceVBO
}

func main() {
	// 4 and 2 there is a bug with the square in the middle
	num_x := 4
	num_y := 3
	
	window := glfwBoilerplate.InitGLFW("Naive method fps = ",
		800, 600, false)
	defer glfw.Terminate()
	
	ourShader := shader.MakeShaders("instancing.vs", "instancing.fs")
	
	quadVAO, quadVBO, instanceVBO := makeBuffers(num_x, num_y)
	defer gl.DeleteVertexArrays(1, &quadVAO)
	defer gl.DeleteVertexArrays(1, &quadVBO)
	defer gl.DeleteVertexArrays(1, &instanceVBO)
	
	// Program loop
	for !window.ShouldClose() {
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		// Draw
		ourShader.Use()
		gl.BindVertexArray(quadVAO)
		gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, int32(num_x * num_y))
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

