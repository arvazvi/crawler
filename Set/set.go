package set

import (
	"hash/fnv"
	"sync"
)

type HashSet struct {
	hashmap map[uint64]bool
	size int
	mu sync.Mutex
}

func (s *HashSet) Insert(value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.hashmap[Hash(value)] = true
	s.size++
}

func (s *HashSet) Contains(value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hashmap[Hash(value)]
}

func (s *HashSet) Size(value string) int {
	return s.size
}

func Hash(url string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(url))

	return h.Sum64()
}