package utils

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type OrderedMap[K comparable, V any] struct {
	keys   []K
	values map[K]V
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		keys:   []K{},
		values: make(map[K]V),
	}
}

func (om *OrderedMap[K, V]) Set(key K, value V) {
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}

func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	value, exists := om.values[key]
	return value, exists
}

func (om *OrderedMap[K, V]) Delete(key K) {
	delete(om.values, key)
	for i, k := range om.keys {
		if k == key {
			om.keys = append(om.keys[:i], om.keys[i+1:]...)
			break
		}
	}
}

func (om *OrderedMap[K, V]) ForEach(fn func(K, V)) {
	for _, key := range om.keys {
		fn(key, om.values[key])
	}
}

func (om *OrderedMap[K, V]) MarshalYAML() (any, error) {
	node := &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: []*yaml.Node{},
	}

	for _, key := range om.keys {
		kNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: fmt.Sprintf("%v", key),
		}

		vNode := &yaml.Node{}
		val := om.values[key]
		if err := vNode.Encode(val); err != nil {
			return nil, err
		}

		node.Content = append(node.Content, kNode, vNode)
	}

	return node, nil
}
