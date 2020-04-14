package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// TYPES

// Base Block type
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

// BlockChain type
// The chain will be represented through a slice of pointers to individual blocks
type BlockChain struct {
	blocks []*Block
}

//AddBlock will allow us to add a block to the chain
func (chain *BlockChain) AddBlock(data string) {
	// get previous block
	prevBlock := chain.blocks[len(chain.blocks)-1]

	//create current blockchain - need the data and the previous Hash
	new := CreateBlock(data, prevBlock.Hash)

	// append new block to chain
	chain.blocks = append(chain.blocks, new)

}

//Genesis is a function to create the genisis block
// first block has hard coded values, does not have link to previous hash
func Genesis() *Block {
	// block will have an empty previous hash
	return CreateBlock("Genesis", []byte{})
}

//InitBlockChain will initialize a new chain with the genesis block
func InitBlockChain() *BlockChain {
	genesisBlock := Genesis()

	newChain := &BlockChain{[]*Block{genesisBlock}}
	//return &BlockChain{[]*Block{Genesis()}}
	return newChain
}

//Block methods
//DeriveHash will create a hash based on the current blocks
// data and the previous blocks hash
func (b *Block) DeriveHash() {

	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{}) // join slices of bytes into 2D slice
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

// CreateBlock will take in a string of data and the previous hash
// and outputs a pointer to a new block
func CreateBlock(data string, prevHash []byte) *Block {
	// Create reference to a new block object
	block := &Block{
		[]byte{},
		[]byte(data),
		prevHash,
	}

	block.DeriveHash()

	return block
}

func main() {
	fmt.Println("hello")

	chain := InitBlockChain()

	chain.AddBlock("First after Gen")
	chain.AddBlock("Second after Gen")
	chain.AddBlock("Third after Gen")

	// lets look at what is in the chain

	for _, block := range chain.blocks {
		//fmt.Printf("Previous hash: %x\n", block.PrevHash) // format x is for base 16 with lower-case letters

		//fmt.Printf("Previous hash in base 10: %d\n", block.PrevHash)

		fmt.Printf("Data in current block: %s\n", block.Data)
		fmt.Printf("Hash of current block: %x\n", block.Hash)

	}

}

// need to be able to compare the hashes to see how the've changed to determine
// validity
