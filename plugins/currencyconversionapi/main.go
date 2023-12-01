package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"oracle-test/plugins"
	"os"

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

func (e *CurrencyConversionApi) Initialize() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file %v", err)
	}

	e.address = os.Getenv("CURRENCY_CONVERSION_ADDRESS")
	e.api_key = os.Getenv("CURRENCY_CONVERSION_API_KEY")

	return nil
}

func (e *CurrencyConversionApi) GetMethods() ([]plugins.Method, error) {
	return plugins.DefaultGetMethods(e)
}

func (e *CurrencyConversionApi) CallMethod(methodName string, params ...interface{}) (interface{}, error) {
	return plugins.DefaultCallMethod(e, methodName, params...)
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
