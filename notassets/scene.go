package notassets

import (
	"fmt"
	"sync"
)

type EntityManager struct {
	Name     string
	Entities map[string]*Entity
	mu       sync.RWMutex
}

func NewScene(name string) *EntityManager {
	return &EntityManager{
		Name:     name,
		Entities: make(map[string]*Entity),
	}
}

// Add adds an entity to the scene
func (s *EntityManager) Add(entity *Entity) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entity == nil {
		return fmt.Errorf("cannot add nil entity")
	}

	if entity.ID == "" {
		return fmt.Errorf("entity must have an ID")
	}

	if _, exists := s.Entities[entity.ID]; exists {
		return fmt.Errorf("entity with ID '%s' already exists in scene", entity.ID)
	}

	s.Entities[entity.ID] = entity
	return nil
}

// Remove removes an entity from the scene
func (s *EntityManager) Remove(entityID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Entities[entityID]; !exists {
		return fmt.Errorf("entity with ID '%s' not found in scene", entityID)
	}

	delete(s.Entities, entityID)
	return nil
}

// Get retrieves an entity by ID
func (s *EntityManager) Get(entityID string) (*Entity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entity, exists := s.Entities[entityID]
	if !exists {
		return nil, fmt.Errorf("entity with ID '%s' not found in scene", entityID)
	}

	return entity, nil
}

// Count returns the number of entities in the scene
func (s *EntityManager) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.Entities)
}

// Clear removes all entities from the scene
func (s *EntityManager) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Entities = make(map[string]*Entity)
}

// GetActiveEntities returns all active entities
func (s *EntityManager) GetActiveEntities() []*Entity {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var active []*Entity
	for _, entity := range s.Entities {
		if entity.Active {
			active = append(active, entity)
		}
	}

	return active
}

// GetVisibleEntities returns all visible entities
func (s *EntityManager) GetVisibleEntities() []*Entity {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var visible []*Entity
	for _, entity := range s.Entities {
		if entity.Visible {
			visible = append(visible, entity)
		}
	}

	return visible
}
