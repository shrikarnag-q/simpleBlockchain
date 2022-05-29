package main

import (
	"fmt"
	"log"

	"github.com/bc/block"
	"github.com/bc/wallet"
)

func init() {
	log.SetPrefix("Blockchain_Logger: ")
}

func main() {
	// Creating wallet
	minerWallet := wallet.NewWallet()
	fmt.Println("minerWallet Blockchain Address\n", minerWallet.BlockchainAddress())

	personA := wallet.NewWallet()
	fmt.Println("personA Blockchain Address\n", personA.BlockchainAddress())

	personB := wallet.NewWallet()
	fmt.Println("personB Blockchain Address\n", personB.BlockchainAddress())

	var value float32 = 2.0

	// *************Creating Transaction********************//

	//Creating transaction on the Wallet side
	t := wallet.NewTransaction(personA.PrivateKey(), personA.PublicKey(), personA.BlockchainAddress(), personB.BlockchainAddress(), value)
	fmt.Printf("Signature: %s\n", t.GenerateSignature())

	//Creating transaction on the blockchain node side
	blockChain := block.NewBlockChain(minerWallet.BlockchainAddress())
	isAdded := blockChain.AddTransaction(personA.BlockchainAddress(), personB.BlockchainAddress(), value, personA.PublicKey(), t.GenerateSignature())
	log.Println("Is it Added? ", isAdded)

	// ***************Miner transactions********//
	blockChain.Mining()
	blockChain.Print()
	fmt.Printf("Wallet of PersonA %1f\n", blockChain.CalculateTotal(personA.BlockchainAddress()))
	fmt.Printf("Wallet of PersonB %1f\n", blockChain.CalculateTotal(personB.BlockchainAddress()))
	fmt.Printf("Wallet of miner %1f\n", blockChain.CalculateTotal(minerWallet.BlockchainAddress()))

}
