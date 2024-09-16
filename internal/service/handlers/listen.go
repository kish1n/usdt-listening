package handlers

import (
	"context"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/structs"
	"github.com/kish1n/usdt_listening/internal/data"
	"github.com/kish1n/usdt_listening/internal/service/errors/apierrors"
	"github.com/kish1n/usdt_listening/internal/service/helpers"
)

type Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Index int
	Time  time.Time
}

var Client *ethclient.Client

func ListenForTransfers(w http.ResponseWriter, r *http.Request) {
	ProjectID := os.Getenv("API_KEY")

	Log(r).Info(ProjectID)
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/" + ProjectID)

	if err != nil {
		Log(r).Fatalf("Failed to connect to the Ethereum client: %v", err)
		apierrors.ErrorConstructor(w, *Log(r), err, "Server error", "500", "Server error 500", "Unpredictable behavior")
		return
	}

	Client = client
	Log(r).Infof("Connected to Ethereum client")

	contractAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)

	sub, err := Client.SubscribeFilterLogs(context.Background(), query, logs)

	if err != nil {
		Log(r).Fatalf("Failed to subscribe to logs: %v", err)
		apierrors.ErrorConstructor(w, *Log(r), err, "Server error", "500", "Server error 500", "Unpredictable behavior")
		return
	}

	contractABIJSON, err := helpers.ReadABIFile("/usr/local/bin/contractABI.json")

	if err != nil {
		Log(r).Fatalf("Failed to read contract ABI file: %v", err)
		apierrors.ErrorConstructor(w, *Log(r), err, "Server error", "500", "Server error 500", "Unpredictable behavior")
		return
	}

	contractABI, err := abi.JSON(strings.NewReader(contractABIJSON))

	if err != nil {
		Log(r).Fatalf("Failed to parse contract ABI: %v", err)
		apierrors.ErrorConstructor(w, *Log(r), err, "Server error", "500", "Server error 500", "Unpredictable behavior")
		return
	}

	for {
		select {
		case err := <-sub.Err():
			Log(r).Fatalf("Error: %v", err)
		case vLog := <-logs:
			Log(r).Infof("Log: %v", vLog)

			var transferEvent Transfer
			err := contractABI.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				Log(r).Fatalf("Failed to unpack log: %v", err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			stmt := data.Transaction{
				FromAddress: transferEvent.From.Hex(),
				ToAddress:   transferEvent.To.Hex(),
				Value:       transferEvent.Value.Int64(),
				Id:          helpers.GenerateUUID(),
				CreatedAt:   time.Now().UTC(),
			}

			test := structs.Map(stmt)
			Log(r).Infof("test %s", test)

			err = TransactionQ(r).Insert(stmt)
			if err != nil {
				apierrors.ErrorConstructor(w, *Log(r), err, "Server error", "500", "Server error 500", "Unpredictable behavior")
				return
			}

			Log(r).Infof("Transfer event: from %s to %s for %d tokens", stmt.FromAddress,
				stmt.ToAddress, stmt.Value)
		}
	}
}
