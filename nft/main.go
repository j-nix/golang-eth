package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Config struct {
	InfuraApiUrl     string `json:"infura_api_url"`
	Value            int64  `json:"value"`
	WalletPrivateKey string `json:"wallet_private_key"`
	ContractAddress  string `json:"contract_address"`
	GasLimit         int64  `json:"gas_limit"`
	ABI              string `json:"abi"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func GetKey(key string) (*ecdsa.PrivateKey, common.Address, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return privateKey, fromAddress, nil
}

func checkTransactionReceipt(client *ethclient.Client, _txHash string) int {
	txHash := common.HexToHash(_txHash)
	tx, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return -1
	}

	return int(tx.Status)
}

func main() {
	config := LoadConfiguration("config.json")
	ctx := context.Background()

	// Dial new client
	client, err := ethclient.DialContext(ctx, config.InfuraApiUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Get value to spend? Or quantity?
	value := big.NewInt(config.Value)

	// Get current suggested gas price
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Returns ECDSA private key for use in TXN and the address it belongs to (send mint here)
	privateKey, fromAddress, _ := GetKey(config.WalletPrivateKey)
	fmt.Printf("Address: %s\n", fromAddress)

	// Get the contract address (hex)
	toAddress := common.HexToAddress(config.ContractAddress)
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Nonce is the number of total TXNs an address has completed, we must be sequential in this number as it is the blockchain!
	nonce, err := client.NonceAt(ctx, fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

	// ABI is basically the spec in JSON format of a contract, detailling all methods and input/output values and types
	abiP, err := abi.JSON(strings.NewReader(config.ABI))
	if err != nil {
		log.Fatal(err)
	}

	data, err := abiP.Pack(
		"mint",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new TXN request object using struct provided by eth-golang library
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      big.NewInt(config.GasLimit).Uint64(),
		GasPrice: gasPrice,
		Data:     data,
	})

	// Sign the transaction as us (private key - our wallet is asking for it)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Inject the transaction into the pending pool on the eth network
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())

	// Wait until we have some kind of receipt from the TXN (note this cannot return PENDING state, have to check that manually)
	for {
		transactionStatus := checkTransactionReceipt(client, signedTx.Hash().Hex())
		fmt.Printf("tx status: %d\n", transactionStatus)

		if transactionStatus == 1 {
			break
		}
	}

	fmt.Println("tx confirmed!")
}
