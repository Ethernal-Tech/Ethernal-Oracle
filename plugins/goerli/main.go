package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"oracle-test/plugins"
	"os"

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

func (g *Goerli) Initialize() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file %v", err)
	}

	g.address = os.Getenv("GOERLI_NODE_URL")

	return nil
}

func (g *Goerli) GetMethods() ([]plugins.Method, error) {
	return plugins.DefaultGetMethods(g)
}

func (g *Goerli) CallMethod(methodName string, params ...interface{}) (interface{}, error) {
	return plugins.DefaultCallMethod(g, methodName, params...)
}

func (g *Goerli) Eth_blockNumber() (uint64, error) {
	client, err := ethclient.Dial(g.address)
	if err != nil {
		return 0, err
	}

	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

func (g *Goerli) Eth_getBlockByNumber(params []byte) ([]byte, error) {
	blockNumber := new(big.Int)
	blockNumber.SetBytes(params)

	client, err := ethclient.Dial(g.address)
	if err != nil {
		return nil, err
	}

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
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
		return nil, err
	}

	return blockInfoJSON, nil
}

// Export as plugin
func main() {}

var ExportPlugin = Goerli{}
