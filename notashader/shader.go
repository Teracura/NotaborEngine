package notashader

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type ShaderProgram struct {
	Type uint32
	Err  error
}

var Shaders map[string]uint32 = make(map[string]uint32)

type Program interface {
	CompileShader(source string, shaderType uint32) ShaderProgram
	CreateProgram(vertexSrc, fragmentSrc string) ShaderProgram
}

func CompileShader(source string, shaderType uint32) ShaderProgram {
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
		return ShaderProgram{
			Type: 0,
			Err:  fmt.Errorf("failed to compile shader: %s", log),
		}
	}

	return ShaderProgram{
		Type: shader,
		Err:  nil,
	}
}

func CreateProgram(vertexSrc, fragmentSrc string) ShaderProgram {
	vert := CompileShader(vertexSrc, gl.VERTEX_SHADER)
	if vert.Err != nil {
		return ShaderProgram{
			Type: 0,
			Err:  vert.Err,
		}
	}

	frag := CompileShader(fragmentSrc, gl.FRAGMENT_SHADER)
	if frag.Err != nil {
		return ShaderProgram{
			Type: 0,
			Err:  frag.Err,
		}
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vert.Type)
	gl.AttachShader(prog, frag.Type)
	gl.LinkProgram(prog)

	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength+1)
		gl.GetProgramInfoLog(prog, logLength, nil, &log[0])
		return ShaderProgram{
			Type: 0,
			Err:  fmt.Errorf("failed to link program: %s", log),
		}
	}

	// deleting is intentional
	gl.DeleteShader(vert.Type)
	gl.DeleteShader(frag.Type)

	return ShaderProgram{Type: prog, Err: nil}
}
