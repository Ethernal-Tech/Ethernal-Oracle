package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"oracle-test/plugins"
	"os"
	"reflect"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

type BlockInfo struct {
	BlockNumber  *big.Int          `json:"blockNumber"`
	BlockHash    string            `json:"blockHash"`
	ParentHash   string            `json:"parentHash"`
	Timestamp    uint64            `json:"timestamp"`
	Transactions []TransactionInfo `json:"transactions"`
}

type TransactionInfo struct {
	TxHash    string   `json:"txHash"`
	GasUsed   uint64   `json:"gasUsed"`
	GasPrice  *big.Int `json:"gasPrice"`
	Value     *big.Int `json:"value"`
	InputData []byte   `json:"inputData"`
}

type Goerli struct {
	address string
}

func (g *Goerli) Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	g.address = os.Getenv("GOERLI_NODE_URL")
}

func (g *Goerli) GetMethods() []plugins.Method {
	structType := reflect.TypeOf(g)

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

func (g *Goerli) CallMethod(methodName string, params ...interface{}) (interface{}, error) {
	methodValue := reflect.ValueOf(g).MethodByName(methodName)

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

func (g *Goerli) Eth_blockNumber() (uint64, error) {
	client, err := ethclient.Dial(g.address)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return blockNumber, nil
}

func (g *Goerli) Eth_getBlockByNumber(params []byte) ([]byte, error) {
	blockNumber := new(big.Int)
	blockNumber.SetBytes(params)

	client, err := ethclient.Dial(g.address)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	blockInfo := BlockInfo{
		BlockNumber: block.Number(),
		BlockHash:   block.Hash().Hex(),
		ParentHash:  block.ParentHash().Hex(),
		Timestamp:   block.Time(),
	}

	for _, tx := range block.Transactions() {
		transactionInfo := TransactionInfo{
			TxHash:    tx.Hash().Hex(),
			GasUsed:   tx.Gas(),
			GasPrice:  tx.GasPrice(),
			Value:     tx.Value(),
			InputData: tx.Data(),
		}
		blockInfo.Transactions = append(blockInfo.Transactions, transactionInfo)
	}

	blockInfoJSON, err := json.Marshal(blockInfo)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return blockInfoJSON, nil
}

// Export as plugin
func main() {}

var ExportPlugin = Goerli{}
