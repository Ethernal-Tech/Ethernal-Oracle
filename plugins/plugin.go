package plugins

import (
	"fmt"
	"reflect"
)

type IPlugin interface {
	Initialize()
	GetMethods() []Method
	CallMethod(methodName string, params ...interface{}) (interface{}, error)
}

type Method struct {
	MethodName   string
	InputParams  []Param
	OutputParams []Param
}

type Param struct {
	ParamType reflect.Type
}

func DefaulGetMethods(structPointer interface{}) []Method {
	structType := reflect.TypeOf(structPointer)

	numMethods := structType.NumMethod()
	var methods []Method

	for i := 0; i < numMethods; i++ {
		method := structType.Method(i)

		if IsIPluginMethod(method.Name) {
			continue
		}

		numInMethods := method.Type.NumIn()
		var inParams []Param

		for j := 0; j < numInMethods; j++ {
			inParams = append(inParams, Param{
				ParamType: method.Type.In(j),
			})
		}

		numOutMethods := method.Type.NumOut()
		var outParams []Param

		for j := 0; j < numOutMethods; j++ {
			outParams = append(outParams, Param{
				ParamType: method.Type.Out(j),
			})
		}

		methods = append(methods, Method{
			MethodName:   method.Name,
			InputParams:  inParams,
			OutputParams: outParams,
		})
	}

	return methods
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
