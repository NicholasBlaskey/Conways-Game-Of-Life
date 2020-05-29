// Translated from
// https://github.com/JoeyDeVries/LearnOpenGL/blob/master/src/2.lighting/1.colors/colors.cpp

package main

import(
	"runtime"
	"fmt"
	
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	
	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

	"github.com/nicholasblaskey/Conways-Game-Of-Life/glfwBoilerplate"	
	"github.com/nicholasblaskey/Conways-Game-Of-Life/conways"
)

func init() {
	runtime.LockOSThread()
}

func getColorRow(board, colors []float32, offset, amount int) {
	endPoint := offset + amount
	for i := offset; i < endPoint; i++ {
		colors[i * 6] = board[i]
		colors[i * 6 + 1] = board[i]
		colors[i * 6 + 2] = board[i]
		colors[i * 6 + 3] = board[i]
		colors[i * 6 + 4] = board[i]
		colors[i * 6 + 5] = board[i]
	}	
}

func getColors(board []float32, numX, numY int) []float32 {
	colors := make([]float32, 6 * numX * numY)
	for i := 0; i < numY - 1; i++ {
		go getColorRow(board, colors, i * numY, numX)
	}
	getColorRow(board, colors, (numY - 1) * numY, numX)

	return colors
}

func makeBuffers(board []float32, numX, numY int) (uint32, uint32, uint32) {
	
	translations := conways.GetPositions(numX, numY)
	xOffset := 1.0 / float32(numX)
	yOffset := 1.0 / float32(numY)
	
	Vertices := []float32{}
	for i := 0; i < len(translations); i++ {
		Vertices = append(Vertices,
			// positions     
			-xOffset + translations[i][0],  yOffset + translations[i][1], 
			xOffset + translations[i][0], -yOffset + translations[i][1],
			-xOffset + translations[i][0], -yOffset + translations[i][1], 
			
			-xOffset + translations[i][0],  yOffset + translations[i][1], 
			xOffset + translations[i][0], -yOffset + translations[i][1], 
			xOffset + translations[i][0],  yOffset + translations[i][1],
		)
	}

	colors := getColors(board, numX, numY)
	var colorVBO uint32
	gl.GenBuffers(1, &colorVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(colors) * 4,
		gl.Ptr(colors), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	
	var VAO, VBO uint32		
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	
	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(Vertices) * 4,
		gl.Ptr(Vertices), gl.STATIC_DRAW)
	
	gl.EnableVertexAttribArray(0)	
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2 * 4, gl.PtrOffset(0))
	
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorVBO)
	gl.VertexAttribPointer(1, 1, gl.FLOAT, false, 4 * 1,
		gl.PtrOffset(0))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.VertexAttribDivisor(2, 1)

	return VAO, VBO, colorVBO
}

func main() {
	numX := 1000
	numY := 1000
	title := "One object method"
	fmt.Println("Starting")
	
	window := glfwBoilerplate.InitGLFW(title,
		800, 600, false)
	defer glfw.Terminate()

	ourShader := shader.MakeShaders("oneObject.vs", "oneObject.fs")
	board := conways.CreateBoard(192921, numX, numY)

	VAO, VBO, colorVBO := makeBuffers(board, numX, numY)
	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteVertexArrays(1, &VBO)
	defer gl.DeleteVertexArrays(1, &colorVBO)
	
	lastTime := 0.0
	numFrames := 0.0
	verticesPerSquare := int32(6)
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	for !window.ShouldClose() {
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)	
		
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		
		// Update board and VBO
		conways.UpdateBoard(board, numX, numY)
		colors := getColors(board, numX, numY)
		
		gl.BindBuffer(gl.ARRAY_BUFFER, colorVBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(colors) * 4,
			gl.Ptr(colors), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		
		// Render cubes
		ourShader.Use()
		ourShader.SetFloat("fragColor", 1.0)
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0,
			verticesPerSquare * int32(numX) * int32(numY))
		gl.BindVertexArray(0)
		
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
