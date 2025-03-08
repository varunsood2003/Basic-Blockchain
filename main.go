package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Block struct {
	Index        int      `json:"index"`
	Timestamp    string   `json:"timestamp"`
	Transactions []string `json:"transactions"`
	PrevHash     string   `json:"prev_hash"`
	Hash         string   `json:"hash"`
	Proof        int      `json:"proof"`
}

type Blockchain struct {
	Chain       []Block
	PendingTxns []string
	mu          sync.Mutex
}

func CalculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%v%s%d", block.Index, block.Timestamp, block.Transactions, block.PrevHash, block.Proof)
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

func (bc *Blockchain) CreateGenesisBlock() {
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().Format(time.RFC3339),
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

func (bc *Blockchain) MineBlock() Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	lastBlock := bc.GetLastBlock()
	newProof := ProofOfWork(lastBlock.Proof)

	newBlock := Block{
		Index:        lastBlock.Index + 1,
		Timestamp:    time.Now().Format(time.RFC3339),
		Transactions: bc.PendingTxns,
		PrevHash:     lastBlock.Hash,
		Proof:        newProof,
	}
	newBlock.Hash = CalculateHash(newBlock)
	bc.Chain = append(bc.Chain, newBlock)
	bc.PendingTxns = []string{}
	return newBlock
}

func (bc *Blockchain) AddTransaction(txn string) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.PendingTxns = append(bc.PendingTxns, txn)
}

func main() {
	r := gin.Default()
	blockchain := Blockchain{}
	blockchain.CreateGenesisBlock()

	r.GET("/chain", func(c *gin.Context) {
		c.JSON(200, gin.H{"chain": blockchain.Chain})
	})

	r.POST("/transactions", func(c *gin.Context) {
		var req struct {
			Transaction string `json:"transaction"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}
		blockchain.AddTransaction(req.Transaction)
		c.JSON(200, gin.H{"message": "Transaction added"})
	})

	r.GET("/mine", func(c *gin.Context) {
		newBlock := blockchain.MineBlock()
		c.JSON(200, gin.H{"message": "Block Mined!", "block": newBlock})
	})

	r.Run(":8080")
}
