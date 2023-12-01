package main

import (
	"oracle-test/plugins"
)

type BestApi struct{}

func (b *BestApi) Initialize() error {
	return nil
}

func (b *BestApi) GetMethods() ([]plugins.Method, error) {
	return plugins.DefaultGetMethods(b)
}

func (b *BestApi) CallMethod(methodName string, params ...interface{}) (interface{}, error) {
	return plugins.DefaultCallMethod(b, methodName, params...)
}

func (b *BestApi) SayHello() (string, error) {
	return "Hello from BestAPI", nil
}

func main() {}

var ExportPlugin = BestApi{}
