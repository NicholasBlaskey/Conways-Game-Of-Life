// Translated from
// https://github.com/JoeyDeVries/LearnOpenGL/blob/master/src/2.lighting/1.colors/colors.cpp

package main

import(
	"runtime"
	"fmt"
	"math"
	
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	
	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

	"github.com/nicholasblaskey/Conways-Game-Of-Life/glfwBoilerplate"	
	"github.com/nicholasblaskey/Conways-Game-Of-Life/conways"
)

const tolerance = 0.00001
const windowWidth = 1920
const windowHeight = 986

func init() {
	runtime.LockOSThread()
}

func makeBuffers(board []float32, numX,
	numY int) (uint32, uint32, uint32, uint32) {	

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
	quadVertices := []float32{
		// Positions // texture coords
		-1.0,  1.0,  0.0, 1.0,
        -1.0, -1.0,  0.0, 0.0,
         1.0, -1.0,  1.0, 0.0,

        -1.0,  1.0,  0.0, 1.0,
         1.0, -1.0,  1.0, 0.0,
         1.0,  1.0,  1.0, 1.0,
	}
	// TileVAO
	var VAO, VBO uint32		
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(Vertices) * 4,
		gl.Ptr(Vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)	
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2 * 4, gl.PtrOffset(0))
	// quadVAO (screen coords to render our framebuffer to)
	var quadVAO, quadVBO uint32	
	gl.GenVertexArrays(1, &quadVAO)
	gl.GenBuffers(1, &quadVBO)
	gl.BindVertexArray(quadVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, quadVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadVertices) * 4,
		gl.Ptr(quadVertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)	
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4 * 4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4 * 4,
		gl.PtrOffset(2 * 4))
	gl.BindVertexArray(0)
	
	return VAO, VBO, quadVAO, quadVBO
}

func createFramebuffer() (uint32, uint32) {
	// Frambuffer config
	var framebuffer uint32
	gl.GenFramebuffers(1, &framebuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
	// Create a color attachment texture
	var textureColorbuffer uint32
	gl.GenTextures(1, &textureColorbuffer)
	gl.BindTexture(gl.TEXTURE_2D, textureColorbuffer)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, windowWidth, windowHeight,
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(nil))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0,
		gl.TEXTURE_2D, textureColorbuffer, 0)
	// Create a renderbuffer object for depth and stencil atachment
	var rbo uint32
	gl.GenRenderbuffers(1, &rbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, rbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8,
		windowWidth, windowHeight)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT,
		gl.RENDERBUFFER, rbo)
	// Ensure framebuffer is complete
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	}

	return framebuffer, textureColorbuffer
}

func main() {
	numX := 100
	numY := 100
	title := "Only updated method"
	fmt.Println("Starting")

	window := glfwBoilerplate.InitGLFW(title,
		windowWidth, windowHeight, false)
	defer glfw.Terminate()

	ourShader := shader.MakeShaders("onlyUpdated.vs", "onlyUpdated.fs")
	screenShader := shader.MakeShaders("screen.vs", "screen.fs")
	translations := conways.GetPositions(numX, numY)

	curBoard := conways.CreateBoard(192921, numX, numY)
	prevBoard := make([]float32, len(curBoard))

	VAO, VBO, quadVAO, quadVBO := makeBuffers(curBoard, numX, numY)
	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteVertexArrays(1, &VBO)
	defer gl.DeleteVertexArrays(1, &quadVAO)
	defer gl.DeleteVertexArrays(1, &quadVBO)

	framebuffer, textureColorbuffer := createFramebuffer()
	gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
	gl.Enable(gl.DEPTH_TEST)
	// Draw every vertex as the first iteration
	ourShader.Use()
	for i := 0; i < len(translations); i++ {
		ourShader.SetVec2("aOffset", translations[i])
		ourShader.SetFloat("fragColor", curBoard[i])
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		gl.BindVertexArray(0)
	}
	
	// Update loop
	lastTime := 0.0
	numFrames := 0.0
	for !window.ShouldClose() {
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)			
		// Update board
		copy(prevBoard, curBoard)
		conways.UpdateBoard(curBoard, numX, numY)
		
		// Bind frame buffer
		gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
		gl.Enable(gl.DEPTH_TEST)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		
		// Render only updated vertices to framebuffer
		ourShader.Use()
		for i := 0; i < len(translations); i++ {
			if math.Abs(float64(prevBoard[i] - curBoard[i])) > tolerance {
				ourShader.SetVec2("aOffset", translations[i])
				ourShader.SetFloat("fragColor", curBoard[i])
				gl.BindVertexArray(VAO)
				gl.DrawArrays(gl.TRIANGLES, 0, 6)
				gl.BindVertexArray(0)
			}
		}
		
		// Bind back framebuffer for the screen
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
		gl.Disable(gl.DEPTH_TEST)
		// Clear all relevant buffers
		gl.ClearColor(1.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		// Render our other framebuffer as a texture to the screen
		screenShader.Use()
		gl.BindVertexArray(quadVAO)
		gl.BindTexture(gl.TEXTURE_2D, textureColorbuffer)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

