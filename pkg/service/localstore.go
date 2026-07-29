// Package service handles room allocation, persistent storage, and streaming I/O endpoints.
package service

import (
	"context"
	"errors"
	"sync"

	"github.com/agenticsfu/agentic-sfu/pkg/rtc"
)

var ErrRoomNotFoundInStore = errors.New("room not found in store")

// LocalStore implements an in-memory session persistence store.
type LocalStore struct {
	mu    sync.RWMutex
	rooms map[string]*rtc.Room
}

func NewLocalStore() *LocalStore {
	return &LocalStore{
		rooms: make(map[string]*rtc.Room),
	}
}

func (s *LocalStore) StoreRoom(_ context.Context, room *rtc.Room) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rooms[room.Name()] = room
	return nil
}

func (s *LocalStore) LoadRoom(_ context.Context, name string) (*rtc.Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	room, ok := s.rooms[name]
	if !ok {
		return nil, ErrRoomNotFoundInStore
	}
	return room, nil
}

func (s *LocalStore) DeleteRoom(_ context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.rooms, name)
	return nil
}

func (s *LocalStore) ListRooms(_ context.Context) ([]*rtc.Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*rtc.Room, 0, len(s.rooms))
	for _, r := range s.rooms {
		list = append(list, r)
	}
	return list, nil
}
