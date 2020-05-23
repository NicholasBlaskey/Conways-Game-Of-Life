package glfwBoilerplate

import(
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func InitGLFW(windowTitle string,
	width, height int, useDepthTest bool) *glfw.Window {
	
    if err := glfw.Init(); err != nil {
        panic(err)
    }
	
    glfw.WindowHint(glfw.ContextVersionMajor, 4)
    glfw.WindowHint(glfw.ContextVersionMinor, 1)
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
    window, err := glfw.CreateWindow(
        width, height, windowTitle, nil, nil)
    if err != nil {
        panic(err)
    }
    window.MakeContextCurrent()

	window.SetFramebufferSizeCallback(
		glfw.FramebufferSizeCallback(framebuffer_size_callback))
	window.SetKeyCallback(keyCallback)

    if err := gl.Init(); err != nil {
        panic(err)
    }

	if useDepthTest {
		gl.Enable(gl.DEPTH_TEST)
	}
	
	return window
}

func framebuffer_size_callback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int,
	action glfw.Action, mods glfw.ModifierKey) {
	// Escape closes window
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

