package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/0xrjman/web3-bot-go/config"
	"github.com/0xrjman/web3-bot-go/pkg/bot"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "bot",
		Short: "Gas reimbursement bot",
	}

	var generateWalletCmd = &cobra.Command{
		Use:   "generate-wallet",
		Short: "Generate a new Ethereum wallet and output the keystore JSON",
		Run:   generateWallet,
	}

	// generateWalletCmd.Flags().StringP("privateKey", "", "", "The password used to encrypt the keystore")
	// generateWalletCmd.Flags().StringP("password", "", "", "The password used to encrypt the keystore")
	rootCmd.AddCommand(generateWalletCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}

	cfg := config.LoadConfig()
	b, err := bot.NewBot(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	b.Run()
}

func generateWallet(cmd *cobra.Command, args []string) {
	fmt.Print("Enter PrivateKey: ")
	privateKeyBytes, _ := term.ReadPassword(int(syscall.Stdin))
	privateKey := string(privateKeyBytes)

	fmt.Print("Enter Password: ")
	passwordBytes, _ := term.ReadPassword(int(syscall.Stdin))
	password := string(passwordBytes)

	ks, err := GenerateKeyStore(privateKey, password)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ks)

	os.Exit(0)
}

func GenerateKeyStore(privateKey string, password string) (*keystore.KeyStore, error) {
	PK, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	ks := keystore.NewKeyStore(".", keystore.StandardScryptN, keystore.StandardScryptP)

	account, err := ks.ImportECDSA(PK, password)
	if err != nil {
		return nil, err
	}

	err = ks.Unlock(account, password)
	if err != nil {
		return nil, err
	}

	return ks, nil
}
