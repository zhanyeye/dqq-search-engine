package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

// 初始化测试用的 ConcurrentHashMap 和 sync.Map
var conMp = NewConcurrentHashMap(8, 1000)
var synMp = sync.Map{}

// readConMap 从 ConcurrentHashMap 中读取数据
func readConMap() {
	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		conMp.Get(key)
	}
}

// writeConMap 向 ConcurrentHashMap 中写入数据
func writeConMap() {
	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		conMp.Set(key, 1)
	}
}

// readSynMap 从 sync.Map 中读取数据
func readSynMap() {
	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		synMp.Load(key)
	}
}

// writeSynMap 向 sync.Map 中写入数据
func writeSynMap() {
	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		synMp.Store(key, 1)
	}
}

// BenchmarkConMap 并发读写测试 ConcurrentHashMap
func BenchmarkConMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const P = 300
		wg := sync.WaitGroup{}
		wg.Add(2 * P)
		for i := 0; i < P; i++ { // 300 个协程一直读
			go func() {
				defer wg.Done()
				readConMap()
			}()
		}
		for i := 0; i < P; i++ { // 300 个协程一直写
			go func() {
				defer wg.Done()
				writeConMap()
				// time.Sleep(100 * time.Millisecond)   //写很少时速度差1倍，一直写时速度差3倍
			}()
		}
		wg.Wait()
	}
}

// BenchmarkSynMap 并发读写测试 sync.Map
func BenchmarkSynMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const P = 300
		wg := sync.WaitGroup{}
		wg.Add(2 * P)
		for i := 0; i < P; i++ { // 300 个协程一直读
			go func() {
				defer wg.Done()
				readSynMap()
			}()
		}
		for i := 0; i < P; i++ { // 300 个协程一直写
			go func() {
				defer wg.Done()
				writeSynMap()
				// time.Sleep(100 * time.Millisecond)
			}()
		}
		wg.Wait()
	}
}

// TestConcurrentHashMap 基本功能测试
func TestConcurrentHashMap(t *testing.T) {
	m := NewConcurrentHashMap(10, 100)

	// 测试 Set 和 Get 方法
	m.Set("key1", "value1")
	value, exists := m.Get("key1")
	if !exists || value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}

	// 测试不存在的 key
	value, exists = m.Get("key2")
	if exists || value != nil {
		t.Errorf("Expected nil, got %v", value)
	}

	// 测试迭代器
	m.Set("key2", "value2")
	iterator := m.CreateIterator()
	count := 0
	for entry := iterator.Next(); entry != nil; entry = iterator.Next() {
		count++
	}
	if count != 2 {
		t.Errorf("Expected 2 entries, got %d", count)
	}
}

// TestConcurrentHashMapIterator 测试 ConcurrentHashMap 的迭代器
func TestConcurrentHashMapIterator(t *testing.T) {
	for i := 0; i < 10; i++ {
		conMp.Set(strconv.Itoa(i), i)
	}
	iterator := conMp.CreateIterator()
	entry := iterator.Next()
	for entry != nil {
		fmt.Println(entry.Key, entry.Value)
		entry = iterator.Next()
	}
}

// TestConcurrentHashMapConcurrency 并发读写测试
func TestConcurrentHashMapConcurrency(t *testing.T) {
	const P = 100
	var wg sync.WaitGroup
	wg.Add(2 * P)

	// 用于验证写入的数据是否能被读取
	keys := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		keys[i] = strconv.Itoa(int(rand.Int63()))
		conMp.Set(keys[i], i)
	}

	// 启动 P 个协程进行读操作
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				key := keys[j%1000]
				value, exists := conMp.Get(key)
				if !exists || value != j%1000 {
					t.Errorf("Expected %d, got %v", j%1000, value)
				}
			}
		}()
	}

	// 启动 P 个协程进行写操作
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				key := strconv.Itoa(int(rand.Int63()))
				conMp.Set(key, rand.Int())
			}
		}()
	}

	wg.Wait()
}

// go test ./utils -bench=Map -run=^$ -count=1 -benchmem -benchtime=3s
/*
goos: windows
goarch: amd64
pkg: dqq-search-engine/utils
cpu: 13th Gen Intel(R) Core(TM) i7-1360P
BenchmarkConMap-16             6         613649983 ns/op        522988856 B/op   6089069 allocs/op
BenchmarkSynMap-16             2        4505936150 ns/op        911848880 B/op  18099574 allocs/op

*/
