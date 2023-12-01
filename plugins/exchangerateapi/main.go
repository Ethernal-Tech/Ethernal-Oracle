package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"oracle-test/plugins"
	"os"
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

type ExchangeRateApi struct {
	api_key string
	address string
}

func (e *ExchangeRateApi) Initialize() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file %v", err)
	}

	e.address = os.Getenv("EXCHANGE_RATE_ADDRESS")
	e.api_key = os.Getenv("EXCHANGE_RATE_API_KEY")

	return nil
}

func (e *ExchangeRateApi) GetMethods() ([]plugins.Method, error) {
	return plugins.DefaultGetMethods(e)
}

func (e *ExchangeRateApi) CallMethod(methodName string, params ...interface{}) (interface{}, error) {
	return plugins.DefaultCallMethod(e, methodName, params...)
}

func (e *ExchangeRateApi) Exhange_Rate(base string, target string) (float64, error) {
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

var ExportPlugin = ExchangeRateApi{}
