package notacore

import (
	"NotaborEngine/notashader"
	"errors"
)

func (w *GlfwWindow2D) DeleteShader(name string) uint32 {
	program := w.Shaders[name]
	delete(w.Shaders, name)
	return program
}

func (w *GlfwWindow2D) UpdateShader(shader notashader.Shader) error {
	if w.Shaders == nil {
		w.Shaders = make(map[string]uint32)
	}
	_, err := w.GetShader(shader.Name)
	if err != nil {
		return err
	}

	w.Shaders[shader.Name] = notashader.CreateProgram(shader.VertexString, shader.FragmentString)
	return nil
}

func (w *GlfwWindow2D) GetShader(name string) (uint32, error) {
	value, found := w.Shaders[name]
	if !found {
		return 0, errors.New("shader with name " + name + " is not found")
	}
	return value, nil
}

func (w *GlfwWindow2D) CreateShader(shader notashader.Shader) error {
	if w.Shaders == nil {
		w.Shaders = make(map[string]uint32)
	}
	_, err := w.GetShader(shader.Name)
	if err == nil {
		return errors.New("shader with name " + shader.Name + " already exists")
	}
	w.MakeContextCurrent()
	w.Shaders[shader.Name] = notashader.CreateProgram(shader.VertexString, shader.FragmentString)
	return nil
}

func (w *GlfwWindow3D) DeleteShader(name string) uint32 {
	program := w.Shaders[name]
	delete(w.Shaders, name)
	return program
}

func (w *GlfwWindow3D) UpdateShader(shader notashader.Shader) error {
	if w.Shaders == nil {
		w.Shaders = make(map[string]uint32)
	}
	_, err := w.GetShader(shader.Name)
	if err != nil {
		return err
	}

	w.Shaders[shader.Name] = notashader.CreateProgram(shader.VertexString, shader.FragmentString)
	return nil
}

func (w *GlfwWindow3D) GetShader(name string) (uint32, error) {
	value, notFound := w.Shaders[name]
	if notFound {
		return 0, errors.New("shader with name " + name + " is not found")
	}
	return value, nil
}

func (w *GlfwWindow3D) CreateShader(shader notashader.Shader) error {
	if w.Shaders == nil {
		w.Shaders = make(map[string]uint32)
	}
	_, err := w.GetShader(shader.Name)
	if err == nil {
		return errors.New("shader with name " + shader.Name + " already exists")
	}

	w.Shaders[shader.Name] = notashader.CreateProgram(shader.VertexString, shader.FragmentString)
	return nil
}
