package notacore

import (
	"NotaborEngine/notagl"
	"fmt"
	"path/filepath"
)

// LoadTexture loads a texture and creates OpenGL texture (with context)
func (w *GlfwWindow2D) LoadTexture(name, path string) (*notagl.Texture, error) {
	w.MakeContextCurrent()

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	tex, err := w.RunTime.TextureMgr.Load(name, absPath, true) // true = create GL texture
	if err != nil {
		return nil, fmt.Errorf("failed to load texture: %w", err)
	}

	return tex, nil
}

// GetTexture retrieves a loaded texture
func (w *GlfwWindow2D) GetTexture(name string) (*notagl.Texture, error) {
	return w.RunTime.TextureMgr.Get(name)
}

// UnloadTexture removes a texture
func (w *GlfwWindow2D) UnloadTexture(name string) error {
	return w.RunTime.TextureMgr.Unload(name)
}

func (w *GlfwWindow3D) LoadTexture(name, path string) (*notagl.Texture, error) {
	w.MakeContextCurrent()

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	tex, err := w.RunTime.TextureMgr.Load(name, absPath, true)
	if err != nil {
		return nil, fmt.Errorf("failed to load texture: %w", err)
	}

	return tex, nil
}

func (w *GlfwWindow3D) GetTexture(name string) (*notagl.Texture, error) {
	return w.RunTime.TextureMgr.Get(name)
}

func (w *GlfwWindow3D) UnloadTexture(name string) error {
	return w.RunTime.TextureMgr.Unload(name)
}
