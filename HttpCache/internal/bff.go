package internal

import (
	"log"
	"sync"
)

/*
@Time    : 2021/3/3 22:31
@Author  : austsxk
@Email   : austsxk@163.com
@File    : bff.go
@Software: GoLand
*/

// defined union cache interface
type Cache interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	GetStat() Stat
}

// cache k-v pair
type Stat struct {
	Count     int64
	KeySize   int64
	ValueSize int64
}

// cache key then add key size and value size
func (s *Stat) Add(key, value string) {
	s.Count++
	s.KeySize += int64(len([]byte(key)))
	s.ValueSize += int64(len([]byte(value)))
}

// delete key then reduce key size and value size
func (s *Stat) Del(key, value string) {
	s.Count--
	s.KeySize -= int64(len([]byte(key)))
	s.ValueSize -= int64(len([]byte(value)))
}

// create a struct to implement cache interface
type MemoryCacheImpl struct {
	// because go raw map not security, so use sync.RWMutex  ensure safety
	mut     sync.RWMutex
	hashMap map[string]string
	// embed Stat struct
	Stat
}

func (m *MemoryCacheImpl) Set(key, value string) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	if v, ok := m.hashMap[key]; ok {
		m.Del(key, v)
	}
	m.hashMap[key] = value
	m.Add(key, value)
	return nil
}

func (m *MemoryCacheImpl) Get(key string) (string, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.hashMap[key], nil
}

func (m *MemoryCacheImpl) Delete(key string) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	if v, ok := m.hashMap[key]; ok {
		delete(m.hashMap, key)
		m.Del(key, v)
	}
	return nil
}

func (m *MemoryCacheImpl) GetStat() Stat {
	return m.Stat
}

// build func to new MemoryCacheImpl
func NewMemoryCacheImpl() *MemoryCacheImpl {
	return &MemoryCacheImpl{
		sync.RWMutex{},
		make(map[string]string),
		Stat{}}
}

func Make(t string) Cache {
	var c Cache
	if t == "cache" {
		c = NewMemoryCacheImpl()
	}
	if c == nil {
		log.Fatal("unknown cache type: ", t)
	}
	log.Println("cache now is already!")
	return c
}
