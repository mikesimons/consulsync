package main

import (
	"fmt"
	"reflect"
)

type Traverser struct {
	Map  func(keys []string, key string, data interface{})
	Node func(keys []string, data interface{})
}

func (gt *Traverser) Traverse(data interface{}, keys []string) {
	valueOfData := reflect.ValueOf(data)
	switch valueOfData.Kind() {
	case reflect.Map:
		d := data.(map[interface{}]interface{})
		for k, v := range d {
			ks := fmt.Sprintf("%v", k)
			if gt.Map != nil {
				gt.Map(keys, ks, v)
			} else {
				gt.Traverse(v, append(keys, ks))
			}
		}
	case reflect.Slice:
		d := data.([]interface{})
		for k, v := range d {
			lastKey := keys[len(keys)-1]
			gt.Traverse(v, append(keys, fmt.Sprintf("%v%d", lastKey, k)))
		}
	default:
		gt.Node(keys, data)
	}
}
