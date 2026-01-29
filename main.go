// cmd/test_entity/main.go
package main

import (
	"fmt"
	"log"

	"NotaborEngine/notagl"
	"NotaborEngine/notamath"
	"NotaborEngine/notashader"
	"NotaborEngine/notassets"
)

func main() {
	// Create a mock texture
	mockTexture := &notassets.Texture{
		ID:     1,
		Width:  64,
		Height: 64,
	}

	// Create a sprite
	sprite := &notassets.Sprite{
		Texture: mockTexture,
		Name:    "player_sprite",
		X:       0,
		Y:       0,
	}

	// Test 1: Create entity with sprite
	entity1 := notassets.NewEntity("ent1", "Player", sprite)

	fmt.Printf("Entity 1:\n")
	fmt.Printf("  ID: %s, Name: %s\n", entity1.ID, entity1.Name)
	fmt.Printf("  Active: %v, Visible: %v\n", entity1.Active, entity1.Visible)
	fmt.Printf("  Sprite position: (%d, %d)\n", entity1.Sprite.X, entity1.Sprite.Y)

	// Test position setting
	entity1.SetPosition(100, 150)
	x, y := entity1.GetPosition()
	fmt.Printf("  After SetPosition(100, 150): (%.0f, %.0f)\n", x, y)

	// Test 2: Create entity with polygon
	rect := notagl.CreateRectangle(notamath.Po2{X: 0, Y: 0}, 50, 50)
	rect.SetColor(notashader.Color{R: 1, G: 0, B: 0, A: 1})

	entity2 := notassets.NewEntityWithPolygon("ent2", "Red Square", rect)

	fmt.Printf("\nEntity 2:\n")
	fmt.Printf("  ID: %s, Name: %s\n", entity2.ID, entity2.Name)
	fmt.Printf("  Has polygon: %v (vertices: %d)\n",
		len(entity2.Polygon.Vertices) > 0,
		len(entity2.Polygon.Vertices))

	entity2.SetPosition(200, 300)
	x2, y2 := entity2.GetPosition()
	fmt.Printf("  Position: (%.0f, %.0f)\n", x2, y2)

	// Test 3: Toggle visibility
	entity1.Visible = false
	entity2.Active = false

	fmt.Printf("\nState changes:\n")
	fmt.Printf("  Entity1 visible: %v\n", entity1.Visible)
	fmt.Printf("  Entity2 active: %v\n", entity2.Active)

	log.Println("Entity test completed successfully")
}
