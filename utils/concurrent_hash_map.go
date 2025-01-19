package util

import (
	farmhash "github.com/leemcloughlin/gofarmhash"
	"golang.org/x/exp/maps"
	"sync"
)

// ConcurrentHashMap 支持并发读写的哈希表。key是string，value是any
type ConcurrentHashMap struct {
	segments     []map[string]any //由多个小map构成,
	segmentCount int              //小map的个数
	segmentLocks []sync.RWMutex   //每个小map配一把读写锁。避免全局只有一把锁，影响性能
	seed         uint32           //每次执行farmhash传统一的seed
}

// NewConcurrentHashMap 创建一个新的ConcurrentHashMap
// capacity 预估map中容纳多少元素，segmentCount 内部包含几个小map
func NewConcurrentHashMap(segmentCount, capacity int) *ConcurrentHashMap {
	segments := make([]map[string]any, segmentCount)
	segmentLocks := make([]sync.RWMutex, segmentCount)
	for i := 0; i < segmentCount; i++ {
		segments[i] = make(map[string]any, capacity/segmentCount)
	}
	return &ConcurrentHashMap{
		segments:     segments,
		segmentCount: segmentCount,
		seed:         0,
		segmentLocks: segmentLocks,
	}
}

// getSegmentIndex 判断key对应到哪个小map
func (m *ConcurrentHashMap) getSegmentIndex(key string) int {
	hash := int(farmhash.Hash32WithSeed([]byte(key), m.seed)) // FarmHash是google开源的Hash算法
	return hash % m.segmentCount
}

// Set 写入<key, value>
// todo：注释解释为啥用 Lock() 锁
func (m *ConcurrentHashMap) Set(key string, value any) {
	index := m.getSegmentIndex(key)
	m.segmentLocks[index].Lock()
	defer m.segmentLocks[index].Unlock()
	m.segments[index][key] = value
}

// Get 根据key读取value
// todo：注释解释为啥用 RLock() 锁
func (m *ConcurrentHashMap) Get(key string) (any, bool) {
	index := m.getSegmentIndex(key)
	m.segmentLocks[index].RLock()
	defer m.segmentLocks[index].RUnlock()
	value, exists := m.segments[index][key]
	return value, exists
}

// Delete 删除指定key的元素
func (m *ConcurrentHashMap) Delete(key string) {
    index := m.getSegmentIndex(key)
    m.segmentLocks[index].Lock()
    defer m.segmentLocks[index].Unlock()
    delete(m.segments[index], key)
}

// CreateIterator 创建一个迭代器
func (m *ConcurrentHashMap) CreateIterator() *ConcurrentHashMapIterator {
	keys := make([][]string, 0, len(m.segments))
	for _, segment := range m.segments {
		row := maps.Keys(segment)
		keys = append(keys, row)
	}
	return &ConcurrentHashMapIterator{
		concurrentMap: m,
		keys:          keys,
		rowIndex:      0,
		colIndex:      0,
	}
}

// MapEntry 表示一个键值对
type MapEntry struct {
	Key   string
	Value any
}

// MapIterator 迭代器接口
type MapIterator interface {
	Next() *MapEntry
}

// ConcurrentHashMapIterator 并发哈希表的迭代器
type ConcurrentHashMapIterator struct {
	concurrentMap *ConcurrentHashMap
	keys          [][]string
	rowIndex      int
	colIndex      int
}

// Next 获取下一个键值对
// go标准库的container/list也是通过Next()来遍历，go标准库database/sql规定按Rows.Next()来遍历查询结果
func (iter *ConcurrentHashMapIterator) Next() *MapEntry {
	if iter.rowIndex >= len(iter.keys) {
		return nil
	}
	row := iter.keys[iter.rowIndex]
	if len(row) == 0 { // 本行为空
		iter.rowIndex += 1
		return iter.Next() // 进入递归，因为下一行可能依然为空
	}
	key := row[iter.colIndex] // 根据下标访问切片元素时，一定注意不要出现数组越界异常。即使下标为0，当切片为空时依然会出现数组越界异常
	value, _ := iter.concurrentMap.Get(key)
	if iter.colIndex >= len(row)-1 {
		iter.rowIndex += 1
		iter.colIndex = 0
	} else {
		iter.colIndex += 1
	}
	return &MapEntry{key, value}
}
