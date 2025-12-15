package server

import "sync"

type Store struct {
	mu        sync.RWMutex
	documents map[string]*Document
}

func NewStore() *Store {
	return &Store{
		documents: make(map[string]*Document),
	}
}

func (s *Store) CreateDocument(id string) *Document {
	s.mu.Lock()
	defer s.mu.Unlock()

	doc := &Document{ID: id}
	s.documents[id] = doc
	return doc
}

func (s *Store) GetDocument(id string) (*Document, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doc, ok := s.documents[id]
	return doc, ok
}
