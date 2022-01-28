package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data string, PrevHash []byte) *Block {
	block := &Block{Hash: []byte{}, Data: []byte(data), PrevHash: PrevHash, Nonce: 0}

	pow := NewProof(block)

	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

func Genesis() *Block {
	return CreateBlock("Boris Johnson yet to receive Sue Gray report and says it’s ‘total rhubarb’ he authorised Kabul animal airlift – as it happened", []byte{})
}

func (b *Block) Serilize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)

	HandleError(err)
	return res.Bytes()
}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Deserilize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	HandleError(err)

	return &block
}
