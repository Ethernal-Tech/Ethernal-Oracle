package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
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

func (e *ExhangeRateApi) CallMethod(methodName string, paramBytes ...[]byte) ([]byte, error) {
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

func (e *ExhangeRateApi) Exhange_Rate(first []byte, second []byte) ([]byte, error) {
	baseCurrency := string(first)
	targetCurrency := string(second)

	apiUrl := e.address
	apiUrl = strings.Replace(apiUrl, "[API_KEY]", e.api_key, -1)
	apiUrl = strings.Replace(apiUrl, "[BASE_CURRENCY]", baseCurrency, -1)

	response, err := http.Get(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %s", err)
	}

	var rates ExchangeRateResponse
	err = json.Unmarshal(body, &rates)
	if err != nil {
		return nil, fmt.Errorf("Error decoding JSON: %s", err)
	}

	floatString := strconv.FormatFloat(rates.ConversionRates[targetCurrency], 'f', -1, 64)

	return []byte(floatString), nil
}

func main() {}

var ExportPlugin = ExhangeRateApi{}
