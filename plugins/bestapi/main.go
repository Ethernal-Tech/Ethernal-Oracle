package main

import (
	"oracle-test/plugins"
)

type BestApi struct{}

func (b *BestApi) Initialize() {}

func (b *BestApi) GetMethods() []plugins.Method {
	return plugins.DefaulGetMethods(b)
}

func (b *BestApi) CallMethod(methodName string, params ...interface{}) (interface{}, error) {
	return plugins.DefaultCallMethod(b, methodName, params...)
}

func (b *BestApi) SayHello() (string, error) {
	return "Hello from BestAPI", nil
}

func main() {}

var ExportPlugin = BestApi{}
