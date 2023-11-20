package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
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

func (g *Goerli) CallMethod(methodName string, paramBytes ...[]byte) ([]byte, error) {
	methodValue := reflect.ValueOf(g).MethodByName(methodName)

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

func (g *Goerli) Eth_blockNumber() ([]byte, error) {
	client, err := ethclient.Dial(g.address)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, blockNumber)
	return bytes, nil
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
