package plugins

import (
	"fmt"
	"reflect"
)

type IPlugin interface {
	Initialize() error
	GetMethods() ([]Method, error)
	CallMethod(methodName string, params ...interface{}) (interface{}, error)
}

type Method struct {
	MethodName   string
	InputParams  []reflect.Type
	OutputParams []reflect.Type
}

func DefaultGetMethods(structPointer interface{}) ([]Method, error) {
	structType := reflect.TypeOf(structPointer)
	if structType == nil {
		return nil, fmt.Errorf("Failed to get type of a struct")
	}

	numMethods := structType.NumMethod()
	var methods []Method

	for i := 0; i < numMethods; i++ {
		method := structType.Method(i)

		if IsIPluginMethod(method.Name) {
			// Skip IPlugin methods
			continue
		}

		numInMethods := method.Type.NumIn()
		var inParams []reflect.Type

		for j := 0; j < numInMethods; j++ {
			inParams = append(inParams, method.Type.In(j))
		}

		numOutMethods := method.Type.NumOut()
		var outParams []reflect.Type

		for j := 0; j < numOutMethods; j++ {
			outParams = append(outParams, method.Type.Out(j))
		}

		methods = append(methods, Method{
			MethodName:   method.Name,
			InputParams:  inParams,
			OutputParams: outParams,
		})
	}

	return methods, nil
}

func DefaultCallMethod(structPointer interface{}, methodName string, params ...interface{}) (interface{}, error) {
	methodValue := reflect.ValueOf(structPointer).MethodByName(methodName)

	if methodValue.IsValid() {
		var methodParams []reflect.Value
		for _, param := range params {
			methodParams = append(methodParams, reflect.ValueOf(param))
		}

		result := methodValue.Call(methodParams)

		if len(result) > 0 {
			value, _ := result[0].Interface().(interface{})
			err, _ := result[1].Interface().(error)
			return value, err
		}
		return nil, fmt.Errorf("Method %s did not return expected values", methodName)
	}

	return nil, fmt.Errorf("Method %s not found", methodName)
}

func IsIPluginMethod(name string) bool {
	ipluginMethods := []string{"Initialize", "GetMethods", "CallMethod"}

	for _, method := range ipluginMethods {
		if name == method {
			return true
		}
	}

	return false
}
