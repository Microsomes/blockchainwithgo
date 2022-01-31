package main

import (
	"os"

	"github.com/microsomes/blockchainwithgo/cli"
)

const bootstrapIRC = "http://www.dal.net:9090/"

func main() {

	// wallet.MakeWallet()
	// wallet.MakeWallet().Address()
	// fmt.Println(hex.EncodeToString(wallet.Address()))

	defer os.Exit(0)

	cmd := cli.CommandLine{}

	cmd.Run()

}
