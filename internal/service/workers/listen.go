package workers

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/structs"
	"github.com/kish1n/usdt_listening/internal/config"
	"github.com/kish1n/usdt_listening/internal/data"
	"github.com/kish1n/usdt_listening/internal/data/pg"
	"github.com/kish1n/usdt_listening/internal/service/helpers"
)

type Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Index int
	Time  time.Time
}

func ListenForTransfers(ctx context.Context, cfg config.Config) {
	ProjectID := os.Getenv("api_key")
	log := cfg.Log()

	log.Info(ProjectID)
	TransactionQ := pg.NewTransaction(cfg.DB().Clone())
	Client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/" + ProjectID)

	if err != nil {
		panic(fmt.Errorf("failed to connect to the Ethereum client: %v", err))
	}

	log.Info("Connected to Ethereum client")

	contractAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)

	sub, err := Client.SubscribeFilterLogs(ctx, query, logs)

	if err != nil {
		panic(fmt.Errorf("failed to subscribe to logs: %v", err))
	}

	contractABIJSON, err := helpers.ReadABIFile("./contractABI.json")

	if err != nil {
		panic(fmt.Errorf("failed to read contract ABI file: %v", err))
	}

	contractABI, err := abi.JSON(strings.NewReader(contractABIJSON))

	if err != nil {
		panic(fmt.Errorf("failed to parse contract ABI: %v", err))
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("Context cancelled, exiting event listener")
				return
			case err := <-sub.Err():
				log.WithError(err).Error("Subscription error")
				return
			case vLog := <-logs:
				log.Info("Log: %v", vLog)

				var transferEvent Transfer
				err := contractABI.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
				if err != nil {
					log.WithError(err).Errorf("Failed to unpack log")
					continue
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
				log.Info("Transaction data: %s", test)

				err = TransactionQ.Insert(stmt)
				if err != nil {
					log.WithError(err).Errorf("Failed to insert transaction")
				}

				log.Infof("Transfer event: from %s to %s for %d tokens", stmt.FromAddress, stmt.ToAddress, stmt.Value)
			}
		}
	}()

	select {
	case <-signalChan:
		log.Info("Received shutdown signal, shutting down...")
	case <-ctx.Done():
		log.Info("Context cancelled, shutting down...")
	}
	log.Info("Process terminated gracefully")
}
