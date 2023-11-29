package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"oracle-test/plugins"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type ConvertRateResponse struct {
	Success bool    `json:"success"`
	Query   Query   `json:"query"`
	Info    Info    `json:"info"`
	Date    string  `json:"date"`
	Result  float64 `json:"result"`
}

type Query struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int32  `json:"amount"`
}

type Info struct {
	Timestamp int64   `json:"timestamp"`
	Rate      float64 `json:"rate"`
}

type CurrencyConversionApi struct {
	api_key string
	address string
}

func (e *CurrencyConversionApi) Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e.address = os.Getenv("CURRENCY_CONVERSION_ADDRESS")
	e.api_key = os.Getenv("CURRENCY_CONVERSION_API_KEY")
}

func (e *CurrencyConversionApi) GetMethods() []plugins.Method {
	structType := reflect.TypeOf(e)

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

func (e *CurrencyConversionApi) CallMethod(methodName string, params ...interface{}) (interface{}, error) {
	methodValue := reflect.ValueOf(e).MethodByName(methodName)

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

func (e *CurrencyConversionApi) Exhange_Rate(base string, target string) (float64, error) {
	apiUrl := e.address + fmt.Sprintf("convert?from=%s&to=%s&amount=1.0", base, target)

	request, err := http.NewRequest("GET", apiUrl, nil)
	request.Header.Add("X-RapidAPI-Key", e.api_key)
	request.Header.Add("X-RapidAPI-Host", "currency-conversion-and-exchange-rates.p.rapidapi.com")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, fmt.Errorf("Error making request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("Error making request: %s", err)
	}

	var rate ConvertRateResponse
	err = json.Unmarshal(body, &rate)
	if err != nil {
		return 0, fmt.Errorf("Error decoding JSON: %s", err)
	}

	return rate.Result, nil
}

func main() {}

var ExportPlugin = CurrencyConversionApi{}
