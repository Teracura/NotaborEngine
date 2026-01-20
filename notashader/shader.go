package notashader

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Shader struct {
	Name           string
	VertexString   string
	FragmentString string
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	sources, free := gl.Strs(source + "\x00")
	defer free()
	gl.ShaderSource(shader, 1, sources, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength+1)
		gl.GetShaderInfoLog(shader, logLength, nil, &log[0])
		return 0, fmt.Errorf("failed to compile shader: %s", log)
	}

	return shader, nil
}

func CreateProgram(vertexSrc, fragmentSrc string) uint32 {
	vert, err := compileShader(vertexSrc, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	frag, err := compileShader(fragmentSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vert)
	gl.AttachShader(prog, frag)
	gl.LinkProgram(prog)

	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength+1)
		gl.GetProgramInfoLog(prog, logLength, nil, &log[0])
		panic(fmt.Sprintf("failed to link program: %s", log))
	}

	// deleting is intentional
	gl.DeleteShader(vert)
	gl.DeleteShader(frag)

	return prog
}
