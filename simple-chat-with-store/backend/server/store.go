package server

import (
	"sync"
	"time"
)

// Very small in-memory store. Replace with DB if needed.

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type Group struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Users []string `json:"users"`
}

type Store struct {
	mu     sync.RWMutex
	users  map[string]User
	groups map[string]Group
}

func NewStore() *Store {
	return &Store{users: make(map[string]User), groups: make(map[string]Group)}
}

func (s *Store) AddUser(u User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[u.ID] = u
}

func (s *Store) CreateUser(u User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[u.ID] = u
}

func (s *Store) GetUser(id string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	return u, ok
}

func (s *Store) AddGroup(g Group) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.groups[g.ID] = g
}

func (s *Store) GetGroup(id string) (Group, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	g, ok := s.groups[id]
	return g, ok
}

func (s *Store) ListGroups() []Group {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Group, 0, len(s.groups))
	for _, g := range s.groups {
		out = append(out, g)
	}
	return out
}
