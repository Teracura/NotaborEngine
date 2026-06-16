package notaentity

import (
	"NotaborEngine/internal/notacollision"
	"NotaborEngine/internal/notacolor"
	"NotaborEngine/internal/notageometry"
	"NotaborEngine/internal/notamath"
	"NotaborEngine/internal/notashader"
	"NotaborEngine/internal/notatexture"
)

// ComponentType represents a unique identifier for component types.
type ComponentType uint32

const (
	ComponentTransform ComponentType = iota
	ComponentSprite
	ComponentPolygon
	ComponentCollider
	ComponentColor
	ComponentShader
	ComponentMaterial
	ComponentCount
)

// TransformComponent stores the transform data for an entity.
type TransformComponent struct {
	Position notamath.Vec2
	Rotation float32
	Scale    notamath.Vec2
}

// SpriteComponent stores the sprite reference for an entity.
type SpriteComponent struct {
	Sprite *notatexture.Sprite
}

// PolygonComponent stores the polygon geometry for an entity.
type PolygonComponent struct {
	Polygon *notageometry.Polygon
}

// ColliderComponent stores the collider for collision detection.
type ColliderComponent struct {
	Collider notacollision.Collider
}

// ColorComponent stores the color tint for an entity.
type ColorComponent struct {
	Color notacolor.Color
}

// ShaderComponent stores the custom shader for an entity.
type ShaderComponent struct {
	Shader *notashader.Shader
}

// MaterialComponent stores the material (shader instance) for an entity.
type MaterialComponent struct {
	Material *notashader.Material
}

// Archetype is a bitmask representing which components an entity has.
type Archetype uint32

// Archetype masks for each component type
const (
	ArchetypeTransform Archetype = 1 << ComponentTransform
	ArchetypeSprite    Archetype = 1 << ComponentSprite
	ArchetypePolygon   Archetype = 1 << ComponentPolygon
	ArchetypeCollider  Archetype = 1 << ComponentCollider
	ArchetypeColor     Archetype = 1 << ComponentColor
	ArchetypeShader    Archetype = 1 << ComponentShader
	ArchetypeMaterial  Archetype = 1 << ComponentMaterial
)

// Has checks if this archetype includes the given component type.
func (a Archetype) Has(componentType ComponentType) bool {
	return (a & (1 << componentType)) != 0
}

// Add adds a component type to this archetype.
func (a Archetype) Add(componentType ComponentType) Archetype {
	return a | (1 << componentType)
}

// Remove removes a component type from this archetype.
func (a Archetype) Remove(componentType ComponentType) Archetype {
	return a &^ (1 << componentType)
}
