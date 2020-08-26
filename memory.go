// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:08 PM - UTC/GMT+08:00

package session

import (
	"sync"
	"time"
)

// 垃圾回收
type garbage struct {
	id     string     // session id
	expire *time.Time // expire time
}

// MemoryStore 内存存储实现
type MemoryStore struct {
	// lock When parallel, ensure data independence, consistency and safety
	mx sync.Mutex
	// sid:key:data save serialize data
	values map[string]map[string][]byte
	// garbage of session
	garbageList []*garbage
}

// newMemoryStore 创建一个内存存储 开辟内存
func newMemoryStore() *MemoryStore {
	ms := &MemoryStore{values: make(map[string]map[string][]byte, MemoryMaxSize)}
	//ms.values[""] = make(map[string]interface{},maxSize)
	// init GARBAGE
	ms.garbageList = make([]*garbage, 0, MemoryMaxSize)
	go ms.gc()
	return ms
}

// Writer 写入数据方法
func (m *MemoryStore) Writer(id, key string, data interface{}) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	// check map pointer is exist
	if m.values[id] == nil {
		m.values[id] = make(map[string][]byte, maxSize)
	}
	serialize, err := Serialize(data)
	if err != nil {
		return err
	}
	m.values[id][key] = serialize
	//log.Printf("%p",m.values[id])
	//log.Println(m.values[id][key])
	return nil

}

// Reader 读取数据 通过id和key
func (m *MemoryStore) Reader(id, key string) ([]byte, error) {
	return m.values[id][key], nil
}

// Remove 通过id和key移除数据
func (m *MemoryStore) Remove(id, key string) {
	delete(m.values[id], key)

}

// Clean 通过id清空data
func (m *MemoryStore) Clean(id string) {
	m.values[id] = make(map[string][]byte, maxSize)
}

func (m *MemoryStore) garbage(g *garbage) {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.garbageList = append(m.garbageList, g)
}

// gc GarbageCollection
func (m *MemoryStore) gc() {
	// 每10分钟进行一次垃圾清理  session过期的全部清理掉
	for {
		time.Sleep(10 * time.Second)

	}

}
