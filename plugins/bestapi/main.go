package main

import (
	"oracle-test/plugins"
	"reflect"
)

type BestApi struct{}

func (b *BestApi) Initialize() {}

func (b *BestApi) GetMethods() []plugins.Method {
	structType := reflect.TypeOf(b)

	numMethods := structType.NumMethod()
	methodCount := 0
	var methods = make([]plugins.Method, numMethods-3)
	for i := 0; i < numMethods; i++ {
		method := structType.Method(i)

		if method.Name == "Initialize" ||
			method.Name == "GetMethods" ||
			method.Name == "CallMethod" {
			continue
		}

		var newMethod = plugins.Method{}
		newMethod.MethodName = method.Name

		numParams := method.Type.NumIn()
		var inputParams = make([]plugins.Param, numParams)
		for j := 0; j < numParams; j++ {
			inputParams[j].ParamType = method.Type.In(j).String()
		}

		numOut := method.Type.NumOut()
		var outputParams = make([]plugins.Param, numOut)
		for j := 0; j < numOut; j++ {
			outputParams[j].ParamType = method.Type.Out(j).String()
		}

		newMethod.InputParams = inputParams
		newMethod.OutputParams = outputParams
		methods[methodCount] = newMethod
		methodCount++
	}

	return methods
}

func (b *BestApi) CallMethod(methodName string, paramBytes ...[]byte) ([]byte, error) {
	return []byte("Hello from BestAPI"), nil
}

func main() {}

var ExportPlugin = BestApi{}
