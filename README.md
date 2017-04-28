# consulsync

Syncs a collection of YAML / JSON files to consul KV. 

## Usage

```
consulsync -datacenter ams2 -consul-address demo.consul.io -kv-prefix test jsonfile1.json yamlfile1.yaml yamlfile2.yaml
```

You may specify `-dryrun` to only show the changes that would be made and `-quiet` to silence any output.

## Status
Only tested on a prototype project but PRs welcome.

ACL tokens nor CA certificates are supported at this time but since we are using `github.com/hashicorp/consul/api` for the consul client, support is there; we just need to expose the options to configure it.
