package notaentity

// SparseSet implements a sparse-set/dense-set data structure for efficient entity storage.
// It provides O(1) lookup, insertion, and removal operations.
// The sparse array maps entity IDs to indices in the dense array.
// The dense array stores entity data contiguously for cache-friendly iteration.
type SparseSet struct {
	sparse []uint32 // sparse[entityID] = index in dense array
	dense  []uint32 // dense[index] = entityID
	size   uint32   // number of active entities
}

// NewSparseSet creates a new sparse set with the given initial capacity.
func NewSparseSet(capacity uint32) *SparseSet {
	return &SparseSet{
		sparse: make([]uint32, capacity),
		dense:  make([]uint32, capacity),
		size:   0,
	}
}

// Contains returns true if the entity exists in the set.
func (s *SparseSet) Contains(entity EntityID) bool {
	id := uint32(entity)
	return id < uint32(len(s.sparse)) && s.sparse[id] < s.size && s.dense[s.sparse[id]] == id
}

// Add adds an entity to the set. Returns true if the entity was newly added.
func (s *SparseSet) Add(entity EntityID) bool {
	id := uint32(entity)

	if s.Contains(entity) {
		return false
	}

	// Grow sparse array if needed
	if id >= uint32(len(s.sparse)) {
		newSparse := make([]uint32, max(id*2+1, uint32(len(s.sparse)*2)))
		copy(newSparse, s.sparse)
		s.sparse = newSparse
	}

	// Grow dense array if needed
	if s.size >= uint32(len(s.dense)) {
		newDense := make([]uint32, max(uint32(len(s.dense)*2), 64))
		copy(newDense, s.dense)
		s.dense = newDense
	}

	s.sparse[id] = s.size
	s.dense[s.size] = id
	s.size++

	return true
}

// Remove removes an entity from the set. Returns true if the entity was removed.
func (s *SparseSet) Remove(entity EntityID) bool {
	id := uint32(entity)

	if !s.Contains(entity) {
		return false
	}

	// Get the index of this entity in the dense array
	idx := s.sparse[id]
	s.size--

	// If this is not the last element, swap with the last element
	if idx != s.size {
		lastID := s.dense[s.size]
		s.dense[idx] = lastID
		s.sparse[lastID] = idx
	}

	return true
}

// GetEntities returns a slice of all entities in the set (dense array portion).
// This is efficient for iteration as entities are stored contiguously.
func (s *SparseSet) GetEntities() []uint32 {
	return s.dense[:s.size]
}

// Size returns the number of entities in the set.
func (s *SparseSet) Size() uint32 {
	return s.size
}

// Clear removes all entities from the set.
func (s *SparseSet) Clear() {
	s.size = 0
}
