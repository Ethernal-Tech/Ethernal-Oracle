package main

import (
	"encoding/binary"
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
	firstBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(firstBytes, first)
	var second uint64 = 123
	secondBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secondBytes, second)

	res, err := mathapi.CallMethod("Add_Numbers", firstBytes, secondBytes)
	fmt.Println("Add_Numbers:", binary.BigEndian.Uint64(res))

	// exchangerate api
	exchange, err := loadPlugin("build/plugins/exchangerateapi.so")
	if err != nil {
		fmt.Println(err)
		return
	}

	exchange.Initialize()
	res, err = exchange.CallMethod("Exhange_Rate", []byte("USD"), []byte("EUR"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Exhange_Rate:", string(res))

	// currencyexhange api
	exchange2, err := loadPlugin("build/plugins/currencyexchangeapi.so")
	if err != nil {
		fmt.Println(err)
		return
	}

	exchange2.Initialize()
	res, err = exchange2.CallMethod("Exhange_Rate", []byte("USD"), []byte("EUR"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Exhange_Rate:", string(res))

	// currencyconversion api
	exchange3, err := loadPlugin("build/plugins/currencyconversionapi.so")
	if err != nil {
		fmt.Println(err)
		return
	}

	exchange3.Initialize()
	res, err = exchange3.CallMethod("Exhange_Rate", []byte("USD"), []byte("EUR"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Exhange_Rate:", string(res))

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
	// fmt.Println("Get_All_Leagues:", string(res))

	// // bestapi (doesn't work)
	// best, err := loadPlugin("build/plugins/bestapi.so")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// best.Initialize()
	// res, err = best.CallMethod("whatever", nil)
	// fmt.Println(string(res))

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
	// fmt.Println("Eth_blockNumber: ", binary.BigEndian.Uint64(res))

	blockNumberBytes := big.NewInt(1000000).Bytes()
	res, err = goerli.CallMethod("Eth_getBlockByNumber", blockNumberBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("Eth_getBlockByNumber: ", res)
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
