package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// this will hold the proof of work algorithm

// we want to force the network to do work to add blocks to the chain
// the work must be hard, but the valication of work done must be easy

// Step 1:
//	Take data from the block

// Step 2: Create the nonce (just a counter) which starts at 0

// step 3: Create a hash of the data plus the nonce

// Step 4: Check that the hash meets a set of requirements
//		This is where the dificulty lies
//		if the hash meets the req's, we use it
//		other wise we create a new hash until we meet the req's

//		the dificulty can go up over time to make it more difficult

// Difine Requirements
//	1. The first few bytes must contain 0s

// For this implementation we will use a constant dificulty
//		in a real blockchain, an alg. would slowely increment the difficulty
//		to account for increased miners on the network and the increased computing power of
//		new technologies

const Difficulty = 16

// Types

//ProofOfWork type
//	The block is the current block in the chain
//	The Target is a number that represents the requirement which is derived from the difficulty
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProof a func to produce a new proof
func NewProof(b *Block) *ProofOfWork {
	// cast 1 to a new big int
	target := big.NewInt(1)

	// the use Lsh to left shift the target by 256 (number of bytes in one of the hashes)
	// minus the difficulty
	target.Lsh(target, uint(256-Difficulty))

	// put the block and the target into an instance of ProofOfWork

	pow := &ProofOfWork{
		b,
		target,
	}

	return pow
}

// now we want to take the block and create a new hash
// Function to combine the current nonce with the previous hash see if
// we meet the requirements
//
func (pow *ProofOfWork) InitData(nonce int) []byte {
	// Combine the block's previous hash and its data
	// we also need to add in the nonce because we need to combine our
	//	block data with the counter
	// we also need to add in the difficulty
	data := bytes.Join([][]byte{
		pow.Block.PrevHash,
		pow.Block.Data,
		ToHex(int64(nonce)),
		ToHex(int64(Difficulty)),
	},
		[]byte{},
	)
	return data

}

// Run is the main computational function
// 		this will loop until the correct hash is found
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	// want an infinite loop
	nonce := 0

	// prepare the data and hash into sha256
	// then convert to big int
	// then compare to target big int which is stored in the POW structuer

	for nonce < math.MaxInt64 { // a very large number
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 { // -1 means are hash is smaller than the target which means we signed the block
			break

		} else {

			nonce++
		}
	}

	fmt.Println()

	return nonce, hash[:]
}

// Validate function
//	 after we run POW run() func we will have the correct nonce that will allow us to derive the hash
//	 that meets the target we needed
//		then run the cycle one more time to show that this hash is actually valid
// so finidng the nonce is hard, but validating is easy
func (pow *ProofOfWork) Validate() bool {

	// set up big in version of harsh
	var intHash big.Int

	// get data
	data := pow.InitData(pow.Block.Nonce)

	//convert data into a shash
	hash := sha256.Sum256(data)

	// convert the hash into bigInt
	intHash.SetBytes(hash[:])

	// compare to target to see if block is valid
	return intHash.Cmp(pow.Target) == -1

}

// ToHex is a utility function to decode a number into bytes
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
