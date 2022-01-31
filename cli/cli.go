package cli

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/microsomes/blockchainwithgo/blockchain"
	"github.com/microsomes/blockchainwithgo/wallet"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS - get the blanace for an address ")
	fmt.Println(" createblockchain -address ADDRESS creates a blockchain and sends the genesis reward to address")
	fmt.Println(" printchain - Prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send amount of coins")
	fmt.Println(" getblock -hash HASH- get block info")
	fmt.Println(" createwallet - Creates a new Wallet")
	fmt.Println(" listaddresses - Lists the addresses in our wallet file")
}

func (cli *CommandLine) listAddresses() {
	wallets, _ := wallet.CreateWallet()
	addresses := wallets.GetAllAddresses()
	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CommandLine) createWalet() {
	wallets, _ := wallet.CreateWallet()
	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Printf("New address is: %s\n", address)
}

func (cli *CommandLine) getBlock(hash string) {
	// byte, err := hex.DecodeString(hash)
	// blockchain.HandleError(err)

	chain := blockchain.ContinueBlockChain("")

	iter := chain.Iterator()

	var foundBlock *blockchain.Block

	for {
		block := iter.Next()

		blockID := hex.EncodeToString(block.Hash)

		if blockID == hash {
			foundBlock = block
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	fmt.Println("--Block Found--")
	fmt.Println("Total TXes:", len(foundBlock.Transactions))
	fmt.Println("Nonce:", foundBlock.Nonce)
	fmt.Println("PrevHash:", hex.EncodeToString(foundBlock.PrevHash))
	fmt.Println("Transaction Hash:", hex.EncodeToString(foundBlock.HashTransactions()))

}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) createBlockchain(address string) {
	chain := blockchain.InitBlockchain(address)
	chain.Database.Close()
	fmt.Println("Finished creating blockchain")

}
func (cli *CommandLine) printChain() {

	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()

	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}

	}
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBlockCmd := flag.NewFlagSet("getblock", flag.ExitOnError)

	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get the blanace for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send the genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	getBlockHash := getBlockCmd.String("hash", "", "Hash of block")

	switch os.Args[1] {
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)
	case "listaddresses":
		err := listAddressCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)
	case "getblock":
		err := getBlockCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if createWalletCmd.Parsed() {
		cli.createWalet()
	}

	if listAddressCmd.Parsed() {
		cli.listAddresses()
	}

	if getBlockCmd.Parsed() {
		if *getBlockHash == "" {
			getBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.getBlock(*getBlockHash)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}

		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)

	}

}
