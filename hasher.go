package cmap

import (
	"cmp"
	"hash/fnv"
	"math"
	"unsafe"

	"github.com/spf13/cast"
)

// ---------------------------------------------------------------------------------------------------------------------

// hasher 哈希计算器接口
type hasher[K cmp.Ordered] interface {
	Hash(key K) uint32
}

// ---------------------------------------------------------------------------------------------------------------------

type hashKeyUint32 interface {
	~int8 | ~int16 | ~int32 | ~uint8 | ~uint16 | ~uint32
}

type hashKeyUint64 interface {
	~int | ~int64 | ~uint | ~uint64 | ~uintptr
}

// ---------------------------------------------------------------------------------------------------------------------

type int8Hasher struct{}
type int16Hasher struct{}
type int32Hasher struct{}
type uint8Hasher struct{}
type uint16Hasher struct{}
type uint32Hasher struct{}

func (h int8Hasher) Hash(key int8) uint32     { return uint32Hash(key) }
func (h int16Hasher) Hash(key int16) uint32   { return uint32Hash(key) }
func (h int32Hasher) Hash(key int32) uint32   { return uint32Hash(key) }
func (h uint8Hasher) Hash(key uint8) uint32   { return uint32Hash(key) }
func (h uint16Hasher) Hash(key uint16) uint32 { return uint32Hash(key) }
func (h uint32Hasher) Hash(key uint32) uint32 { return uint32Hash(key) }

// ---------------------------------------------------------------------------------------------------------------------

type intHasher struct{}
type int64Hasher struct{}
type uintHasher struct{}
type uint64Hasher struct{}
type uintptrHasher struct{}

func (h intHasher) Hash(key int) uint32         { return uint64Hash(key) }
func (h int64Hasher) Hash(key int64) uint32     { return uint64Hash(key) }
func (h uintHasher) Hash(key uint) uint32       { return uint64Hash(key) }
func (h uint64Hasher) Hash(key uint64) uint32   { return uint64Hash(key) }
func (h uintptrHasher) Hash(key uintptr) uint32 { return uint64Hash(key) }

// ---------------------------------------------------------------------------------------------------------------------

type float32Hasher struct{}
type float64Hasher struct{}

func (h float32Hasher) Hash(key float32) uint32 { return float32Hash(key) }
func (h float64Hasher) Hash(key float64) uint32 { return float64Hash(key) }

// ---------------------------------------------------------------------------------------------------------------------

type stringHasher struct{}

func (h stringHasher) Hash(key string) uint32 {
	v := fnv.New32()
	_, _ = v.Write(unsafe.Slice(unsafe.StringData(key), len(key)))
	return v.Sum32()
}

// ---------------------------------------------------------------------------------------------------------------------

type genericHasher[T cmp.Ordered] struct {
	value T
}

func (h genericHasher[T]) Hash(key T) uint32 {
	str := cast.ToString(key)
	v := fnv.New32()
	_, _ = v.Write(unsafe.Slice(unsafe.StringData(str), len(str)))
	return v.Sum32()
}

// ---------------------------------------------------------------------------------------------------------------------

func getHasher[K cmp.Ordered]() hasher[K] {
	var zero K
	switch any(zero).(type) {
	case int8:
		return any(int8Hasher{}).(hasher[K])
	case int16:
		return any(int16Hasher{}).(hasher[K])
	case int32:
		return any(int32Hasher{}).(hasher[K])
	case uint8:
		return any(uint8Hasher{}).(hasher[K])
	case uint16:
		return any(uint16Hasher{}).(hasher[K])
	case uint32:
		return any(uint32Hasher{}).(hasher[K])
	case int:
		return any(intHasher{}).(hasher[K])
	case int64:
		return any(int64Hasher{}).(hasher[K])
	case uint:
		return any(uintHasher{}).(hasher[K])
	case uint64:
		return any(uint64Hasher{}).(hasher[K])
	case uintptr:
		return any(uintptrHasher{}).(hasher[K])
	case float32:
		return any(float32Hasher{}).(hasher[K])
	case float64:
		return any(float64Hasher{}).(hasher[K])
	case string:
		return any(stringHasher{}).(hasher[K])
	default:
		return any(genericHasher[K]{}).(hasher[K])
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func uint32Hash[T hashKeyUint32](key T) uint32 {
	x := uint32(key)
	x = ((x >> 16) ^ x) * 0x45d9f3b
	x = ((x >> 16) ^ x) * 0x45d9f3b
	x = (x >> 16) ^ x
	return x
}

func uint64Hash[T hashKeyUint64](key T) uint32 {
	x := uint64(key)
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	x = x ^ (x >> 31)
	return uint32(x)
}

func float32Hash(key float32) uint32 {
	if math.IsNaN(float64(key)) {
		return 0x7FFFFFFF
	}
	if key == 0 {
		return 0
	}
	v := *(*uint32)(unsafe.Pointer(&key))
	return uint32Hash(v)
}

func float64Hash(key float64) uint32 {
	if math.IsNaN(key) {
		return 0x7FFFFFFF
	}
	if key == 0 {
		return 0
	}
	v := *(*uint64)(unsafe.Pointer(&key))
	return uint64Hash(v)
}

// ---------------------------------------------------------------------------------------------------------------------
