package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"regexp"
	"strings"
)

type stringMarshaller func(data interface{}) string

func applyIfKeyMatch(key string, pattern string, fn stringMarshaller, data interface{}) (string, bool) {
	if matches, _ := regexp.MatchString(pattern, key); matches {
		return fn(data), true
	}
	return "", false
}

func envMarshal(data interface{}) string {
	var keys []string
	out := &([]string{})

	gt := &Traverser{Node: func(keys []string, data interface{}) {
		*out = append(*out, fmt.Sprintf("%v=%v", strings.ToUpper(strings.Join(keys, "_")), data))
	}}
	gt.Traverse(data, keys)
	return strings.Join(*out, "\n")
}

func yamlMarshal(data interface{}) string {
	yaml, _ := yaml.Marshal(data)
	return string(yaml)
}

func jsonMarshal(data interface{}) string {
	tempMap := make(map[string]interface{})
	for k, v := range data.(map[interface{}]interface{}) {
		tempMap[fmt.Sprintf("%v", k)] = v
	}
	json, _ := json.Marshal(tempMap)
	return string(json)

}
