package notaentity

import (
	"NotaborEngine/internal/notacolor"
	"NotaborEngine/notamath"
	"testing"
)

func TestSparseSet(t *testing.T) {
	set := NewSparseSet(16)

	// Test Add
	if !set.Add(5) {
		t.Error("Expected Add(5) to return true")
	}
	if !set.Add(10) {
		t.Error("Expected Add(10) to return true")
	}
	if !set.Add(3) {
		t.Error("Expected Add(3) to return true")
	}

	// Test Contains
	if !set.Contains(5) || !set.Contains(10) || !set.Contains(3) {
		t.Error("Expected all added entities to be contained")
	}
	if set.Contains(7) {
		t.Error("Expected entity 7 to not be contained")
	}

	// Test duplicate Add
	if set.Add(5) {
		t.Error("Expected Add(5) to return false for duplicate")
	}

	// Test Size
	if set.Size() != 3 {
		t.Errorf("Expected size 3, got %d", set.Size())
	}

	// Test Remove
	if !set.Remove(10) {
		t.Error("Expected Remove(10) to return true")
	}
	if set.Contains(10) {
		t.Error("Expected entity 10 to not be contained after removal")
	}
	if set.Size() != 2 {
		t.Errorf("Expected size 2 after removal, got %d", set.Size())
	}

	// Test duplicate Remove
	if set.Remove(10) {
		t.Error("Expected Remove(10) to return false for non-existent entity")
	}

	// Test GetEntities
	entities := set.GetEntities()
	if len(entities) != 2 {
		t.Errorf("Expected 2 entities, got %d", len(entities))
	}
}

func TestEntityManager_CreateEntity(t *testing.T) {
	em := NewEntityManager()

	// Create entities
	e1 := em.CreateEntity("player")
	if e1 == nil {
		t.Fatal("Expected entity to be created")
	}
	if e1.GetId() != "player" {
		t.Errorf("Expected entity name 'player', got '%s'", e1.GetId())
	}

	e2 := em.CreateEntity("enemy")
	if e2 == nil {
		t.Fatal("Expected entity to be created")
	}
	if e2.GetId() != "enemy" {
		t.Errorf("Expected entity name 'enemy', got '%s'", e2.GetId())
	}

	// Test GetEntity
	retrieved := em.GetEntity("player")
	if retrieved == nil {
		t.Fatal("Expected to retrieve entity 'player'")
	}
	if retrieved.ID != e1.ID {
		t.Error("Expected retrieved entity to match created entity")
	}

	// Test duplicate name (should create new entity with same name reference)
	e3 := em.CreateEntity("player")
	if e3 == nil {
		t.Fatal("Expected entity to be created with duplicate name")
	}
	// The new entity should overwrite the name mapping
	retrieved = em.GetEntity("player")
	if retrieved.ID != e3.ID {
		t.Error("Expected name 'player' to map to the latest entity")
	}
}

func TestEntityManager_Transforms(t *testing.T) {
	em := NewEntityManager()
	e := em.CreateEntity("test")

	// Test initial transform
	pos := e.Position()
	if pos != (notamath.Vec2{}) {
		t.Errorf("Expected initial position (0,0), got (%f,%f)", pos.X, pos.Y)
	}

	rot := e.Rotation()
	if rot != 0 {
		t.Errorf("Expected initial rotation 0, got %f", rot)
	}

	scale := e.ScaleValue()
	if scale.X != 1 || scale.Y != 1 {
		t.Errorf("Expected initial scale (1,1), got (%f,%f)", scale.X, scale.Y)
	}

	// Test Move
	e.Move(notamath.Vec2{X: 10, Y: 20})
	em.Flush()
	pos = e.Position()
	if pos.X != 10 || pos.Y != 20 {
		t.Errorf("Expected position (10,20) after move, got (%f,%f)", pos.X, pos.Y)
	}

	// Test Rotate
	e.Rotate(1.57)
	em.Flush()
	rot = e.Rotation()
	if rot < 1.56 || rot > 1.58 {
		t.Errorf("Expected rotation ~1.57, got %f", rot)
	}

	// Test Scale
	e.Scale(notamath.Vec2{X: 2, Y: 3})
	em.Flush()
	scale = e.ScaleValue()
	if scale.X != 2 || scale.Y != 3 {
		t.Errorf("Expected scale (2,3), got (%f,%f)", scale.X, scale.Y)
	}
}

func TestEntityManager_Components(t *testing.T) {
	em := NewEntityManager()
	e := em.CreateEntity("test")

	// Test Color
	e.WithColor(notacolor.Red)
	color := em.GetColor(e.ID)
	if color.R != notacolor.Red.R || color.G != notacolor.Red.G || color.B != notacolor.Red.B {
		t.Error("Expected color to be red")
	}

	// Test Active/Visible
	if !e.IsActive() {
		t.Error("Expected entity to be active by default")
	}
	if !e.IsVisible() {
		t.Error("Expected entity to be visible by default")
	}

	e.SetActive(false)
	if e.IsActive() {
		t.Error("Expected entity to be inactive after SetActive(false)")
	}

	e.SetVisible(false)
	if e.IsVisible() {
		t.Error("Expected entity to be invisible after SetVisible(false)")
	}
}

func TestEntityManager_CollisionGroups(t *testing.T) {
	em := NewEntityManager()
	e1 := em.CreateEntity("entity1")
	e2 := em.CreateEntity("entity2")

	// Add to collision group
	em.AddToCollisionGroup("group1", e1)
	em.AddToCollisionGroup("group1", e2)

	// Verify group was created
	if _, ok := em.collisionGroups["group1"]; !ok {
		t.Error("Expected collision group 'group1' to exist")
	}
}

func TestEntityManager_Remove(t *testing.T) {
	em := NewEntityManager()
	em.CreateEntity("test")

	e := em.GetEntity("test")
	if e == nil {
		t.Fatal("Expected entity to exist")
	}

	em.Remove("test")

	e = em.GetEntity("test")
	if e != nil {
		t.Error("Expected entity to be removed")
	}
}

func TestEntityManager_EntityReuse(t *testing.T) {
	em := NewEntityManager()

	// Create and remove entity
	e1 := em.CreateEntity("test1")
	id1 := e1.ID
	em.Remove("test1")

	// Create new entity - should reuse the ID
	e2 := em.CreateEntity("test2")
	if e2.ID != id1 {
		t.Errorf("Expected ID reuse, got %d instead of %d", e2.ID, id1)
	}
}

func TestArchetype(t *testing.T) {
	var a Archetype = 0

	a = a.Add(ComponentTransform)
	if !a.Has(ComponentTransform) {
		t.Error("Expected archetype to have Transform component")
	}

	a = a.Add(ComponentSprite)
	if !a.Has(ComponentSprite) {
		t.Error("Expected archetype to have Sprite component")
	}

	a = a.Remove(ComponentTransform)
	if a.Has(ComponentTransform) {
		t.Error("Expected archetype to not have Transform component after removal")
	}
}
