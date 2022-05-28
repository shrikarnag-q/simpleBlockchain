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
	w := wallet.NewWallet()
	fmt.Println("PrivateKey\n", w.PrivateKeyStr())
	fmt.Println("PublicKey\n", w.PublicKeyStr())
	fmt.Println("Blockchain Address\n", w.BlockchainAddress())

}
