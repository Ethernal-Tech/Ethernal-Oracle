package main

import (
	"oracle-test/plugins"
)

type MathApi struct{}

func (m *MathApi) Initialize() error {
	return nil
}

func (m *MathApi) GetMethods() ([]plugins.Method, error) {
	return plugins.DefaultGetMethods(m)
}

func (m *MathApi) CallMethod(methodName string, params ...interface{}) (interface{}, error) {
	return plugins.DefaultCallMethod(m, methodName, params...)
}

func (m *MathApi) Add_Numbers(first uint64, second uint64) (uint64, error) {
	var added = first + second

	return added, nil
}

func main() {}

var ExportPlugin = MathApi{}
