package bugs

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var ErrorNotFound = errors.New("not found")
var ErrorInvalidInput = errors.New("invalid input")

// Bugs datastore interface that can be implemented in memory, (no)SQL, etc.
type IBugService interface {
	// GetAllBugs returns all bugs as a slice
	GetAllBugs() ([]Bug, error)
	// GetBugByID returns a bug by its ID or error if not found
	GetBugByID(id uint64) (Bug, error)
	// CreateBug creates a new bug and returns it or error if failed
	CreateBug(bug Bug) (Bug, error)
	// UpdateBug updates an existing bug and returns it or error if failed
	UpdateBug(bug Bug) (Bug, error)
	// DeleteBug deletes a bug by its ID or error if not found
	DeleteBug(id uint64) error
}

// BugService is an in-memory implementation of IBugService
type BugService struct {
	bugs     []Bug
	_mutex   sync.Mutex
	_counter atomic.Uint64
}

// NewBugService creates a new in-memory BugService
// and returns it as IBugService
func NewBugService() IBugService {
	return &BugService{bugs: []Bug{}} // initialize to empty list
}

// in-memory implementation of IBugService :: GetAllBugs()
func (s *BugService) GetAllBugs() ([]Bug, error) {
	return s.bugs, nil
}

// in-memory implementation of IBugService :: GetBugByID()
func (s *BugService) GetBugByID(id uint64) (Bug, error) {
	for _, bug := range s.bugs {
		if bug.ID == id {
			return bug, nil
		}
	}
	return Bug{}, fmt.Errorf("bug with id %d not found, error: %w", id, ErrorNotFound)
}

func (s *BugService) _inc_counter() uint64 {
	s._mutex.Lock()
	defer s._mutex.Unlock()
	return s._counter.Add(1)
}

// in-memory implementation of IBugService :: CreateBug()
func (s *BugService) CreateBug(bug Bug) (Bug, error) {
	bug.ID = s._inc_counter()
	bug.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	bug.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	s.bugs = append(s.bugs, bug)
	return bug, nil
}

// in-memory implementation of IBugService :: UpdateBug()
func (s *BugService) UpdateBug(bug Bug) (Bug, error) {
	for i, b := range s.bugs {
		if b.ID == bug.ID {
			id := s.bugs[i].ID
			ca := s.bugs[i].CreatedAt
			s.bugs[i] = bug
			bug.ID = id
			bug.CreatedAt = ca
			bug.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
			return bug, nil
		}
	}
	return Bug{}, fmt.Errorf("bug with id %d not found, error: %w", bug.ID, ErrorNotFound)
}

// in-memory implementation of IBugService :: DeleteBug()
func (s *BugService) DeleteBug(id uint64) error {
	for i, bug := range s.bugs {
		if bug.ID == id {
			s.bugs = append(s.bugs[:i], s.bugs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("bug with id %d not found, error: %w", id, ErrorNotFound)
}
