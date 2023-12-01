package main

import (
	"fmt"
	"oracle-test/plugins"
	"plugin"
)

type CombinedExchange struct {
	CurrencyConversionApi plugins.IPlugin
	CurrencyExchangeApi   plugins.IPlugin
	ExchangeRateApi       plugins.IPlugin
}

func (c *CombinedExchange) Initialize() error {
	var err error

	c.CurrencyConversionApi, err = loadPlugin("build/plugins/exchangerateapi.so")
	if err != nil {
		return fmt.Errorf("Error loading ExchangeRateAPI %v", err)
	}

	c.CurrencyExchangeApi, err = loadPlugin("build/plugins/currencyexchangeapi.so")
	if err != nil {
		return fmt.Errorf("Error loading CurrencyExchangeAPI %v", err)
	}

	c.ExchangeRateApi, err = loadPlugin("build/plugins/currencyconversionapi.so")
	if err != nil {
		return fmt.Errorf("Error loading CurrencyConversionAPI %v", err)
	}

	c.CurrencyConversionApi.Initialize()
	c.CurrencyExchangeApi.Initialize()
	c.ExchangeRateApi.Initialize()

	return nil
}

func (c *CombinedExchange) GetMethods() ([]plugins.Method, error) {
	return plugins.DefaultGetMethods(c)
}

func (c *CombinedExchange) CallMethod(methodName string, params ...interface{}) (interface{}, error) {
	return plugins.DefaultCallMethod(c, methodName, params...)
}

func (c *CombinedExchange) Exhange_Rate(first string, second string) (float64, error) {
	res1, _ := c.CurrencyConversionApi.CallMethod("Exhange_Rate", first, second)
	res2, _ := c.CurrencyExchangeApi.CallMethod("Exhange_Rate", first, second)
	res3, _ := c.ExchangeRateApi.CallMethod("Exhange_Rate", first, second)

	average := (res1.(float64) + res2.(float64) + res3.(float64)) / 3

	return average, nil
}

func loadPlugin(path string) (plugins.IPlugin, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	sym, err := p.Lookup("ExportPlugin")
	if err != nil {
		return nil, err
	}

	return sym.(plugins.IPlugin), nil
}

func main() {}

var ExportPlugin = CombinedExchange{}
