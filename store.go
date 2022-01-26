package shorturl

import (
	"encoding/json"
	"errors"
	"io"
	"net/rpc"
	"os"
	"sync"
)

const saveQueueLength = 1000

type Store interface {
	Put(url, key *string) error
	Get(key, url *string) error
}

type ProxyStore struct {
	urls   *URLStore // local cache
	client *rpc.Client
}

type URLStore struct {
	urls map[string]string
	mu   sync.RWMutex // 读写锁，保证线程安全
	save chan record
}

type record struct {
	Key, Url string
}

// NewURLStore 工厂函数
func NewURLStore(fileName string) *URLStore {
	s := &URLStore{ // & 取地址，即变为指针
		urls: make(map[string]string),
	}
	if fileName == "" {
		return s
	}
	s.save = make(chan record, saveQueueLength) // 带缓冲 channel
	err := s.load(fileName)
	DropError(err, "Error loading URLStore:")
	go s.saveLoop(fileName) // 保存文件协程
	return s
}

func (s *URLStore) Get(key, url *string) error {
	s.mu.RLock() // 防止读-写冲突
	defer s.mu.RUnlock()
	if u, ok := s.urls[*key]; ok {
		*url = u
		return nil
	}
	return errors.New("key not found")
}

func (s *URLStore) Set(key, url *string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.urls[*key]; ok {
		return errors.New("key already exists")
	}
	s.urls[*key] = *url
	return nil
}

func (s *URLStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.urls)
}

func (s *URLStore) Put(url, key *string) error {
	for {
		*key = genKey(s.Count())
		if err := s.Set(key, url); err == nil {
			break
		}
	}
	if s.save != nil {
		s.save <- record{*key, *url} // 发送到 channel
	}
	return nil
}

// 解码保存的文件，并载入内存
func (s *URLStore) load(fileName string) error {
	f, err := os.Open(fileName)
	//DropError(err, "Error opening URLStore:")
	if err != nil {
		return err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	for err == nil { // 循环解码
		var r record // 记录
		if err = d.Decode(&r); err == nil {
			s.Set(&r.Key, &r.Url)
		}
	}
	if err == io.EOF { // 解码成功
		return nil
	}
	//DropError(err, "Error decoding URLStore:")
	return err
}

// 从 channel 中获取记录并且编码到文件
func (s *URLStore) saveLoop(fileName string) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	DropError(err, "Error opening URLStore:")
	defer f.Close()
	e := json.NewEncoder(f)
	for {
		r := <-s.save // 从 channel 接收数据
		err = e.Encode(r)
		DropError(err, "Error saving to URLStore:")
	}
}

func NewProxyStore(addr string) *ProxyStore {
	client, err := rpc.DialHTTP("tcp", addr)
	DropError(err, "Error constructing ProxyStore:")
	return &ProxyStore{urls: NewURLStore(""), client: client}
}

func (s *ProxyStore) Get(key, url *string) error {
	if err := s.urls.Get(key, url); err == nil {
		return nil
	}
	// rpc call to master:
	if err := s.client.Call("Store.Get", key, url); err != nil {
		return err
	}
	s.urls.Set(key, url) // update local cache
	return nil
}

func (s *ProxyStore) Put(url, key *string) error {
	// rpc call to master:
	if err := s.client.Call("Store.Put", url, key); err != nil {
		return err
	}
	s.urls.Set(key, url) // update local cache
	return nil
}
