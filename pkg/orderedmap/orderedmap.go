package orderedmap

import (
	"fmt"
	"iter"

	"gopkg.in/yaml.v3"
)

// OrderedMap keeps the order of keys and values.
// This is developed for dependencies value.
type OrderedMap struct {
	keys   []string
	values map[string][]string
}

// New creates a new OrderedMap.
func New() *OrderedMap {
	return &OrderedMap{
		keys:   make([]string, 0),
		values: make(map[string][]string),
	}
}

// Iter returns an iterator over the OrderedMap.
func (om *OrderedMap) Iter() iter.Seq2[string, []string] {
	return func(yield func(string, []string) bool) {
		for _, key := range om.keys {
			if !yield(key, om.values[key]) {
				break
			}
		}
	}
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (om *OrderedMap) UnmarshalYAML(value *yaml.Node) error {
	*om = *New()

	// If the document node, process the content.
	if value.Kind == yaml.DocumentNode && len(value.Content) > 0 {
		value = value.Content[0]
	}

	// If the value is not a mapping node, return an error.
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("value is not mapping node")
	}

	// Process the key and value pairs.
	for i := 0; i < len(value.Content); i += 2 {
		keyNode := value.Content[i]
		valueNode := value.Content[i+1]

		if keyNode.Kind != yaml.ScalarNode {
			return fmt.Errorf("key is not scalar node")
		}
		if valueNode.Kind != yaml.SequenceNode {
			return fmt.Errorf("value is not sequence node")
		}

		key := keyNode.Value
		var parsedValue []string
		if err := valueNode.Decode(&parsedValue); err != nil {
			return err
		}

		om.set(key, parsedValue)
	}

	return nil
}

// set adds a key and value to the OrderedMap.
func (om *OrderedMap) set(key string, value []string) {
	// If the key does not exist, add it.
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}
