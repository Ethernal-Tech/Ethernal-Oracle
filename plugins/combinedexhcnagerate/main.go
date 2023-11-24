package main

import (
	"fmt"
	"oracle-test/plugins"
	"plugin"
	"reflect"
	"strconv"
)

type CombinedExchange struct {
	CurrencyConversionApi plugins.IPlugin
	CurrencyExchangeApi   plugins.IPlugin
	ExchangeRateApi       plugins.IPlugin
}

func (c *CombinedExchange) Initialize() {
	var err error

	c.CurrencyConversionApi, err = loadPlugin("build/plugins/exchangerateapi.so")
	if err != nil {
		fmt.Println(err)
		return
	}

	c.CurrencyExchangeApi, err = loadPlugin("build/plugins/currencyexchangeapi.so")
	if err != nil {
		fmt.Println(err)
		return
	}

	c.ExchangeRateApi, err = loadPlugin("build/plugins/currencyconversionapi.so")
	if err != nil {
		fmt.Println(err)
		return
	}

	c.CurrencyConversionApi.Initialize()
	c.CurrencyExchangeApi.Initialize()
	c.ExchangeRateApi.Initialize()
}

func (c *CombinedExchange) GetMethods() []plugins.Method {
	structType := reflect.TypeOf(c)

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

func (c *CombinedExchange) CallMethod(methodName string, paramBytes ...[]byte) ([]byte, error) {
	methodValue := reflect.ValueOf(c).MethodByName(methodName)

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

func (c *CombinedExchange) Exhange_Rate(first []byte, second []byte) ([]byte, error) {
	res1, _ := c.CurrencyConversionApi.CallMethod("Exhange_Rate", first, second)
	res2, _ := c.CurrencyExchangeApi.CallMethod("Exhange_Rate", first, second)
	res3, _ := c.ExchangeRateApi.CallMethod("Exhange_Rate", first, second)

	rate1, _ := strconv.ParseFloat(string(res1), 64)
	rate2, _ := strconv.ParseFloat(string(res2), 64)
	rate3, _ := strconv.ParseFloat(string(res3), 64)

	average := (rate1 + rate2 + rate3) / 3

	floatString := strconv.FormatFloat(average, 'f', -1, 64)

	return []byte(floatString), nil
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
