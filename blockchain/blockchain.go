package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BLockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitMainChain() *BlockChain {

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)

	HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()

			err = txn.Set(genesis.Hash, genesis.Serilize())
			HandleError(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err

		} else {
			item, err := txn.Get([]byte("lh"))
			err = item.Value(func(val []byte) error {
				lastHash = val
				return err
			})

			return err

		}
	})

	HandleError(err)

	blockchain := BlockChain{LastHash: lastHash, Database: db}

	return &blockchain
}

func (chain *BlockChain) AddBlock(data string) {

	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleError(err)
		item.Value(func(val []byte) error {
			lastHash = val
			return err
		})
		return err
	})
	HandleError(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serilize())
		HandleError(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		HandleError(err)

		chain.LastHash = newBlock.Hash

		return err
	})

	HandleError(err)

}

func (chain *BlockChain) Iterator() *BLockChainIterator {
	iter := &BLockChainIterator{chain.LastHash, chain.Database}
	return iter
}

func (iter *BLockChainIterator) Next() *Block {
	var block *Block
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		err = item.Value(func(val []byte) error {
			block = Deserilize(val)
			return err
		})
		HandleError(err)
		return err
	})
	HandleError(err)

	iter.CurrentHash = block.PrevHash

	return block

}
