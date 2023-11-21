package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type CurrencyExchangeApi struct {
	api_key string
	address string
}

func (e *CurrencyExchangeApi) Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e.address = os.Getenv("CURRENCY_EXCHANGE_ADDRESS")
	e.api_key = os.Getenv("CURRENCY_EXCHANGE_API_KEY")
}

func (e *CurrencyExchangeApi) CallMethod(methodName string, paramBytes ...[]byte) ([]byte, error) {
	methodValue := reflect.ValueOf(e).MethodByName(methodName)

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

func (e *CurrencyExchangeApi) Exhange_Rate(first []byte, second []byte) ([]byte, error) {
	baseCurrency := string(first)
	targetCurrency := string(second)

	apiUrl := e.address + fmt.Sprintf("exchange?from=%s&to=%s&q=1.0", baseCurrency, targetCurrency)

	request, err := http.NewRequest("GET", apiUrl, nil)
	request.Header.Add("X-RapidAPI-Key", e.api_key)
	request.Header.Add("X-RapidAPI-Host", "currency-exchange.p.rapidapi.com")

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

var ExportPlugin = CurrencyExchangeApi{}
