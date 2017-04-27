package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type Diff struct {
	Add    map[string]interface{}
	Change map[string]interface{}
	Remove []string
}

func (d *Diff) HasChanges() bool {
	return len(d.Add) > 0 || len(d.Change) > 0 || len(d.Remove) > 0
}

type dataMap map[string]interface{}

func diffConsulData(kvPrefix string, consulData api.KVPairs, data map[string]interface{}) *Diff {
	result := &Diff{
		Add:    make(map[string]interface{}),
		Change: make(map[string]interface{}),
	}

	for _, pair := range consulData {
		k := pair.Key[len(kvPrefix)+1 : len(pair.Key)]
		if _, ok := data[k]; ok {
			if data[k] != string(pair.Value) {
				result.Change[pair.Key] = data[k]
			}
		} else {
			result.Remove = append(result.Remove, pair.Key)
		}
		delete(data, k)
	}

	for k, v := range data {
		k = fmt.Sprintf("%v/%v", kvPrefix, k)
		result.Add[k] = v
	}

	return result
}

func (d *Diff) Print() {
	if !d.HasChanges() {
		fmt.Printf("No changes\n")
		return
	}

	if len(d.Add) > 0 {
		fmt.Printf("To add:\n%v\n", d.Add)
	}

	if len(d.Change) > 0 {
		fmt.Printf("To change:\n%v\n", d.Change)
	}

	if len(d.Remove) > 0 {
		fmt.Printf("To remove:\n%v\n", d.Remove)
	}
}
