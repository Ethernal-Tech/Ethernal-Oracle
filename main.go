package main

import (
	"fmt"
	"oracle-test/plugins"
	"plugin"
)

func main() {
	best, err := loadPlugin("build/plugins/bestapi.so")
	if err != nil {
		fmt.Println(err)
		return
	}

	best.Initialize()
	res := best.CallMethod("wgatever")
	fmt.Println(res)
}

func loadPlugin(path string) (plugins.IPlugin, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	sym, err := p.Lookup("ExportPlugin")
	if err != nil {
		return nil, err
	}

	return sym.(plugins.IPlugin), nil
}
