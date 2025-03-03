package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Index        int
	TimeStamp    string
	Transactions []string
	PrevHash     string
	Hash         string
	Proof        int
}

type Blockchain struct {
	Chain []Block
}

func CalculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%v%s%d", block.Index, block.TimeStamp, block.Transactions, block.PrevHash, block.Proof)
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

func (bc *Blockchain) CreateGenesisBlock() {
	genesisBlock := Block{
		Index:        0,
		TimeStamp:    time.Now().String(),
		Transactions: []string{"Genesis Block"},
		PrevHash:     "",
		Proof:        100,
	}
	genesisBlock.Hash = CalculateHash(genesisBlock)
	bc.Chain = append(bc.Chain, genesisBlock)
}

func (bc *Blockchain) GetLastBlock() Block {
	return bc.Chain[len(bc.Chain)-1]
}

func ProofOfWork(lastProof int) int {
	proof := 0
	for {
		hashAttempt := sha256.Sum256([]byte(fmt.Sprintf("%d%d", lastProof, proof)))
		hashString := hex.EncodeToString(hashAttempt[:])
		if hashString[:4] == "0000" {
			break
		}
		proof++
	}
	return proof
}

func (bc *Blockchain) AddBlock(transactions []string) {
	lastBlock := bc.GetLastBlock()
	newProof := ProofOfWork(lastBlock.Proof)

	newBlock := Block{
		Index:        lastBlock.Index + 1,
		TimeStamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     lastBlock.Hash,
		Proof:        newProof,
	}
	newBlock.Hash = CalculateHash(newBlock)
	bc.Chain = append(bc.Chain, newBlock)
}

func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Chain); i++ {
		prevBlock := bc.Chain[i-1]
		currBlock := bc.Chain[i]

		if currBlock.Hash != CalculateHash(currBlock) {
			return false
		}
		if currBlock.PrevHash != prevBlock.Hash {
			return false
		}
	}
	return true
}

func main() {
	blockchain := Blockchain{}
	blockchain.CreateGenesisBlock()

	blockchain.AddBlock([]string{"Varun sent Aryan 5 BTC"})
	blockchain.AddBlock([]string{"Rahul sent Varun 5 BTC"})

	fmt.Println("Blockchain valid:", blockchain.IsValid())

	for _, block := range blockchain.Chain {
		fmt.Printf("\nIndex: %d\nTimestamp: %s\nTransactions: %v\nPrevHash: %s\nHash: %s\nProof: %d\n",
			block.Index, block.TimeStamp, block.Transactions, block.PrevHash, block.Hash, block.Proof)
	}
}
