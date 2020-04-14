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
// The chain will be implented through a slice of pointers to individual blocks
type BlockChain struct {
	blocks []*Block
}

//Block methods
//DeriveHash will create a hash based on the current blocks
// data and the previous blocks hash
func (b *Block) DeviveHash() {

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

	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{}) // join slices of bytes into 2D slice

}
