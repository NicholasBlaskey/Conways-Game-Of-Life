package main

import (
	"fmt"
	"math"
	"runtime"

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

func makeChunkBuffer(chunkSize, layout int, xOffset,
	yOffset float32) (uint32, int) {

	// First iteration we will always draw it
	Vertices := []float32{}
	numVertices := 2
	prevColor := float32(layout % 2)
	prevColorIndex := 0
	Vertices = append(Vertices,
		-xOffset, -yOffset, prevColor, // Bot left
		-xOffset, yOffset, prevColor) // Top left

	for i := 1; i < chunkSize; i++ {
		layout /= 2
		curColor := float32(layout % 2)

		if curColor != prevColor {
			Vertices = append(Vertices,
				// Finish previous square
				// Bot right of first triangle
				xOffset+xOffset*float32((i-1)*2), -yOffset, prevColor,

				xOffset+xOffset*float32((i-1)*2), -yOffset, prevColor,
				xOffset+xOffset*float32((i-1)*2), yOffset, prevColor,
				-xOffset+xOffset*float32(prevColorIndex*2), yOffset, prevColor)

			// Start new square
			Vertices = append(Vertices,
				-xOffset+xOffset*float32(i*2), yOffset, curColor,
				-xOffset+xOffset*float32(i*2), -yOffset, curColor)
			numVertices += 6
			prevColorIndex = i
		}
		prevColor = curColor
	}

	// Final iteration we will always finish the square
	Vertices = append(Vertices,
		xOffset+xOffset*float32((chunkSize-1)*2), yOffset, prevColor,
		// Top right most
		xOffset+xOffset*float32((chunkSize-1)*2), yOffset, prevColor,
		// Bot right most
		xOffset+xOffset*float32((chunkSize-1)*2), -yOffset, prevColor,
		// Bot left
		-xOffset+xOffset*float32(prevColorIndex*2), -yOffset, prevColor,
	)
	numVertices += 4

	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(Vertices)*4,
		gl.Ptr(Vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 1, gl.FLOAT, false, 3*4,
		gl.PtrOffset(2*4))

	return VAO, numVertices

}

func makeBuffers(numX, numY, chunkSize int) ([]uint32, []int) {
	VAOS := []uint32{}
	sizes := []int{}
	xOffset := 1.0 / float32(numX)
	yOffset := 1.0 / float32(numY)
	for i := 0; i < int(math.Pow(2, float64(chunkSize))); i++ {
		VAO, numVertices := makeChunkBuffer(chunkSize, i, xOffset, yOffset)
		VAOS = append(VAOS, VAO)
		sizes = append(sizes, numVertices)
	}

	return VAOS, sizes
}

func getPositions(numX, numY, chunkSize int) []mgl.Vec2 {
	translations := []mgl.Vec2{}
	xOffset := 1.0 / float32(numX)
	yOffset := 1.0 / float32(numY)
	for y := -numY; y < numY; y += 2 {
		for x := -numX; x < numX; x += 2 * chunkSize {
			translations = append(translations,
				mgl.Vec2{float32(x)/float32(numX) + xOffset,
					float32(y)/float32(numY) + yOffset})
		}
	}
	return translations
}

func main() {
	numX := 1000
	numY := 1000
	chunkSize := 5
	title := "Chunks method"
	fmt.Println("Starting")

	window := glfwBoilerplate.InitGLFW(title,
		800, 600, false)
	defer glfw.Terminate()

	ourShader := shader.MakeShaders("chunks.vs", "chunks.fs")
	translations := getPositions(numX, numY, chunkSize)
	board := conways.CreateBoard(192921, numX, numY)
	VAOS, numVertices := makeBuffers(numX, numY, chunkSize)

	lastTime := 0.0
	numFrames := 0.0
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	for !window.ShouldClose() {
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)

		gl.ClearColor(0.3, 0.5, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		// Update board
		conways.UpdateBoard(board, numX, numY)
		ourShader.Use()

		// Render VAO
		for i := 0; i < len(translations); i++ {
			ourShader.SetVec2("aOffset", translations[i])

			// We need to hard code the indexOfVAO calculation
			// Making it work for any chunk size causes too much of a
			// preformance hit. Even when storing powers of 2.
			indexOfVAO := int(board[i*5] + board[i*5+1]*2 +
				board[i*5+2]*4 + board[i*5+3]*8 + board[i*5+4]*16)
			gl.BindVertexArray(VAOS[indexOfVAO])
			gl.DrawArrays(gl.TRIANGLES, 0, int32(numVertices[indexOfVAO]))
			gl.BindVertexArray(0)
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
