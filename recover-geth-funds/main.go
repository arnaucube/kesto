package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

var cfg Config

type Config struct {
	Web3 struct {
		Url string
	}
	KeyStore struct {
		Path        string
		Address     string
		Password    string
		KeyJsonPath string
	}
}

var cliCommands = []cli.Command{
	{
		Name:    "info",
		Aliases: []string{},
		Usage:   "get info",
		Action:  cmdInfo,
	},
	{
		Name:    "sendall",
		Aliases: []string{},
		Usage:   "send all eth to address",
		Action:  cmdSendAll,
	},
}

func cmdInfo(c *cli.Context) error {
	if err := mustRead(c); err != nil {
		return err
	}

	ks, acc := loadKeyStore()
	ethSrv := loadWeb3(ks, &acc)
	balance, err := ethSrv.GetBalance(acc.Address)
	if err != nil {
		fmt.Println("error getting balance")
		return err
	}
	fmt.Println("Current balance " + balance.String() + " ETH")
	return nil
}

func cmdSendAll(c *cli.Context) error {
	if err := mustRead(c); err != nil {
		return err
	}

	toAddrStr := c.GlobalString("address")
	if toAddrStr == "" {
		return fmt.Errorf("no address to send the ETH specified")
	}
	fmt.Println("Sending to:", toAddrStr)

	ks, acc := loadKeyStore()
	ethSrv := loadWeb3(ks, &acc)

	toAddr := common.HexToAddress(toAddrStr)
	if err := ethSrv.SendAll(toAddr); err != nil {
		return err
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "recover-geth-funds"
	app.Usage = "cli to send all the funds from a geth keystore to an address"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config"},
		&cli.StringFlag{Name: "address"},
	}

	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, cliCommands...)

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}

// load

const (
	passwdPrefix = "passwd:"
	filePrefix   = "file:"
)

func Assert(msg string, err error) {
	if err != nil {
		fmt.Println(msg, " ", err.Error())
		os.Exit(1)
	}
}

func mustRead(c *cli.Context) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".") // adding home directory as first search path

	if c.GlobalString("config") != "" {
		viper.SetConfigFile(c.GlobalString("config"))
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	return nil
}

func loadKeyStore() (*keystore.KeyStore, accounts.Account) {
	var err error
	var passwd string

	ks := keystore.NewKeyStore(cfg.KeyStore.Path, keystore.StandardScryptN, keystore.StandardScryptP)

	if strings.HasPrefix(cfg.KeyStore.Password, passwdPrefix) {
		passwd = cfg.KeyStore.Password[len(passwdPrefix):]
	} else {
		filename := cfg.KeyStore.Password
		if strings.HasPrefix(filename, filePrefix) {
			filename = cfg.KeyStore.Password[len(filePrefix):]
		}
		passwdbytes, err := ioutil.ReadFile(filename)
		Assert("Cannot read password ", err)
		passwd = string(passwdbytes)
	}

	acc, err := ks.Find(accounts.Account{
		Address: common.HexToAddress(cfg.KeyStore.Address),
	})
	Assert("Cannot find keystore account", err)

	Assert("Cannot unlock account", ks.Unlock(acc, passwd))
	fmt.Println("Keystore and account unlocked successfully:", acc.Address.Hex())

	return ks, acc
}

func loadWeb3(ks *keystore.KeyStore, acc *accounts.Account) *ethService {
	// Create geth client
	url := cfg.Web3.Url
	hidden := strings.HasPrefix(url, "hidden:")
	if hidden {
		url = url[len("hidden:"):]
	}
	passwd, err := readPassword(cfg.KeyStore.Password)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	ethsrv := newEthService(url, ks, acc, cfg.KeyStore.KeyJsonPath, passwd)
	if hidden {
		fmt.Println("Connection to web3 server opened", "url", "(hidden)")
	} else {
		fmt.Println("Connection to web3 server opened", "url", cfg.Web3.Url)
	}
	return ethsrv
}

func readPassword(configPassword string) (string, error) {
	var passwd string
	if strings.HasPrefix(cfg.KeyStore.Password, passwdPrefix) {
		passwd = cfg.KeyStore.Password[len(passwdPrefix):]
	} else {
		filename := cfg.KeyStore.Password
		if strings.HasPrefix(filename, filePrefix) {
			filename = cfg.KeyStore.Password[len(filePrefix):]
		}
		passwdbytes, err := ioutil.ReadFile(filename)
		if err != nil {
			return passwd, err
		}
		passwd = string(passwdbytes)
	}
	return passwd, nil
}

// eth

type ethService struct {
	ks       *keystore.KeyStore
	acc      *accounts.Account
	client   *ethclient.Client
	KeyStore struct {
		Path     string
		Password string
	}
}

func newEthService(url string, ks *keystore.KeyStore, acc *accounts.Account, keystorePath, password string) *ethService {
	client, err := ethclient.Dial(url)
	if err != nil {
		fmt.Println("Can not open connection to web3 (config.Web3.Url: " + url + ")\n" + err.Error() + "\n")
		os.Exit(0)
	}

	service := &ethService{
		ks:     ks,
		acc:    acc,
		client: client,
		KeyStore: struct {
			Path     string
			Password string
		}{
			Path:     keystorePath,
			Password: password,
		},
	}

	return service
}

func (ethSrv *ethService) GetBalance(address common.Address) (*big.Float, error) {
	balance, err := ethSrv.client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}
	ethBalance := gweiToEth(balance)
	return ethBalance, nil
}

func gweiToEth(g *big.Int) *big.Float {
	f := new(big.Float)
	f.SetString(g.String())
	e := new(big.Float).Quo(f, big.NewFloat(math.Pow10(18)))
	return e
}

func (ethSrv *ethService) SendAll(toAddr common.Address) error {
	balance, err := ethSrv.client.BalanceAt(context.Background(), ethSrv.acc.Address, nil)
	if err != nil {
		return err
	}
	ethBalance := gweiToEth(balance)
	fmt.Println("Current balance:", ethBalance, "ETH")

	nonce, err := ethSrv.client.PendingNonceAt(context.Background(), ethSrv.acc.Address)
	if err != nil {
		return err
	}
	fmt.Println("Nonce:", nonce)
	gasLimit := uint64(21000)
	gasLimitBI := big.NewInt(int64(gasLimit))
	gasPrice, err := ethSrv.client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	value := new(big.Int).Sub(balance, new(big.Int).Mul(gasLimitBI, gasPrice))
	fmt.Println("balance", gweiToEth(balance), "ETH")
	fmt.Println(" tosend", gweiToEth(value), "ETH (substracting the fees)")

	confirmed := askConfirmation()
	if !confirmed {
		return fmt.Errorf("operation cancelled")
	}
	fmt.Println("operation confirmed")

	var data []byte
	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, data)

	chainID, err := ethSrv.client.NetworkID(context.Background())
	if err != nil {
		return err
	}

	signedTx, err := ethSrv.ks.SignTx(*ethSrv.acc, tx, chainID)
	if err != nil {
		return err
	}

	err = ethSrv.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())

	return nil
}

func askConfirmation() bool {
	var s string

	fmt.Printf("(y/N): ")
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}
