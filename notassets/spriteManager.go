package notassets

import (
	"fmt"
	"sync"
)

type TextureManager struct {
	textures map[string]*Texture
	mu       sync.RWMutex
}

func NewTextureManager() *TextureManager {
	return &TextureManager{
		textures: make(map[string]*Texture),
	}
}

func (tm *TextureManager) Load(name, path string) (*Texture, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tex, ok := tm.textures[name]; ok {
		return tex, nil
	}

	tex, err := LoadTexture(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load texture '%s' from '%s': %w", name, path, err)
	}

	tm.textures[name] = tex
	return tex, nil
}

func (tm *TextureManager) Get(name string) (*Texture, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tex, ok := tm.textures[name]
	if !ok {
		return nil, fmt.Errorf("texture '%s' not found", name)
	}

	return tex, nil
}

func (tm *TextureManager) Unload(name string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tex, ok := tm.textures[name]
	if !ok {
		return fmt.Errorf("texture '%s' not found", name)
	}

	tex.Delete()
	delete(tm.textures, name)
	return nil
}

func (tm *TextureManager) Clear() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for name, tex := range tm.textures {
		tex.Delete()
		delete(tm.textures, name)
	}
}

func (tm *TextureManager) Count() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return len(tm.textures)
}

func (tm *TextureManager) List() []string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	names := make([]string, 0, len(tm.textures))
	for name := range tm.textures {
		names = append(names, name)
	}
	return names
}

type SpriteManager struct {
	sprites  map[string]*Sprite
	textures *TextureManager
	mu       sync.RWMutex
}

func NewSpriteManager(textureManager *TextureManager) *SpriteManager {
	return &SpriteManager{
		sprites:  make(map[string]*Sprite),
		textures: textureManager,
	}
}

// Create creates a new sprite from a loaded texture
func (sm *SpriteManager) Create(name string, texture *Texture) (*Sprite, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.sprites[name]; exists {
		return nil, fmt.Errorf("sprite '%s' already exists", name)
	}

	sprite := &Sprite{
		Texture: texture,
		Name:    name,
		X:       0,
		Y:       0,
	}

	sm.sprites[name] = sprite
	return sprite, nil
}

// LoadAndCreate loads a texture and creates a sprite from it
func (sm *SpriteManager) LoadAndCreate(spriteName, texturePath string) (*Sprite, error) {
	texture, err := sm.textures.Load(spriteName, texturePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load texture: %w", err)
	}

	return sm.Create(spriteName, texture)
}

// Get retrieves a sprite by name
func (sm *SpriteManager) Get(name string) (*Sprite, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sprite, exists := sm.sprites[name]
	if !exists {
		return nil, fmt.Errorf("sprite '%s' not found", name)
	}

	return sprite, nil
}

// Remove removes a sprite (does not delete the texture)
func (sm *SpriteManager) Remove(name string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.sprites[name]; !exists {
		return fmt.Errorf("sprite '%s' not found", name)
	}

	delete(sm.sprites, name)
	return nil
}

// Clear removes all sprites (does not delete textures)
func (sm *SpriteManager) Clear() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.sprites = make(map[string]*Sprite)
}

// Count returns the number of sprites
func (sm *SpriteManager) Count() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.sprites)
}
