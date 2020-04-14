package main

import (
	"fmt"
	"strconv"

	"github.com/jjg-akers/go-blockchain/blockchain"
)

func main() {
	//fmt.Println("hello")

	chain := blockchain.InitBlockChain()

	chain.AddBlock("First after Gen")
	chain.AddBlock("Second after Gen")
	chain.AddBlock("Third after Gen")

	// lets look at what is in the chain

	for _, block := range chain.Blocks {

		//fmt.Printf("Previous hash: %x\n", block.PrevHash) // format x is for base 16 with lower-case letters

		//fmt.Printf("Previous hash in base 10: %d\n", block.PrevHash)

		fmt.Printf("Data in current block: %s\n", block.Data)
		fmt.Printf("Hash of current block: %x\n", block.Hash)

		// add proof stuff

		pow := blockchain.NewProof(block)
		fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

	}

}

// need to be able to compare the hashes to see how the've changed to determine
// validity
