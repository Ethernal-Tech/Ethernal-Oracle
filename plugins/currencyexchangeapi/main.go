package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"oracle-test/plugins"
	"os"
	"strconv"

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

func (e *CurrencyExchangeApi) GetMethods() []plugins.Method {
	return plugins.DefaulGetMethods(e)
}

func (e *CurrencyExchangeApi) CallMethod(methodName string, params ...interface{}) (interface{}, error) {
	return plugins.DefaultCallMethod(e, methodName, params...)
}

func (e *CurrencyExchangeApi) Exhange_Rate(base string, target string) (float64, error) {
	apiUrl := e.address + fmt.Sprintf("exchange?from=%s&to=%s&q=1.0", base, target)

	request, err := http.NewRequest("GET", apiUrl, nil)
	request.Header.Add("X-RapidAPI-Key", e.api_key)
	request.Header.Add("X-RapidAPI-Host", "currency-exchange.p.rapidapi.com")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, fmt.Errorf("Error making request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("Error making request: %s", err)
	}

	rate, err := strconv.ParseFloat(string(body), 8)

	return rate, nil
}

func main() {}

var ExportPlugin = CurrencyExchangeApi{}
