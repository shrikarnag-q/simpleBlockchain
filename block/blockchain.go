package block

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bc/utils"
)

// **************structures*****************//

type Block struct {
	timeStamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}
type BlockChain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
	port              uint16
}

type Transaction struct {
	senderBlockchainAddress   string
	receiverBlockchainAddress string
	value                     float32
}

type TransactionRequest struct {
	SenderPublicKey            *string  `json:"sender_public_key"`
	SenderBlockchainAddress    *string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string  `json:"recipient_blockchain_address"`
	Value                      *float32 `json:"value"`
	Signature                  *string  `json:"signature"`
}

// ************gen****************//
const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "I AM A MINER"
	MINING_REWARD     = 1.0
)

// ******************Block Related****************//
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	// fmt.Println(string(m))
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TimeStamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{TimeStamp: b.timeStamp,
		Nonce:        b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
		Transactions: b.transactions,
	})
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timeStamp = time.Now().UnixNano()
	b.previousHash = previousHash
	b.nonce = nonce
	b.transactions = transactions
	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp    %d\n", b.timeStamp)
	fmt.Printf("nonce    %d\n", b.nonce)
	fmt.Printf("previousHash    %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}

}

// ***********Block Chain Related *******************//

func NewBlockChain(blockchainAddress string, port uint16) *BlockChain {
	bc := new(BlockChain)
	b := &Block{}
	bc.blockchainAddress = blockchainAddress
	bc.port = port
	bc.CreateBlock(0, b.Hash())
	return bc
}
func (bc *BlockChain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chains"`
	}{

		Blocks: bc.chain,
	})
}

func (bc *BlockChain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *BlockChain) AddTransaction(sender, receiver string, value float32, senderPublicKey *ecdsa.PublicKey, signature *utils.Signature) bool {
	t := NewTransaction(sender, receiver, value)
	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}
	if bc.VerifyTransactionSignature(senderPublicKey, signature, t) {
		/* //have commented this for now. You can uncomment when live
		if bc.CalculateTotal(sender) < value {
			log.Println("Error: Insufficient Balance")
			return false
		} */
		bc.transactionPool = append(bc.transactionPool, t)
		return true

	} else {
		log.Println("ERROR: Verify Transaction")
	}
	return false

}

func (bc *BlockChain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, signature *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], signature.R, signature.S)
}

func (bc *BlockChain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions, NewTransaction(t.senderBlockchainAddress, t.receiverBlockchainAddress, t.value))

	}
	return transactions
}

func (bc *BlockChain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{
		timeStamp:    0,
		nonce:        nonce,
		previousHash: previousHash,
		transactions: transactions,
	}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros

}

func (bc *BlockChain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce

}

func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *BlockChain) Print() {

	for i, block := range bc.chain {

		fmt.Printf("%s chain  %d %s\n", strings.Repeat("=", 10), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s \n", strings.Repeat("*", 25))

}

func (bc *BlockChain) CalculateTotal(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			if blockchainAddress == t.receiverBlockchainAddress {
				totalAmount += t.value
			}
			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= t.value
			}
		}
	}
	return totalAmount
}

// ************transactions related******************//

func NewTransaction(sender, receiver string, value float32) *Transaction {
	return &Transaction{senderBlockchainAddress: sender, receiverBlockchainAddress: receiver, value: value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 30))
	fmt.Printf("Sender_Blockchain_Address: %s\n", t.senderBlockchainAddress)
	fmt.Printf("Receiver_Blockchain_Address: %s\n", t.receiverBlockchainAddress)
	fmt.Printf("Value:          %1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderBlockchainAddress   string  `json:"sender_blockchain_address"`
		ReceiverBlockchainAddress string  `json:"receiver_blockchain_address"`
		Value                     float32 `json:"value"`
	}{
		SenderBlockchainAddress:   t.senderBlockchainAddress,
		ReceiverBlockchainAddress: t.receiverBlockchainAddress,
		Value:                     t.value,
	})
}

func (tr TransactionRequest) Validate() bool {
	if tr.RecipientBlockchainAddress == nil || tr.SenderBlockchainAddress == nil || tr.Value == nil || tr.SenderPublicKey == nil || tr.Signature == nil {
		return false
	}
	return true
}

// ****************Mining Related ****************//

func (bc *BlockChain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}
