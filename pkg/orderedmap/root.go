package orderedmap

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

func (om *OrderedMap[K, V]) Values() map[K]V {
	return om.values
}

func (om *OrderedMap[K, V]) Keys() []K {
	return om.keys
}

func (om *OrderedMap[K, V]) Has(key K) bool {
	_, exists := om.values[key]
	return exists
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

func (om *OrderedMap[K, V]) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("expected mapping node but got %v", node.Kind)
	}

	om.keys = make([]K, 0, len(node.Content)/2)
	om.values = make(map[K]V, len(node.Content)/2)

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		var keyStr string
		if err := keyNode.Decode(&keyStr); err != nil {
			return fmt.Errorf("failed to decode key: %w", err)
		}

		var key K
		if err := yaml.Unmarshal([]byte(keyStr), &key); err != nil {
			return fmt.Errorf("failed to parse key string '%s': %w", keyStr, err)
		}

		var value V
		if err := valueNode.Decode(&value); err != nil {
			return fmt.Errorf("failed to decode value for key %v: %w", key, err)
		}

		om.keys = append(om.keys, key)
		om.values[key] = value
	}

	return nil
}