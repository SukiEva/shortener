package shorturl

import (
	"encoding/gob"
	"io"
	"os"
	"sync"
)

type URLStore struct {
	urls map[string]string
	mu   sync.RWMutex // 读写锁，保证线程安全
	file *os.File
}

type record struct {
	Key, Url string
}

// NewURLStore 工厂函数
func NewURLStore(fileName string) *URLStore {
	s := &URLStore{ // & 取地址，即变为指针
		urls: make(map[string]string),
	}
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	DropError(err, "Error opening URLStore:")
	s.file = f
	err = s.load()
	DropError(err, "Error loading URLStore:")
	return s
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
			err := s.save(key, url)
			DropError(err, "Error saving to URLStore:")
			return key
		}
	}
	//return ""
}

// 解码保存的文件，并存入内存
func (s *URLStore) load() error {
	if _, err := s.file.Seek(0, 0); err != nil { // 寻址文件起始位置
		return err
	}
	d := gob.NewDecoder(s.file)
	var err error
	for err == nil { // 循环解码
		var r record // 记录
		if err = d.Decode(&r); err == nil {
			s.Set(r.Key, r.Url)
		}
	}
	if err == io.EOF { // 解码成功
		return nil
	}
	return err
}

// 编码后保存至文件
func (s *URLStore) save(key, url string) error {
	e := gob.NewEncoder(s.file)
	return e.Encode(record{key, url})
}
