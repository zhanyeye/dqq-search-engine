package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

// 初始化并发和同步映射
var concurrentMap = NewConcurrentHashMap(8, 1000)
var syncMap = sync.Map{}

// readFromConcurrentMap 执行对并发映射的读取操作
func readFromConcurrentMap() {
	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		concurrentMap.Get(key)
	}
}

// writeToConcurrentMap 执行对并发映射的写入操作
func writeToConcurrentMap() {
	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		concurrentMap.Set(key, 1)
	}
}

// deleteFromConcurrentMap 执行对并发映射的删除操作
func deleteFromConcurrentMap() {
	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		concurrentMap.Delete(key)
	}
}

// readFromSyncMap 执行对同步映射的读取操作
func readFromSyncMap() {
	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		syncMap.Load(key)
	}
}

// writeToSyncMap 执行对同步映射的写入操作
func writeToSyncMap() {
	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		syncMap.Store(key, 1)
	}
}

// deleteFromSyncMap 执行对同步映射的删除操作
func deleteFromSyncMap() {
	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		syncMap.Delete(key)
	}
}

// BenchmarkConcurrentMap 基准测试并发映射的读取、写入和删除操作
func BenchmarkConcurrentMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const numGoroutines = 300
		wg := sync.WaitGroup{}
		wg.Add(3 * numGoroutines)
		for i := 0; i < numGoroutines; i++ { // 300 个 goroutine 用于读取
			go func() {
				defer wg.Done()
				readFromConcurrentMap()
			}()
		}
		for i := 0; i < numGoroutines; i++ { // 300 个 goroutine 用于写入
			go func() {
				defer wg.Done()
				writeToConcurrentMap()
			}()
		}
		for i := 0; i < numGoroutines; i++ { // 300 个 goroutine 用于删除
			go func() {
				defer wg.Done()
				deleteFromConcurrentMap()
			}()
		}
		wg.Wait()
	}
}

// BenchmarkSyncMap 基准测试同步映射的读取、写入和删除操作
func BenchmarkSyncMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const numGoroutines = 300
		wg := sync.WaitGroup{}
		wg.Add(3 * numGoroutines)
		for i := 0; i < numGoroutines; i++ { // 300 个 goroutine 用于读取
			go func() {
				defer wg.Done()
				readFromSyncMap()
			}()
		}
		for i := 0; i < numGoroutines; i++ { // 300 个 goroutine 用于写入
			go func() {
				defer wg.Done()
				writeToSyncMap()
			}()
		}
		for i := 0; i < numGoroutines; i++ { // 300 个 goroutine 用于删除
			go func() {
				defer wg.Done()
				deleteFromSyncMap()
			}()
		}
		wg.Wait()
	}
}

// TestConcurrentHashMapDelete 测试 ConcurrentHashMap 的删除方法
func TestConcurrentHashMapDelete(t *testing.T) {
	// 设置初始数据
	for i := 0; i < 10; i++ {
		concurrentMap.Set(strconv.Itoa(i), i)
	}

	// 删除数据
	for i := 0; i < 10; i++ {
		concurrentMap.Delete(strconv.Itoa(i))
	}

	// 验证数据是否被删除
	for i := 0; i < 10; i++ {
		if _, exists := concurrentMap.Get(strconv.Itoa(i)); exists {
			t.Errorf("key %d 应该已被删除", i)
		}
	}
}

// TestConcurrentHashMapIterator 测试并发映射的迭代器
func TestConcurrentHashMapIterator(t *testing.T) {
	for i := 0; i < 10; i++ {
		concurrentMap.Set(strconv.Itoa(i), i)
	}
	iterator := concurrentMap.CreateIterator()
	entry := iterator.Next()
	for entry != nil {
		fmt.Println(entry.Key, entry.Value)
		entry = iterator.Next()
	}
}

// go test -v ./util/test -run=^TestConcurrentHashMapIterator$ -count=1
// go test ./utils -bench=Map -run=^$ -count=1 -benchmem -benchtime=3s
/*
goos: windows
goarch: amd64
pkg: dqq-search-engine/utils
cpu: 13th Gen Intel(R) Core(TM) i7-1360P
BenchmarkConcurrentMap-16              4        1081605200 ns/op        550542020 B/op   9110571 allocs/op
BenchmarkSyncMap-16                    1        4764948600 ns/op        804182752 B/op  21111612 allocs/op
*/
