package eth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/term"

	"github.com/0xrjman/web3-bot-go/config"
)

type EthKeystore struct {
	Key *keystore.KeyStore
}

func LoadKeyStore(cfg *config.Config) *EthKeystore {
	ks := keystore.NewKeyStore(cfg.KEYSTORE_PATH, keystore.StandardScryptN, keystore.StandardScryptP)
	fmt.Print("Enter Password: ")
	passwordBytes, _ := term.ReadPassword(int(syscall.Stdin))
	password := string(passwordBytes)
	if err := ks.Unlock(accounts.Account{
		Address: common.HexToAddress(cfg.SIGNER),
	}, password); err != nil {
		log.Fatalf("failed to unlock SIGNER %v", err)
	}
	return &EthKeystore{Key: ks}
}

// func CreateKeyStore(cfg *config.Config) *EthKeystore {
// 	fmt.Print("Enter Password: ")
// 	passwordBytes, _ := term.ReadPassword(int(syscall.Stdin))
// 	password := string(passwordBytes)
// 	key, err := GenerateKeyStore(cfg, password)
// 	if err != nil {
// 		log.Fatalf("failed to generate keystore %v", err)
// 	}
// 	return &EthKeystore{Key: key}
// }

func (e *EthKeystore) SendTransaction(rpcClient *RPCClient, to string, amount int) (string, error) {
	// Parse the destination address
	address := common.HexToAddress(to)

	// Convert the amount to Wei (1 ETH = 1e18 Wei)
	value := big.NewInt(int64(amount))
	value.Mul(value, big.NewInt(1e18))

	// Fetch gas price from the network
	gasPriceResult, err := rpcClient.CallRPCMethod(func(client *ethclient.Client) (interface{}, error) {
		return client.SuggestGasPrice(context.Background())
	})
	if err != nil {
		return "", err
	}
	gasPrice, ok := gasPriceResult.(*big.Int)
	if !ok {
		return "", errors.New("failed to assert gas price to *big.Int")
	}

	// Set a static gas limit
	gasLimit := uint64(50000)

	// Fetch nonce from the network
	var nonce uint64
	_, err = rpcClient.CallRPCMethod(func(client *ethclient.Client) (interface{}, error) {
		nonce, err = client.PendingNonceAt(context.Background(), e.Key.Accounts()[0].Address)
		return nil, err
	})
	if err != nil {
		return "", err
	}

	// Create a new transaction
	tx := types.NewTransaction(nonce, address, value, gasLimit, gasPrice, nil)

	// Sign the transaction
	signedTx, err := e.Key.SignTx(accounts.Account{Address: e.Key.Accounts()[0].Address}, tx, nil)
	if err != nil {
		return "", err
	}

	// Send the transaction
	_, err = rpcClient.CallRPCMethod(func(client *ethclient.Client) (interface{}, error) {
		err := client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			return nil, err
		}
		// return the hash of the transaction
		return signedTx.Hash().Hex(), nil
	})
	if err != nil {
		return "", err
	}

	// Return the transaction hash
	return signedTx.Hash().Hex(), nil
}

// func GenerateKeyStore(password string) (*keystore.Key, error) {
// 	privateKey, err := crypto.GenerateKey()
// 	if err != nil {
// 		return nil, err
// 	}

// 	key := keystore.NewKeyFromECDSA(privateKey)
// 	encryptedKey, err := keystore.EncryptKey(key, password, keystore.StandardScryptN, keystore.StandardScryptP)
// 	if err != nil {
// 		return nil, err
// 	}

// 	key, err = keystore.DecryptKey([]byte(encryptedKey), password)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return key, nil
// }
