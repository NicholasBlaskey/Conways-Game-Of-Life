// Translated from
// https://github.com/JoeyDeVries/LearnOpenGL/blob/master/src/2.lighting/1.colors/colors.cpp

package main

import(
	"runtime"
	"fmt"
	"math"
	"time"
	
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
	
	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

	"github.com/nicholasblaskey/Conways-Game-Of-Life/glfwBoilerplate"	
	"github.com/nicholasblaskey/Conways-Game-Of-Life/conways"
)

func init() {
	runtime.LockOSThread()
}

func makeChunkBuffer(chunkSize, layout int, xOffset,
	yOffset float32) (uint32, int) {

	fmt.Printf("\nGening layout %d\n", layout)
	
	Vertices := []float32{}	
	numVertices := 2

	// First iteration we will always draw it
	prevColor := float32(layout % 2)
	
	Vertices = append(Vertices,
		-xOffset,  -yOffset, prevColor, // Bot left 
		-xOffset, yOffset, prevColor)   // Top left	
	prevColorIndex := 0
	fmt.Printf("first layout=%d,prevColor=%f\n", layout, prevColor)
	for i := 1; i < chunkSize - 1; i++ { 
		layout /= 2
		curColor := float32(layout % 2)

		//fmt.Println(i)
		fmt.Printf("layout=%d,curColor=%f,prevColor=%f\n", layout, curColor, prevColor)
		
		if curColor != prevColor {
			fmt.Printf("Swapping at i=%d\n", i)
			Vertices = append(Vertices,
				// Finish previous square
				// Bot right of first triangle
				xOffset + xOffset * float32((i - 1) * 2), -yOffset, prevColor,
				
				xOffset + xOffset * float32((i - 1) * 2), -yOffset, prevColor,
				xOffset + xOffset * float32((i - 1) * 2), yOffset, prevColor,
				-xOffset + xOffset * float32(prevColorIndex * 2), yOffset, prevColor,

				// Start new square
				-xOffset + xOffset * float32(i * 2),  yOffset, curColor,
				-xOffset + xOffset * float32(i * 2),  -yOffset, curColor,
			)
				
			prevColorIndex = i
			numVertices += 6
		}
		prevColor = curColor
	}

	// Final iteration we will always draw it
	curColor := float32(layout % 2)
	layout /= 2
	Vertices = append(Vertices,
		xOffset + xOffset * float32((chunkSize - 1) * 2), yOffset, prevColor,
		// Top right most
		xOffset + xOffset * float32((chunkSize - 1) * 2), yOffset, curColor, 
		// Bot right most
		xOffset + xOffset * float32((chunkSize - 1) * 2), -yOffset, curColor,
		// Bot left
		-xOffset + xOffset * float32(prevColorIndex * 2), -yOffset, curColor,
	)
	numVertices += 4
	
	
	var VAO, VBO uint32		
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	
	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(Vertices) * 4,
		gl.Ptr(Vertices), gl.STATIC_DRAW)	
	gl.EnableVertexAttribArray(0)	
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 3 * 4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 1, gl.FLOAT, false, 3 * 4,
		gl.PtrOffset(2 * 4))
		
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
				mgl.Vec2{float32(x) / float32(numX) + xOffset,
					float32(y) / float32(numY) + yOffset})
		}
	}

	return translations
}


func main() {
	numX := 35
	numY := 35
	chunkSize := 5
	title := "Chunks method"
	fmt.Println("Starting")

	//makeChunkBuffer(5, 6, 1.0 / float32(numX), 1.0 / float32(numY))
	//panic("damn")
	
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

		/*
		// Render VAO
		for i := 0; i < len(translations); i++ {
		//for i := len(translations) - 2; i < len(translations); i++ {
			//translations[i][0] += 0.01
			//translations[i][1] += 0.01
			
			ourShader.SetVec2("aOffset", translations[i])
			
			// We need to hard code the indexOfVAO calculation
			// Making it work for any chunk size causes too much of a
			// preformance hit. Even when storing powers of 2.
			indexOfVAO := int(board[i * 5] + board[i * 5 + 1] * 2 +
				board[i * 5 + 2] * 4 + board[i * 5 + 3] * 8 +
				board[i * 5 + 4] * 16)
			gl.BindVertexArray(VAOS[indexOfVAO])
			gl.DrawArrays(gl.TRIANGLES, 0, int32(numVertices[indexOfVAO]))
			gl.BindVertexArray(0)
		}
*/

		for i := 0; i < len(translations); i++ {
			translations[0][1] = 1
			for j := 0; j < len(numVertices); j++ {
				ourShader.Use()
				translations[0][1] += float32(-2) / 35.0

				ourShader.SetVec2("aOffset", translations[0])
				
				gl.BindVertexArray(VAOS[j])
				gl.DrawArrays(gl.TRIANGLES, 0, int32(numVertices[j]))
				gl.BindVertexArray(0)
			}
			break
		}
		
		window.SwapBuffers()
		glfw.PollEvents()

		time.Sleep(1 * time.Millisecond)
	}
}

