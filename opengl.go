package main

import (
    "fmt"
    "log"
    "runtime"

    "github.com/go-gl/gl/v3.2-core/gl"
    "github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
    // GLFW and OpenGL need to be run on the main OS thread
    runtime.LockOSThread()
}

func main() {
    if err := glfw.Init(); err != nil {
        log.Fatalln("Failed to initialize GLFW:", err)
    }
    defer glfw.Terminate()

    glfw.WindowHint(glfw.ContextVersionMajor, 3)
    glfw.WindowHint(glfw.ContextVersionMinor, 2)
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

    window, err := glfw.CreateWindow(800, 600, "OpenGL Triangle Example", nil, nil)
    if err != nil {
        log.Fatalln("Failed to create GLFW window:", err)
    }
    window.MakeContextCurrent()

    if err := gl.Init(); err != nil {
        log.Fatalln("Failed to initialize OpenGL:", err)
    }

    version := gl.GoStr(gl.GetString(gl.VERSION))
    fmt.Println("OpenGL version", version)

    // Set up a simple vertex array and vertex buffer
    var vao uint32
    gl.GenVertexArrays(1, &vao)
    gl.BindVertexArray(vao)

    vertices := []float32{
        //  Position         Color
        0.0,  0.5, 0.0,  1.0, 0.0, 0.0, // Top    (red)
        -0.5, -0.5, 0.0,  0.0, 1.0, 0.0, // Left   (green)
        0.5, -0.5, 0.0,  0.0, 0.0, 1.0, // Right  (blue)
    }

    var vbo uint32
    gl.GenBuffers(1, &vbo)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

    // Set the position attribute
    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
    gl.EnableVertexAttribArray(0)

    // Set the color attribute
    gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
    gl.EnableVertexAttribArray(1)

    // Simple shader setup
    vertexShaderSource := `
        #version 330 core
        layout(location = 0) in vec3 aPos;
        layout(location = 1) in vec3 aColor;

        out vec3 ourColor;

        void main() {
            gl_Position = vec4(aPos, 1.0);
            ourColor = aColor;
        }
    ` + "\x00"

    fragmentShaderSource := `
        #version 330 core
        out vec4 FragColor;
        in vec3 ourColor;

        void main() {
            FragColor = vec4(ourColor, 1.0);
        }
    ` + "\x00"

    vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
    cvertexSource, free := gl.Strs(vertexShaderSource)
    gl.ShaderSource(vertexShader, 1, cvertexSource, nil)
    free()
    gl.CompileShader(vertexShader)

    fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
    cfragmentSource, free := gl.Strs(fragmentShaderSource)
    gl.ShaderSource(fragmentShader, 1, cfragmentSource, nil)
    free()
    gl.CompileShader(fragmentShader)

    shaderProgram := gl.CreateProgram()
    gl.AttachShader(shaderProgram, vertexShader)
    gl.AttachShader(shaderProgram, fragmentShader)
    gl.LinkProgram(shaderProgram)

    gl.DeleteShader(vertexShader)
    gl.DeleteShader(fragmentShader)

    // Use the created shader program
    gl.UseProgram(shaderProgram)

    for !window.ShouldClose() {
        // Clear the screen
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        // Draw the triangle
        gl.BindVertexArray(vao)
        gl.DrawArrays(gl.TRIANGLES, 0, 3)

        // Swap buffers
        window.SwapBuffers()
        glfw.PollEvents()
    }
}
