package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"oracle-test/plugins"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type LiveScore struct {
	api_key string
	address string
}

func (s *LiveScore) Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s.address = os.Getenv("LIVESCORE_ADDRESS")
	s.api_key = os.Getenv("LIVESCORE_API_KEY")
}

func (s *LiveScore) GetMethods() []plugins.Method {
	structType := reflect.TypeOf(s)

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

func (s *LiveScore) CallMethod(methodName string, paramBytes ...[]byte) ([]byte, error) {
	methodValue := reflect.ValueOf(s).MethodByName(methodName)

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

func (s *LiveScore) Get_All_Leagues() ([]byte, error) {
	apiUrl := s.address + "list?Category=soccer"

	request, err := http.NewRequest("GET", apiUrl, nil)
	request.Header.Add("X-RapidAPI-Key", s.api_key)
	request.Header.Add("X-RapidAPI-Host", "livescore6.p.rapidapi.com")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %s", err)
	}

	return body, nil
}

func main() {}

var ExportPlugin = LiveScore{}
