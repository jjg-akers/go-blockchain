package blockchain

//"fmt"

// TYPES

// Block is the base block type
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

// BlockChain type
// The chain will be represented through a slice of pointers to individual blocks
type BlockChain struct {
	Blocks []*Block
}

//AddBlock will allow us to add a block to the chain
func (chain *BlockChain) AddBlock(data string) {
	// get previous block
	prevBlock := chain.Blocks[len(chain.Blocks)-1]

	//create current blockchain - need the data and the previous Hash
	new := CreateBlock(data, prevBlock.Hash)

	// append new block to chain
	chain.Blocks = append(chain.Blocks, new)

}

//Genesis is a function to create the genisis block
// first block has hard coded values, does not have link to previous hash
func Genesis() *Block {
	// block will have an empty previous hash
	return CreateBlock("Genesis", []byte{})
}

// InitBlockChain will initialize a new chain with the genesis block
func InitBlockChain() *BlockChain {
	genesisBlock := Genesis()

	newChain := &BlockChain{[]*Block{genesisBlock}}
	//return &BlockChain{[]*Block{Genesis()}}
	return newChain
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
