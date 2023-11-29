package main

import (
	"fmt"
	"math/big"
	"oracle-test/plugins"
	"plugin"
)

func main() {

	// mathapi api
	mathapi, err := loadPlugin("build/plugins/mathapi.so")
	if err != nil {
		fmt.Println(err)
		return
	}

	var first uint64 = 123
	var second uint64 = 123

	res, err := mathapi.CallMethod("Add_Numbers", first, second)
	fmt.Println("Add_Numbers:", res.(uint64))

	// combined api
	exchange, err := loadPlugin("build/plugins/combinedexhcnagerate.so")
	if err != nil {
		fmt.Println(err)
		return
	}

	exchange.Initialize()
	methods := exchange.GetMethods()
	fmt.Println(methods)
	res, err = exchange.CallMethod("Exhange_Rate", "USD", "EUR")
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("Exhange_Rate:", res)

	// livescore api
	sportsdb, err := loadPlugin("build/plugins/livescore.so")
	if err != nil {
		fmt.Println(err)
		return
	}

	sportsdb.Initialize()
	res, err = sportsdb.CallMethod("Get_All_Leagues")
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("Get_All_Leagues:", res)

	// bet365 api
	bet365, err := loadPlugin("build/plugins/bet365.so")
	if err != nil {
		fmt.Println(err)
		return
	}

	bet365.Initialize()
	res, err = bet365.CallMethod("Get_Premier_League_Quotas")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Get_Premiere_League_Quotas:", res)

	// goerli
	goerli, err := loadPlugin("build/plugins/goerli.so")
	if err != nil {
		fmt.Println(err)
		return
	}

	goerli.Initialize()
	res, err = goerli.CallMethod("Eth_blockNumber")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Eth_blockNumber: ", res)

	blockNumberBytes := big.NewInt(1000000).Bytes()
	res, err = goerli.CallMethod("Eth_getBlockByNumber", blockNumberBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Eth_getBlockByNumber: ", res)
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
