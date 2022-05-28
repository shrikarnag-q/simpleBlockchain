package main

import (
	"fmt"
	"log"

	"github.com/bc/wallet"
)

func init() {
	log.SetPrefix("Blockchain_Logger: ")
}

func main() {
	// Creating wallet
	w := wallet.NewWallet()
	fmt.Println("w PrivateKey\n", w.PrivateKeyStr())
	fmt.Println("w PublicKey\n", w.PublicKeyStr())
	fmt.Println("w Blockchain Address\n", w.BlockchainAddress())

	x := wallet.NewWallet()
	fmt.Println("x PrivateKey\n", x.PrivateKeyStr())
	fmt.Println("x PublicKey\n", x.PublicKeyStr())
	fmt.Println("x Blockchain Address\n", x.BlockchainAddress())

	// Creating Transaction
	t := wallet.NewTransaction(w.PrivateKey(), w.PublicKey(), w.BlockchainAddress(), x.BlockchainAddress(), 2.0)
	fmt.Printf("Signature: %s\n", t.GenerateSignature())
}
