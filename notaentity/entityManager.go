package notaentity

import (
	"NotaborEngine/internal/notacolor"
	"NotaborEngine/notacollision"
	"NotaborEngine/notageometry"
	"NotaborEngine/notamath"
	"NotaborEngine/notashader"
	"NotaborEngine/notatexture"
	"sync"

	"github.com/viterin/vek/vek32"
)

// CollisionGroup is a set of entities that can collide with each other.
type CollisionGroup struct {
	Name     string
	Entities []EntityID
}

// collisionTable stores collision results between entity pairs.
type collisionTable struct {
	pairs map[EntityID]map[EntityID]notamath.Vec2
}

// ComponentStorage holds arrays of components for entities with a specific archetype.
type ComponentStorage struct {
	transforms []TransformComponent
	sprites    []SpriteComponent
	polygons   []PolygonComponent
	colliders  []ColliderComponent
	colors     []ColorComponent
	shaders    []ShaderComponent
	materials  []MaterialComponent
}

// EntityManager manages entities using an ECS architecture with sparse-set storage.
// It supports component archetypes for efficient iteration and batched transform updates.
type EntityManager struct {
	mu sync.RWMutex

	// Entity management
	nextID     EntityID
	nameToID   map[string]EntityID
	idToName   map[EntityID]string
	freeIDs    []EntityID
	entities   []Entity
	archetypes map[EntityID]Archetype

	// Component storage - arrays indexed by entity ID
	transforms []TransformComponent
	sprites    []*notatexture.Sprite
	polygons   []*notageometry.Polygon
	colliders  []notacollision.Collider
	colors     []notacolor.Color
	shaders    []*notashader.Shader
	materials  []*notashader.Material

	// Entity state
	active  []bool
	visible []bool

	// Pending transform deltas (for batched updates)
	pendingMove     []notamath.Vec2
	pendingRot      []float32
	pendingScale    []notamath.Vec2
	hasPendingMove  []bool
	hasPendingRot   []bool
	hasPendingScale []bool

	// Collision groups
	collisionGroups map[string]*CollisionGroup
	collisionTable  *collisionTable

	// Dirty flags
	dirtyMove  bool
	dirtyRot   bool
	dirtyScale bool
}

// NewEntityManager creates a new EntityManager with initial capacity.
func NewEntityManager() *EntityManager {
	initialCapacity := 64

	em := &EntityManager{
		nextID:          1, // Start from 1 to allow 0 as invalid
		nameToID:        make(map[string]EntityID),
		idToName:        make(map[EntityID]string),
		freeIDs:         make([]EntityID, 0),
		entities:        make([]Entity, initialCapacity),
		archetypes:      make(map[EntityID]Archetype),
		transforms:      make([]TransformComponent, initialCapacity),
		sprites:         make([]*notatexture.Sprite, initialCapacity),
		polygons:        make([]*notageometry.Polygon, initialCapacity),
		colliders:       make([]notacollision.Collider, initialCapacity),
		colors:          make([]notacolor.Color, initialCapacity),
		shaders:         make([]*notashader.Shader, initialCapacity),
		materials:       make([]*notashader.Material, initialCapacity),
		active:          make([]bool, initialCapacity),
		visible:         make([]bool, initialCapacity),
		pendingMove:     make([]notamath.Vec2, initialCapacity),
		pendingRot:      make([]float32, initialCapacity),
		pendingScale:    make([]notamath.Vec2, initialCapacity),
		hasPendingMove:  make([]bool, initialCapacity),
		hasPendingRot:   make([]bool, initialCapacity),
		hasPendingScale: make([]bool, initialCapacity),
		collisionGroups: make(map[string]*CollisionGroup),
		collisionTable: &collisionTable{
			pairs: make(map[EntityID]map[EntityID]notamath.Vec2),
		},
	}

	// Initialize default values for all potential entities
	for i := 0; i < initialCapacity; i++ {
		em.transforms[i] = TransformComponent{
			Scale: notamath.Vec2{X: 1, Y: 1},
		}
		em.colors[i] = notacolor.White
		em.active[i] = true
		em.visible[i] = true
		em.pendingScale[i] = notamath.Vec2{X: 1, Y: 1}
	}

	return em
}

// ensureCapacity ensures the internal arrays can hold at least the given entity ID.
func (em *EntityManager) ensureCapacity(id EntityID) {
	if int(id) < len(em.transforms) {
		return
	}

	newCapacity := max(id*2+1, EntityID(len(em.transforms)*2))

	// Grow all arrays
	newTransforms := make([]TransformComponent, newCapacity)
	copy(newTransforms, em.transforms)
	em.transforms = newTransforms

	newSprites := make([]*notatexture.Sprite, newCapacity)
	copy(newSprites, em.sprites)
	em.sprites = newSprites

	newPolygons := make([]*notageometry.Polygon, newCapacity)
	copy(newPolygons, em.polygons)
	em.polygons = newPolygons

	newColliders := make([]notacollision.Collider, newCapacity)
	copy(newColliders, em.colliders)
	em.colliders = newColliders

	newColors := make([]notacolor.Color, newCapacity)
	copy(newColors, em.colors)
	em.colors = newColors

	newShaders := make([]*notashader.Shader, newCapacity)
	copy(newShaders, em.shaders)
	em.shaders = newShaders

	newMaterials := make([]*notashader.Material, newCapacity)
	copy(newMaterials, em.materials)
	em.materials = newMaterials

	newActive := make([]bool, newCapacity)
	copy(newActive, em.active)
	em.active = newActive

	newVisible := make([]bool, newCapacity)
	copy(newVisible, em.visible)
	em.visible = newVisible

	newPendingMove := make([]notamath.Vec2, newCapacity)
	copy(newPendingMove, em.pendingMove)
	em.pendingMove = newPendingMove

	newPendingRot := make([]float32, newCapacity)
	copy(newPendingRot, em.pendingRot)
	em.pendingRot = newPendingRot

	newPendingScale := make([]notamath.Vec2, newCapacity)
	copy(newPendingScale, em.pendingScale)
	em.pendingScale = newPendingScale

	newHasPendingMove := make([]bool, newCapacity)
	copy(newHasPendingMove, em.hasPendingMove)
	em.hasPendingMove = newHasPendingMove

	newHasPendingRot := make([]bool, newCapacity)
	copy(newHasPendingRot, em.hasPendingRot)
	em.hasPendingRot = newHasPendingRot

	newHasPendingScale := make([]bool, newCapacity)
	copy(newHasPendingScale, em.hasPendingScale)
	em.hasPendingScale = newHasPendingScale

	// Initialize new slots with default values
	for i := len(em.entities); i < int(newCapacity); i++ {
		em.transforms[i] = TransformComponent{
			Scale: notamath.Vec2{X: 1, Y: 1},
		}
		em.colors[i] = notacolor.White
		em.active[i] = true
		em.visible[i] = true
		em.pendingScale[i] = notamath.Vec2{X: 1, Y: 1}
	}

	newEntities := make([]Entity, newCapacity)
	copy(newEntities, em.entities)
	em.entities = newEntities
}

// CreateEntity creates a new entity with the given name.
func (em *EntityManager) CreateEntity(name string) *Entity {
	em.mu.Lock()
	defer em.mu.Unlock()

	var id EntityID

	// Reuse a free ID or allocate a new one
	if len(em.freeIDs) > 0 {
		id = em.freeIDs[len(em.freeIDs)-1]
		em.freeIDs = em.freeIDs[:len(em.freeIDs)-1]
	} else {
		id = em.nextID
		em.nextID++
	}

	em.ensureCapacity(id)

	// Register name mapping
	em.nameToID[name] = id
	em.idToName[id] = name

	// Create entity handle
	em.entities[id] = Entity{
		ID:      id,
		manager: em,
	}

	// Set default archetype (Transform + Color)
	em.archetypes[id] = ArchetypeTransform | ArchetypeColor

	// Reset components to defaults
	em.transforms[id] = TransformComponent{
		Scale: notamath.Vec2{X: 1, Y: 1},
	}
	em.colors[id] = notacolor.White
	em.active[id] = true
	em.visible[id] = true
	em.sprites[id] = nil
	em.polygons[id] = nil
	em.colliders[id] = nil
	em.shaders[id] = nil
	em.materials[id] = nil

	return &em.entities[id]
}

// SubmitMove queues a movement delta for the entity.
func (em *EntityManager) SubmitMove(id EntityID, delta notamath.Vec2) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if int(id) >= len(em.pendingMove) {
		return
	}

	em.pendingMove[id] = em.pendingMove[id].Add(delta)
	em.hasPendingMove[id] = true
	em.dirtyMove = true
}

// SubmitRotation queues a rotation delta for the entity (in radians).
func (em *EntityManager) SubmitRotation(id EntityID, rad float32) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if int(id) >= len(em.pendingRot) {
		return
	}

	em.pendingRot[id] += rad
	em.hasPendingRot[id] = true
	em.dirtyRot = true
}

// SubmitScale queues a scale factor for the entity (multiplicative).
func (em *EntityManager) SubmitScale(id EntityID, factor notamath.Vec2) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if int(id) >= len(em.pendingScale) {
		return
	}

	em.pendingScale[id] = notamath.Vec2{
		X: em.pendingScale[id].X * factor.X,
		Y: em.pendingScale[id].Y * factor.Y,
	}
	em.hasPendingScale[id] = true
	em.dirtyScale = true
}

// Flush applies all pending transform updates and syncs colliders.
func (em *EntityManager) Flush() {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.flushEntities()
	em.flushColliders()
}

// flushEntities applies pending transform deltas.
func (em *EntityManager) flushEntities() {
	if em.dirtyMove {
		for i := EntityID(0); i < em.nextID; i++ {
			if em.hasPendingMove[i] {
				em.transforms[i].Position = em.transforms[i].Position.Add(em.pendingMove[i])
				em.pendingMove[i] = notamath.Vec2{}
				em.hasPendingMove[i] = false
			}
		}
		em.dirtyMove = false
	}

	if em.dirtyRot {
		for i := EntityID(0); i < em.nextID; i++ {
			if em.hasPendingRot[i] {
				em.transforms[i].Rotation += em.pendingRot[i]
				em.pendingRot[i] = 0
				em.hasPendingRot[i] = false
			}
		}
		em.dirtyRot = false
	}

	if em.dirtyScale {
		for i := EntityID(0); i < em.nextID; i++ {
			if em.hasPendingScale[i] {
				em.transforms[i].Scale = notamath.Vec2{
					X: em.transforms[i].Scale.X * em.pendingScale[i].X,
					Y: em.transforms[i].Scale.Y * em.pendingScale[i].Y,
				}
				em.pendingScale[i] = notamath.Vec2{X: 1, Y: 1}
				em.hasPendingScale[i] = false
			}
		}
		em.dirtyScale = false
	}
}

// flushColliders updates all colliders to match their entity's transform.
func (em *EntityManager) flushColliders() {
	em.syncColliders()

	// Clear collision table
	em.collisionTable = &collisionTable{
		pairs: make(map[EntityID]map[EntityID]notamath.Vec2),
	}
}

// syncColliders updates collider transforms.
func (em *EntityManager) syncColliders() {
	for i := EntityID(1); i < em.nextID; i++ {
		if !em.active[i] {
			continue
		}

		collider := em.colliders[i]
		if collider == nil {
			continue
		}

		t := notamath.Transform2D{}
		t.SetPosition(em.transforms[i].Position)
		t.SetRotation(em.transforms[i].Rotation)
		t.SetScale(em.transforms[i].Scale)

		collider.UpdateFromTransform(&t)
	}
}

// GetPosition returns the position of an entity.
func (em *EntityManager) GetPosition(id string) notamath.Vec2 {
	em.mu.RLock()
	defer em.mu.RUnlock()

	entityID, ok := em.nameToID[id]
	if !ok {
		return notamath.Vec2{}
	}

	return em.transforms[entityID].Position
}

// GetScale returns the scale of an entity.
func (em *EntityManager) GetScale(id string) notamath.Vec2 {
	em.mu.RLock()
	defer em.mu.RUnlock()

	entityID, ok := em.nameToID[id]
	if !ok {
		return notamath.Vec2{}
	}

	return em.transforms[entityID].Scale
}

// GetRotation returns the rotation of an entity (in radians).
func (em *EntityManager) GetRotation(id string) float32 {
	em.mu.RLock()
	defer em.mu.RUnlock()

	entityID, ok := em.nameToID[id]
	if !ok {
		return 0
	}

	return em.transforms[entityID].Rotation
}

// getPosition returns the position of an entity by ID.
func (em *EntityManager) getPosition(id EntityID) notamath.Vec2 {
	if int(id) >= len(em.transforms) {
		return notamath.Vec2{}
	}
	return em.transforms[id].Position
}

// getScale returns the scale of an entity by ID.
func (em *EntityManager) getScale(id EntityID) notamath.Vec2 {
	if int(id) >= len(em.transforms) {
		return notamath.Vec2{}
	}
	return em.transforms[id].Scale
}

// getRotation returns the rotation of an entity by ID.
func (em *EntityManager) getRotation(id EntityID) float32 {
	if int(id) >= len(em.transforms) {
		return 0
	}
	return em.transforms[id].Rotation
}

// GetEntities returns all entity handles.
func (em *EntityManager) GetEntities() []*Entity {
	em.mu.RLock()
	defer em.mu.RUnlock()

	result := make([]*Entity, 0, em.nextID)
	for i := EntityID(1); i < em.nextID; i++ {
		if em.active[i] {
			result = append(result, &em.entities[i])
		}
	}
	return result
}

// GetEntity returns the entity with the given name.
func (em *EntityManager) GetEntity(name string) *Entity {
	em.mu.RLock()
	defer em.mu.RUnlock()

	id, ok := em.nameToID[name]
	if !ok {
		return nil
	}

	if int(id) >= len(em.entities) || !em.active[id] {
		return nil
	}

	return &em.entities[id]
}

// Remove removes an entity by name.
func (em *EntityManager) Remove(name string) {
	em.mu.Lock()
	defer em.mu.Unlock()

	id, ok := em.nameToID[name]
	if !ok {
		return
	}

	em.removeEntity(id)
}

// removeEntity removes an entity by ID (must hold lock).
func (em *EntityManager) removeEntity(id EntityID) {
	delete(em.nameToID, em.idToName[id])
	delete(em.idToName, id)
	delete(em.archetypes, id)

	em.active[id] = false
	em.visible[id] = false
	em.freeIDs = append(em.freeIDs, id)

	// Remove from collision groups
	for _, group := range em.collisionGroups {
		filtered := make([]EntityID, 0, len(group.Entities))
		for _, eid := range group.Entities {
			if eid != id {
				filtered = append(filtered, eid)
			}
		}
		group.Entities = filtered
	}
}

// AddToCollisionGroup adds an entity to a collision group.
func (em *EntityManager) AddToCollisionGroup(group string, e *Entity) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if _, ok := em.collisionGroups[group]; !ok {
		em.collisionGroups[group] = &CollisionGroup{
			Name:     group,
			Entities: []EntityID{},
		}
	}

	em.collisionGroups[group].Entities = append(em.collisionGroups[group].Entities, e.ID)
}

// SolveGroupCollision computes collisions between entities in a group.
func (em *EntityManager) SolveGroupCollision(name string) {
	em.mu.Lock()
	defer em.mu.Unlock()

	group, ok := em.collisionGroups[name]
	if !ok || len(group.Entities) < 2 {
		return
	}

	entities := group.Entities

	for a := 0; a < len(entities); a++ {
		id1 := entities[a]
		if !em.active[id1] {
			continue
		}

		c1 := em.colliders[id1]
		if c1 == nil {
			continue
		}

		for b := a + 1; b < len(entities); b++ {
			id2 := entities[b]
			if !em.active[id2] {
				continue
			}

			c2 := em.colliders[id2]
			if c2 == nil {
				continue
			}

			_, mtv := notacollision.Intersects(c1, c2)
			if mtv != (notamath.Vec2{}) {
				if em.collisionTable.pairs[id1] == nil {
					em.collisionTable.pairs[id1] = make(map[EntityID]notamath.Vec2)
				}
				if em.collisionTable.pairs[id2] == nil {
					em.collisionTable.pairs[id2] = make(map[EntityID]notamath.Vec2)
				}
				em.collisionTable.pairs[id1][id2] = mtv
				em.collisionTable.pairs[id2][id1] = mtv.Neg()
			}
		}
	}
}

// Collides checks if two entities are colliding.
func (em *EntityManager) Collides(a, b *Entity) bool {
	collides, _ := em.CollidesMTV(a, b)
	return collides
}

// GetMTV returns the minimum translation vector between two entities.
func (em *EntityManager) GetMTV(a, b *Entity) notamath.Vec2 {
	_, mtv := em.CollidesMTV(a, b)
	return mtv
}

// CollidesMTV checks collision and returns the MTV.
func (em *EntityManager) CollidesMTV(a, b *Entity) (bool, notamath.Vec2) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	row := em.collisionTable.pairs[a.ID]
	if row == nil {
		return false, notamath.Vec2{}
	}

	mtv, ok := row[b.ID]
	return ok, mtv
}

// Component setters

// SetSprite sets the sprite component for an entity.
func (em *EntityManager) SetSprite(id EntityID, sprite *notatexture.Sprite) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.ensureCapacity(id)
	em.sprites[id] = sprite
}

// SetPolygon sets the polygon component for an entity.
func (em *EntityManager) SetPolygon(id EntityID, polygon *notageometry.Polygon) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.ensureCapacity(id)
	em.polygons[id] = polygon
}

// SetCollider sets the collider component for an entity.
func (em *EntityManager) SetCollider(id EntityID, collider notacollision.Collider) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.ensureCapacity(id)
	em.colliders[id] = collider
}

// SetColor sets the color component for an entity.
func (em *EntityManager) SetColor(id EntityID, color notacolor.Color) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.ensureCapacity(id)
	em.colors[id] = color
}

// SetShader sets the shader component for an entity.
func (em *EntityManager) SetShader(id EntityID, shader *notashader.Shader) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.ensureCapacity(id)
	em.shaders[id] = shader
}

// SetMaterial sets the material component for an entity.
func (em *EntityManager) SetMaterial(id EntityID, material *notashader.Material) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.ensureCapacity(id)
	em.materials[id] = material
}

// Component getters

// GetSprite returns the sprite component for an entity.
func (em *EntityManager) GetSprite(id EntityID) *notatexture.Sprite {
	em.mu.RLock()
	defer em.mu.RUnlock()

	if int(id) >= len(em.sprites) {
		return nil
	}
	return em.sprites[id]
}

// GetPolygon returns the polygon component for an entity.
func (em *EntityManager) GetPolygon(id EntityID) *notageometry.Polygon {
	em.mu.RLock()
	defer em.mu.RUnlock()

	if int(id) >= len(em.polygons) {
		return nil
	}
	return em.polygons[id]
}

// GetColor returns the color component for an entity.
func (em *EntityManager) GetColor(id EntityID) notacolor.Color {
	em.mu.RLock()
	defer em.mu.RUnlock()

	if int(id) >= len(em.colors) {
		return notacolor.White
	}
	return em.colors[id]
}

// GetShader returns the shader component for an entity.
func (em *EntityManager) GetShader(id EntityID) *notashader.Shader {
	em.mu.RLock()
	defer em.mu.RUnlock()

	if int(id) >= len(em.shaders) {
		return nil
	}
	return em.shaders[id]
}

// GetMaterial returns the material component for an entity.
func (em *EntityManager) GetMaterial(id EntityID) *notashader.Material {
	em.mu.RLock()
	defer em.mu.RUnlock()

	if int(id) >= len(em.materials) {
		return nil
	}
	return em.materials[id]
}

// IsActive returns whether an entity is active.
func (em *EntityManager) IsActive(id EntityID) bool {
	em.mu.RLock()
	defer em.mu.RUnlock()

	if int(id) >= len(em.active) {
		return false
	}
	return em.active[id]
}

// SetActive sets the active state of an entity.
func (em *EntityManager) SetActive(id EntityID, active bool) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if int(id) < len(em.active) {
		em.active[id] = active
	}
}

// IsVisible returns whether an entity is visible.
func (em *EntityManager) IsVisible(id EntityID) bool {
	em.mu.RLock()
	defer em.mu.RUnlock()

	if int(id) >= len(em.visible) {
		return false
	}
	return em.visible[id]
}

// SetVisible sets the visibility state of an entity.
func (em *EntityManager) SetVisible(id EntityID, visible bool) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if int(id) < len(em.visible) {
		em.visible[id] = visible
	}
}

// GetEntityName returns the name of an entity by ID.
func (em *EntityManager) GetEntityName(id EntityID) string {
	em.mu.RLock()
	defer em.mu.RUnlock()

	return em.idToName[id]
}

// SetPosition sets the position of an entity directly.
func (em *EntityManager) SetPosition(id EntityID, pos notamath.Vec2) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.ensureCapacity(id)
	em.transforms[id].Position = pos
}

// SetRotation sets the rotation of an entity directly.
func (em *EntityManager) SetRotation(id EntityID, rot float32) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.ensureCapacity(id)
	em.transforms[id].Rotation = rot
}

// SetScaleValue sets the scale of an entity directly.
func (em *EntityManager) SetScaleValue(id EntityID, scale notamath.Vec2) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.ensureCapacity(id)
	em.transforms[id].Scale = scale
}

// ForEachWithArchetype iterates over all entities that have the specified archetype.
func (em *EntityManager) ForEachWithArchetype(archetype Archetype, fn func(EntityID)) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	for id, a := range em.archetypes {
		if (a&archetype) == archetype && em.active[id] {
			fn(id)
		}
	}
}

// Unused but required for interface compatibility
var _ = vek32.Add_Inplace
