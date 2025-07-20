# CMap - 高性能并发安全哈希表

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Lofanmi/cmap)](https://goreportcard.com/report/github.com/Lofanmi/cmap)
[![GoDoc](https://godoc.org/github.com/Lofanmi/cmap?status.svg)](https://godoc.org/github.com/Lofanmi/cmap)
[![Tests](https://img.shields.io/badge/Tests-Passing-brightgreen.svg)](https://github.com/Lofanmi/cmap/actions)
[![Coverage](https://img.shields.io/badge/Coverage-44.3%25-yellow.svg)](https://github.com/Lofanmi/cmap)

一个高性能的并发安全哈希表，完全兼容 [gods](https://github.com/emirpasic/gods) 库的接口设计，支持泛型和多种配置选项。

## ✨ 特性

- 🔒 **完全并发安全** - 通过分片技术实现高并发访问
- 🚀 **高性能** - 在并发场景下性能显著优于标准 map + 互斥锁
- 🔧 **完全兼容 Gods** - 直接使用 gods 库的 Map 接口，无缝集成
- 🎯 **泛型支持** - 完全支持 Go 1.21+ 泛型
- 📊 **多种Map类型** - 支持 HashMap, TreeMap, LinkedHashMap
- 💾 **序列化支持** - 支持 JSON, JSONiter, Sonic, Gob 序列化格式
- 📁 **文件操作** - 支持保存到文件和从文件加载
- ⚙️ **可配置** - 支持分片数量、Map类型、序列化格式等配置
- 🧪 **全面测试** - 包含单元测试、并发测试和性能基准测试

## 📈 性能表现（TODO）

- TODO

## 🚀 快速开始

### 安装

```bash
go get -u github.com/Lofanmi/cmap
```

### 基本使用

```go
package main

import (
    "fmt"
	
    "github.com/Lofanmi/cmap"
)

func main() {
    // 创建并发映射
    cm := cmap.New[string, int]()
    
    // 设置值
    cm.Put("apple", 5)
    cm.Put("banana", 3)
    
    // 获取值
    if value, found := cm.Get("apple"); found {
        fmt.Printf("苹果数量: %d\n", value)
    }
    
    fmt.Printf("总大小: %d\n", cm.Size())
}
```

### 配置选项

```go
// 创建带配置的并发映射
cm := cmap.New[string, int](
    cmap.WithShardCount(64),          // 64个分片
    cmap.WithSerializer(cmap.JsonSerializer()), // JSON序列化
)
```

## 📚 使用场景

### 1. 高并发 Web 服务

```go
// 用户会话存储
sessionStore := cmap.New[string, *UserSession]()
go func() {
    for {
        session := getUserSession()
        sessionStore.Put(session.ID, session)
    }
}()

// 多个 goroutine 同时读取
go func() {
    for {
        if session, found := sessionStore.Get(sessionID); found {
            processSession(session)
        }
    }
}()
```

### 2. 缓存系统

```go
// 分布式缓存
cache := cmap.New[string, interface{}](cmap.WithShardCount(256))
cache.Put("user:123", userData)
cache.Put("product:456", productData)

// 批量操作
keys := []string{"user:123", "product:456", "order:789"}
results := cache.GetMultiple(keys)
```

### 3. 实时数据处理

```go
// 实时指标收集
metrics := cmap.New[string, float64]()
for _, metric := range incomingMetrics {
    metrics.Put(metric.Name, metric.Value)
}

// 定期聚合
go func() {
    for {
        time.Sleep(time.Minute)
        aggregateMetrics(metrics)
    }
}()
```

## 🔧 API 参考

### 构造函数

```go
// 通用构造函数
func New[K comparable, V any](options ...Option) *Map[K, V]

// String 键专用构造函数
func NewStringHashMap[V any](options ...Option) *Map[string, V]
func NewStringTreeMap[V any](options ...Option) *Map[string, V]
func NewStringLinkedHashMap[V any](options ...Option) *Map[string, V]

// Int 键专用构造函数
func NewIntHashMap[V any](options ...Option) *Map[int, V]
func NewIntTreeMap[V any](options ...Option) *Map[int, V]
func NewIntLinkedHashMap[V any](options ...Option) *Map[int, V]

// Int64 键专用构造函数
func NewInt64HashMap[V any](options ...Option) *Map[int64, V]
func NewInt64TreeMap[V any](options ...Option) *Map[int64, V]
func NewInt64LinkedHashMap[V any](options ...Option) *Map[int64, V]

// 通用类型构造函数
func NewHashMap[K comparable, V any](options ...Option) *Map[K, V]
func NewTreeMap[K comparable, V any](options ...Option) *Map[K, V]
func NewLinkedHashMap[K comparable, V any](options ...Option) *Map[K, V]
```

### 配置选项

```go
func WithShardCount(count uint32) Option
func WithSerializer(serializer *SerializerFunc) Option
```

### 核心方法

```go
// 基本操作
Put(key K, value V)
Get(key K) (value V, found bool)
Remove(key K)
Size() int
Empty() bool
Clear()

// 批量操作
PutAll(data map[K]V)
GetMultiple(keys []K) map[K]V
RemoveMultiple(keys []K)

// 序列化
MarshalJSON() ([]byte, error)
UnmarshalJSON(data []byte) error
MarshalWith(serializer *SerializerFunc) ([]byte, error)
UnmarshalWith(data []byte, serializer *SerializerFunc) error

// 文件操作
SaveToFile(filename string) error
LoadFromFile(filename string) error

// 迭代
Keys() []K
Values() []V
```

## ⚡ 性能优化建议

### 1. 分片数量调优

```go
// 高并发场景
highConcurrencyMap := cmap.New[string, int](cmap.WithShardCount(512))

// 低并发场景
lowConcurrencyMap := cmap.New[string, int](cmap.WithShardCount(16))
```

### 2. 选择合适的 Map 类型

```go
// 需要排序的场景
sortedMap := cmap.NewTreeMap[string, int]()

// 需要保持插入顺序的场景
orderedMap := cmap.NewLinkedHashMap[string, int]()

// 一般场景（推荐）
hashMap := cmap.New[string, int]()
```

### 3. 批量操作优化

```go
// 批量添加（推荐）
data := make(map[string]int, 1000)
for i := 0; i < 1000; i++ {
    data[fmt.Sprintf("key%d", i)] = i
}
cm.PutAll(data)

// 批量获取（推荐）
keys := make([]string, 100)
for i := 0; i < 100; i++ {
    keys[i] = fmt.Sprintf("key%d", i)
}
results := cm.GetMultiple(keys)
```

## 🧪 测试

### 运行所有测试

```bash
go test -v
```

### 运行基准测试

```bash
go test -bench=. -benchmem
```

### 运行并发测试

```bash
go test -v -run TestConcurrent
```

## 🤝 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 Apache License Version 2.0 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

---

⭐ 如果这个项目对您有帮助，请给它一个星标！