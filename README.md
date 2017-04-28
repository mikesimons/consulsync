# consulsync

Syncs a collection of YAML / JSON files to consul KV.

The source files are authoritative and sync is one way (for now).
Key & content diffing is performed so only differences will cause KV writes to happen.

## Usage

```
consulsync -datacenter ams2 -consul-address demo.consul.io -kv-prefix test jsonfile1.json yamlfile1.yaml yamlfile2.yaml
```

You may specify `-dryrun` to only show the changes that would be made and `-quiet` to silence any output.

When multiple files are provided, each one is processed in turn in to the same map. This means that the contents of the later files will override the earlier ones.

### String filters

While the main purpose of consulsync is to mirror structures from version controlled config files in to KV prefixes, there are some specific serialisation features to make some things nicer to do. These filters cause the subtree of a key to be marshaled in to a string format and stored in the key that contained the trigger tag.

You can specify a trigger tag by appending either `!env`, `!yaml` or `!json` to a key. Please note that everything in the key before the tag is arbitrary and need not contain `.json` for the json filter for example.

### JSON
The following will store the string `{ "somevalue": "test", "someothervalue": "test" }` in the `service1/configfile.json` key.

The subtree is simply serialized to JSON. Please note that the subtree must be a map / hash type (and not a scalar / list type).

```
service1:
  configfile.json!json:
    somevalue: test
    someothervalue: test
```
### YAML
The following will store the string `somevalue: test\nsomeothervalue: test\n` in the `service1/configfile.yaml` key.

The subtree is simply serialized to YAML.

```
service1:
  configfile.yaml!yaml:
    somevalue: test
    someothervalue: test
```
### ENV
The following will store the string `SOME_VALUE=test\nSOMEOTHERVALUE=test` in the `service1/configfile.env` key.

The subtree is traversed and flattened. All keys are uppercased and nested keys are joined with `_`. Please note that no escaping is performed on keys nor values at this time. 

```
service1:
  configfile.env!env:
    some:
      value: test
    someothervalue: test
```

## Status
Only tested on a prototype project but PRs welcome.

ACL tokens nor CA certificates are supported at this time but since we are using `github.com/hashicorp/consul/api` for the consul client, support is there; we just need to expose the options to configure it.
