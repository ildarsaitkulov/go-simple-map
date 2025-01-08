package main

import (
	"fmt"
	"hash/fnv"
)

type Item[K comparable, V any] struct {
	key   K
	value V
	next  *Item[K, V]
}

type SimpleMap[K comparable, V any] struct {
	buckets []*Item[K, V]
	size    int
}

func NewSimpleMap[K comparable, V any](bucketsCount int) *SimpleMap[K, V] {
	return &SimpleMap[K, V]{
		buckets: make([]*Item[K, V], bucketsCount),
		size:    0,
	}
}

func (m *SimpleMap[K, V]) Put(key K, value V) {
	hash := hashKey(key)
	bkey := int(hash % uint32(len(m.buckets)))

	item := m.buckets[bkey]
	for item != nil {
		if item.key == key { //перезаписываем значение
			item.value = value
			return
		}
		if item.next == nil {
			break
		}
		item = item.next
	}

	newItem := &Item[K, V]{key: key, value: value, next: item}
	m.buckets[bkey] = newItem
	m.size++

}

func (m *SimpleMap[K, V]) Get(key K) (V, bool) {
	hash := hashKey(key)
	bkey := int(hash % uint32(len(m.buckets)))

	item := m.buckets[bkey]

	for item != nil {
		if item.key == key {
			return item.value, true
		}
		item = item.next
	}

	var zeroVal V

	return zeroVal, false
}

func (m *SimpleMap[K, V]) Delete(key K) {
	bkey := int(hashKey(key) % uint32(len(m.buckets)))

	item := m.buckets[bkey]
	var prev *Item[K, V]

	for item != nil {
		if item.key == key {
			// Удаляем элемент.
			if prev == nil {
				m.buckets[bkey] = item.next
			} else {
				prev.next = item.next
			}
			m.size--
			return
		}
		prev = item
		item = item.next
	}
}
func (m *SimpleMap[K, V]) Size() int {
	return m.size
}

// hashKey вычисляет хеш для ключа.
func hashKey[K comparable](key K) uint32 {
	hasher := fnv.New32a()
	_, _ = fmt.Fprintf(hasher, "%v", key) // Хешируем строковое представление ключа.
	return hasher.Sum32()
}

func main() {
	smap := NewSimpleMap[string, int](2)
	smap.Put("1", 1)
	smap.Put("3", 3)
	smap.Put("4", 4)

	fmt.Println(smap.Get("3"))
	fmt.Println(smap.Get("1"))
	fmt.Println(smap.Get("4"))
	smap.Delete("3")
	fmt.Println(smap.Get("3"))

	fmt.Println(smap.Size())
}
