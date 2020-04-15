package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

// We have two main types of data that needs to be stored
//		Blocks: stored with metadata that descibe all the blocks on the chain
//		Chainstate: stores the state of a chain and all of current unspend transactions

// create a path to our db
const (
	dbPath = "./tmp/blocks"
)

// BlockChain type
// 	The chain will be represented through a slice of pointers to individual blocks
//		but will only store the hash of the last block in the chain
//	Uses the badger DB
type BlockChain struct {
	LastHash []byte
	//Blocks   []*Block
	Database *badger.DB
}

// BlockChainIterator will allow us to be able to visualize data layer
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

//AddBlock will allow us to add a block to the chain
func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	// Get a Read db transaction
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		return err

	})

	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	// with new block created, put new block into db
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)

		err = txn.Set([]byte("lh"), newBlock.Hash)

		// set the in memeory chsin last hash value
		chain.LastHash = newBlock.Hash

		return err
	})

	Handle(err)

}

// InitBlockChain will initialize a new chain with the genesis block
func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath      // where db will store keys and metadata
	opts.ValueDir = dbPath //where db stores values
	opts.Logger = nil

	db, err := badger.Open(opts)
	Handle(err)

	// Access db through Update - allows read/write
	//		View - readonly
	// pass a closure to Update
	//		A closure is a function value that references variables
	//		from outside its body. The function may access and assign to
	//		the referenced variables; in this sense the function is "bound" to the variables.
	err = db.Update(func(txn *badger.Txn) error {
		// check if chain has already been stored
		// 		if there is, create a new instance in memory and get the last hash of blockchain on disk
		//		into memory

		// if not, create genisis block store in DB and then save the hash as last hash
		//		then create a  new instance in memorey

		// Check if block chain exists. If an error is returned, we know a db doesn't exist
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No blockchain found")
			genesis := Genesis()

			fmt.Println("Genesis proved")

			// load into DB with hash as key
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)

			// then set the genesis hash as the last has in the DB
			err = txn.Set([]byte("lh"), genesis.Hash)
			// handle this error by passing it back to the closure

			// save in memory
			lastHash = genesis.Hash

			return err

		} else { // if chain already there
			// get last hash from db
			item, err := txn.Get([]byte("lh"))
			Handle(err)

			// store last hash in memory
			lastHash, err = item.ValueCopy(nil)

			return err
		}
	})
	// outside of update, handle and error that was passed back from closure
	Handle(err)

	// create the new blockchain in memory
	blockchain := BlockChain{
		lastHash,
		db,
	}

	return &blockchain
}

// Iterator will allow us to iterate throught the blockchain
//		need to convert the blockcahin into blockchain iterator
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{
		chain.LastHash,
		chain.Database,
	}

	return iter
}

// Next will return the next block on the chain
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)

		encodedBlock, err := item.ValueCopy(nil) // will return a byte representation of the block

		block = Deserialize(encodedBlock)

		return err

	})

	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}
