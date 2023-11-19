package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

type MathApi struct{}

func (m *MathApi) Initialize() {}

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
