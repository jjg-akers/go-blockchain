package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

//"fmt"

// TYPES

// Block is the base block type
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

//Genesis is a function to create the genisis block
// first block has hard coded values, does not have link to previous hash
func Genesis() *Block {
	// block will have an empty previous hash
	return CreateBlock("Genesis", []byte{})
}

// since we are deriving the hash within the proof of work function, we
//		don't need a seperate derive hash funciton

//DeriveHash will create a hash based on the current blocks
// data and the previous blocks hash
// func (b *Block) DeriveHash() {

// 	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{}) // join slices of bytes into 2D slice
// 	hash := sha256.Sum256(info)
// 	b.Hash = hash[:]
// }

// CreateBlock will take in a string of data and the previous hash
// and outputs a pointer to a new block
func CreateBlock(data string, prevHash []byte) *Block {
	// Create reference to a new block object
	block := &Block{
		[]byte{},
		[]byte(data),
		prevHash,
		0,
	}

	// We want to run the proof of work alg
	// on every block that is created
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	//block.DeriveHash()

	return block
}

// Since Badger DB only accepts slicres of bytes we need a serialize and deserialize function

// Serialize
func (b *Block) Serialize() []byte {
	var res bytes.Buffer // Results

	// create an encoder on out results buffer
	encoder := gob.NewEncoder(&res)

	// this will let us call encode on the block itself
	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes() // this will just be a byte representation of the block
}

//Deserialize will take in a slice of bytes and transform it back into block representaiton
func Deserialize(data []byte) *Block {
	var block Block

	// Create a byte reader and pass it to decoder
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	// Return the reference
	return &block
}

//Handle will handle errors
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
