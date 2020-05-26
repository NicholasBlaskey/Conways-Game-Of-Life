// Translated from
// https://github.com/JoeyDeVries/LearnOpenGL/blob/master/src/2.lighting/1.colors/colors.cpp

package main

import(
	"runtime"
	"fmt"
	"math"
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

func makeChunkBuffer(chunkSize, layout int,xOffset, yOffset float32) (uint32, int) {

	fmt.Println(xOffset + yOffset)
	
	Vertices := []float32{}
	fmt.Println(Vertices)

	fmt.Println(layout)
	//prevColor := -1
	for i := 0; i < chunkSize; i++ {
		curColor := layout % 2
		layout /= 2

		// TODO lets implement vertex skipping later on
		// Right now lets get it done
		// Skip vertex if it is the same color
		//if curColor != prevColor {
		//	Vertices
		//} 
		//fmt.Printf("color=%d, layout=%d\n", color, layout)
	}

	/*
	Vertices := []float32{
		// positions     
		-xOffset,  yOffset, 
		xOffset, -yOffset,
		-xOffset, -yOffset, 

		-xOffset,  yOffset, 
		xOffset, -yOffset, 
		xOffset,  yOffset, 
	}
	var VAO, VBO uint32		
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	
	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(Vertices) * 4,
		gl.Ptr(Vertices), gl.STATIC_DRAW)
	
	gl.EnableVertexAttribArray(0)	
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2 * 4, gl.PtrOffset(0))
	
*/

	return 3, 3//VAO, VBO

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
	translations := conways.GetPositions(numX, numY)
	board := conways.CreateBoard(192921, numX, numY)
	fmt.Println(ourShader)
	fmt.Println(translations)
	
	VAO, VBO := makeBuffers(numX, numY, chunkSize)
//	defer gl.DeleteVertexArrays(1, &VAO)
//	defer gl.DeleteVertexArrays(1, &VBO)

	fmt.Println(VBO)
	fmt.Println(VAO)
	
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

		// Render cubes
		/*
		ourShader.Use()
		for i := 0; i < len(translations); i++ {
			ourShader.SetVec2("aOffset", translations[i])
			ourShader.SetFloat("fragColor", board[i])
			gl.BindVertexArray(VAO)
			gl.DrawArrays(gl.TRIANGLES, 0, 6)
			gl.BindVertexArray(0)
		}
*/
		
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

