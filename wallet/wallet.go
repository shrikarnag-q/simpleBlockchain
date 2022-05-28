package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/bc/utils"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// ***************Wallet*********//
type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

type Transaction struct {
	senderPrivateKey          *ecdsa.PrivateKey
	senderPublicKey           *ecdsa.PublicKey
	senderBlockchainAddress   string
	receiverBlockchainAddress string
	value                     float32
}

func NewWallet() *Wallet {
	// 1. Create ECDSA PrivateKey (32bytes) and PublicKey (64bytes)
	w := new(Wallet)
	privatekey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privatekey
	w.publicKey = &w.privateKey.PublicKey
	// 2. Perform SHA-256 Hashing on PublicKey (32bytes)
	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	// 3. Perform RIPEMD-160 hashing on result of SHA256 (20bytes)
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	// 4. Add version Byte infront of RIPEMD-160 - 0x00 on main network
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])
	// 5. Perform SHA-256 hashing on extended RIPEMD-160 result
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	// 6. Perform SHA-256 hashing on result from 5th step
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	// 7. Take first four bytes from SHA-256 hash for checksum
	chsum := digest6[:4]
	// 8. Add four bytes at the end of the result of extended RIPE-160 from step 4 (25bytes)
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[:21], chsum[:])
	// 9. Convert the result into byte string into BASE58
	address := base58.Encode(dc8)
	w.blockchainAddress = address
	return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress

}

// ********************Transaction in Wallet**********************//
func NewTransaction(privatekey *ecdsa.PrivateKey, publickey *ecdsa.PublicKey, sender string, receiver string, value float32) *Transaction {
	return &Transaction{privatekey, publickey, sender, receiver, value}
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{r, s}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender   string  `json:"sender_blockchain_address"`
		Receiver string  `json:"receiver_blockchain_address"`
		Value    float32 `json:"value"`
	}{
		Sender:   t.senderBlockchainAddress,
		Receiver: t.receiverBlockchainAddress,
		Value:    t.value,
	})
}
