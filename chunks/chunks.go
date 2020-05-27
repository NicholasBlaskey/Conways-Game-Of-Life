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

	Vertices := []float32{}	
	numVertices := 0
	//prevColor := -1
	for i := 0; i < chunkSize; i++ {
		curColor := float32(layout % 2)
		layout /= 2

		numVertices += 6
		Vertices = append(Vertices,
			// Position   // Color
			-xOffset + xOffset * float32(i * 2),  yOffset, curColor, 
			xOffset + xOffset * float32(i * 2), -yOffset, curColor,
			-xOffset + xOffset * float32(i * 2), -yOffset, curColor,
			
			-xOffset + xOffset * float32(i * 2),  yOffset, curColor,
			xOffset + xOffset * float32(i * 2), -yOffset, curColor,
			xOffset + xOffset * float32(i * 2),  yOffset, curColor)
		
		// TODO lets implement vertex skipping later on
		// Right now lets get it done
		// Skip vertex if it is the same color
		//if curColor != prevColor {
		//	Vertices
		//} 
		//fmt.Printf("color=%d, layout=%d\n", curColor, layout)
	}

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
	numX := 100
	numY := 100
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

	/*
	powers := []float32{}
	for i := 0; i < chunkSize; i++ {
		powers = append(powers, float32(math.Pow(2, float64(i))))
	}
*/
	
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	for !window.ShouldClose() {
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)	
		
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		// Update board and VBO
		conways.UpdateBoard(board, numX, numY)

		ourShader.Use()
		for i := 0; i < len(translations); i++ {
			ourShader.SetVec2("aOffset", translations[i])

			/*
			indexOfVAO := 0
			for j := 0; j < chunkSize; j++ {
				indexOfVAO += int(board[i * chunkSize + j] * powers[j])
			}
			fmt.Println(indexOfVAO)
            */
			
			indexOfVAO := int(board[i * 5] + board[i * 5 + 1] * 2 +
				board[i * 5 + 2] * 4 + board[i * 5 + 3] * 8 +
				board[i * 5 + 4] * 16)
			gl.BindVertexArray(VAOS[indexOfVAO])
			gl.DrawArrays(gl.TRIANGLES, 0, int32(numVertices[indexOfVAO]))
			gl.BindVertexArray(0)
		}

		time.Sleep(1 * time.Millisecond)
		
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

