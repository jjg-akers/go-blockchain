package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/jjg-akers/go-blockchain/blockchain"
)

// Commandline is a struct to handle user interaction
type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println(" add -block <Block_Data> : adds a block to the chain")
	fmt.Println(" print - Prints the block in the chain")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit() // this will exit the application but does it by shutting down the goroutine
		// this is needed to properly shut down the Badger DB
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Block Added...")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in current block: %s\n", block.Data)
		fmt.Printf("Hash of current block: %x\n", block.Hash)

		// add proof stuff
		pow := blockchain.NewProof(block)
		fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			// the genesis block doesn't have a prevhash so it's length will be 0
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	// create some flags for command line tool
	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data") // subset flag

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	// when flags are parsed, they return a bool
	//	so if addblockcmd has been parsed... do something
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			// if they pass empty data
			addBlockCmd.Usage()
			runtime.Goexit()

		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func main() {
	//fmt.Println("hello")
	defer os.Exit(0) // ensure gracefull shutdown

	// Create the blockchain
	chain := blockchain.InitBlockChain()

	// This defer will only execture if the go channel exits properly
	//	thats why we need to use runtime.goexit instead of os exit
	defer chain.Database.Close()

	cli := CommandLine{chain}

	cli.run()

}

// need to be able to compare the hashes to see how the've changed to determine
// validity
