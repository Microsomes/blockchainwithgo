package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/microsomes/blockchainwithgo/blockchain"
)

const bootstrapIRC = "http://www.dal.net:9090/"

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("add - block BLOCK_DATA - add a block to the chain")
	fmt.Println("print - Prints the blocks in the chain ")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Prev. Hash: %x\n", block.PrevHash)
		fmt.Printf("Data $s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)

		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println("")

	}
}

func main() {

}