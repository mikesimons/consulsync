package main

import (
	"github.com/hashicorp/consul/api"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func readFile(file string) map[interface{}]interface{} {
	contents := make(map[interface{}]interface{})

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Unable to read %s: %s", file, err)
	}

	err = yaml.Unmarshal(data, &contents)
	if err != nil {
		log.Fatalf("Unable to parse %s: %s", file, err)
	}

	return contents
}

func process(data interface{}, output *map[string]interface{}) {
	var keys []string
	if output == nil {
		tmp := make(map[string]interface{})
		output = &tmp
	}

	stringMarshallers := map[string]stringMarshaller{
		"!env$":  envMarshal,
		"!yaml$": yamlMarshal,
		"!json$": jsonMarshal,
	}

	gt := &Traverser{}

	gt.Node = func(keys []string, data interface{}) {
		if data == nil {
			data = (interface{})("")
		}
		(*output)[strings.Join(keys, "/")] = data
	}

	gt.Map = func(keys []string, key string, data interface{}) {
		for pattern, fn := range stringMarshallers {
			if value, applied := applyIfKeyMatch(key, pattern, fn, data); applied {
				(*output)[strings.Join(append(keys, key), "/")] = value
				return
			}
		}

		gt.Traverse(data, append(keys, key))
	}

	gt.Traverse(data, keys)
}

func main() {
	app := cli.NewApp()
	app.Name = "consulsync"
	app.Usage = "Sync YAML / JSON files to consul KV"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "consul-address",
			Usage: "JSON / YAML file",
		},
		cli.StringFlag{
			Name:  "datacenter",
			Usage: "Datacenter",
		},
		cli.StringFlag{
			Name:  "kv-prefix",
			Usage: "KV prefix to sync",
		},
		cli.BoolFlag{
			Name:  "dryrun",
			Usage: "Do not perform changes",
		},
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "Do not print output",
		},
	}

	app.Action = func(c *cli.Context) error {
		out := make(map[string]interface{})
		for _, file := range c.Args() {
			data := readFile(file)
			process(interface{}(data), &out)
		}

		config := &api.Config{
			Address:    c.GlobalString("consul-address"),
			Datacenter: c.GlobalString("datacenter"),
		}

		consul, err := api.NewClient(config)
		if err != nil {
			panic(err)
		}

		consulData := fetchConsulData(consul, c.GlobalString("kv-prefix"))
		diff := diffConsulData(c.GlobalString("kv-prefix"), consulData, out)

		if !c.Bool("quiet") {
			diff.Print()
		}

		if !c.Bool("dryrun") {
			applyDiff(consul, diff)
		}

		return nil
	}

	app.Run(os.Args)
}
