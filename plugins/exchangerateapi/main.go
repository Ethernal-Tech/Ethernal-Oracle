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
	"strings"

	"github.com/joho/godotenv"
)

type ExchangeRateResponse struct {
	Result             string             `json:"result"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeLastUpdateUTC  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	TimeNextUpdateUTC  string             `json:"time_next_update_utc"`
	BaseCode           string             `json:"base_code"`
	ConversionRates    map[string]float64 `json:"conversion_rates"`
}

type ExhangeRateApi struct {
	api_key string
	address string
}

func (e *ExhangeRateApi) Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e.address = os.Getenv("EXCHANGE_RATE_ADDRESS")
	e.api_key = os.Getenv("EXCHANGE_RATE_API_KEY")
}

func (e *ExhangeRateApi) GetMethods() []plugins.Method {
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

func (e *ExhangeRateApi) CallMethod(methodName string, params ...interface{}) (interface{}, error) {
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

func (e *ExhangeRateApi) Exhange_Rate(base string, target string) (float64, error) {
	apiUrl := e.address
	apiUrl = strings.Replace(apiUrl, "[API_KEY]", e.api_key, -1)
	apiUrl = strings.Replace(apiUrl, "[BASE_CURRENCY]", base, -1)

	response, err := http.Get(apiUrl)
	if err != nil {
		return 0, fmt.Errorf("Error making request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("Error making request: %s", err)
	}

	var rates ExchangeRateResponse
	err = json.Unmarshal(body, &rates)
	if err != nil {
		return 0, fmt.Errorf("Error decoding JSON: %s", err)
	}

	return rates.ConversionRates[target], nil
}

func main() {}

var ExportPlugin = ExhangeRateApi{}
