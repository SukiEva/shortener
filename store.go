package shorturl

import (
	"sync"
)

type URLStore struct {
	urls map[string]string
	mu   sync.RWMutex // 读写锁，保证线程安全
}

// NewURLStore 工厂函数
func NewURLStore() *URLStore {
	return &URLStore{ // & 取地址，即变为指针
		urls: make(map[string]string),
	}
}

func (s *URLStore) Get(key string) string {
	s.mu.RLock() // 防止读-写冲突
	defer s.mu.RUnlock()
	return s.urls[key]
}

func (s *URLStore) Set(key, url string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.urls[key]; ok {
		return false
	}
	s.urls[key] = url
	return true
}

func (s *URLStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.urls)
}

func (s *URLStore) Put(url string) string {
	for {
		key := genKey(s.Count())
		if s.Set(key, url) {
			return key
		}
	}
	//return ""
}
