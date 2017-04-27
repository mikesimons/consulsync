package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

func fetchConsulData(client *api.Client, kvPrefix string) api.KVPairs {
	kv := client.KV()
	keys, _, err := kv.List(kvPrefix, nil)
	if err != nil {
		panic(err)
	}
	return keys
}

func applyDiff(client *api.Client, diff *Diff) {
	kv := client.KV()
	for k, v := range diff.Add {
		kv.Put(&api.KVPair{Key: k, Value: []byte(fmt.Sprintf("%v", v))}, nil)
	}

	for k, v := range diff.Change {
		kv.Put(&api.KVPair{Key: k, Value: []byte(fmt.Sprintf("%v", v))}, nil)
	}

	for _, v := range diff.Remove {
		kv.Delete(v, nil)
	}
}
