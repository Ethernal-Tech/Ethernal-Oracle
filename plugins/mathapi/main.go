package main

import (
	"encoding/binary"
	"fmt"
	"oracle-test/plugins"
	"reflect"
)

type MathApi struct{}

func (m *MathApi) Initialize() {}

func (m *MathApi) GetMethods() []plugins.Method {
	structType := reflect.TypeOf(m)

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

func (m *MathApi) CallMethod(methodName string, paramBytes ...[]byte) ([]byte, error) {
	methodValue := reflect.ValueOf(m).MethodByName(methodName)

	if methodValue.IsValid() {
		var methodParams []reflect.Value
		for _, param := range paramBytes {
			methodParams = append(methodParams, reflect.ValueOf(param))
		}

		result := methodValue.Call(methodParams)

		if len(result) > 0 {
			value, _ := result[0].Interface().([]uint8)
			err, _ := result[1].Interface().(error)
			return value, err
		}
		return nil, fmt.Errorf("Method %s did not return expected values", methodName)
	}

	return nil, fmt.Errorf("Method %s not found", methodName)
}

func (m *MathApi) Add_Numbers(first []byte, second []byte) ([]byte, error) {

	firstNum := binary.BigEndian.Uint64(first)
	secondNum := binary.BigEndian.Uint64(second)
	var added = firstNum + secondNum

	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, added)

	return bytes, nil
}

func main() {}

var ExportPlugin = MathApi{}
