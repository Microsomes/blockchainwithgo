package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/microsomes/blockchainwithgo/blockchain"
	"golang.org/x/crypto/ripemd160"
)

const (
	checkSumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)

	blockchain.HandleError(err)

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()

	wallet := Wallet{PrivateKey: private, PublicKey: public}

	return &wallet
}

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	blockchain.HandleError(err)

	publicRipMD := hasher.Sum(nil)

	return publicRipMD

}

func CheckSum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checkSumLength]
}

func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	versionHash := append([]byte{version}, pubHash...)

	checkSum := CheckSum(versionHash)

	fullHash := append(versionHash, checkSum...)

	address := Base58Encode(fullHash)

	return address

}
