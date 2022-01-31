package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger/v3"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "Live  Politics latest news: Met Police tells Sue Gray to only publish ‘minimal’ details on 'partygate' events"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BLockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func DBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func ContinueBlockChain(address string) *BlockChain {
	if DBExists() == false {
		fmt.Println("No existing db exists Create one")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.ValueDir = dbPath
	opts.Logger = nil

	db, err := badger.Open(opts)

	HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleError(err)
		item.Value(func(val []byte) error {
			lastHash = val
			return err
		})
		return err
	})

	HandleError(err)

	chain := BlockChain{LastHash: lastHash, Database: db}

	return &chain
}

func InitBlockchain(address string) *BlockChain {
	var lastHash []byte

	if DBExists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)
	opts.ValueDir = dbPath
	opts.Logger = nil

	db, err := badger.Open(opts)

	HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis Block created")
		err = txn.Set(genesis.Hash, genesis.Serilize())
		HandleError(err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash
		return err
	})

	HandleError(err)

	blockchain := BlockChain{LastHash: lastHash, Database: db}

	return &blockchain

}

func (chain *BlockChain) AddBlock(transactions []*Transaction) {
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

	newBlock := CreateBlock(transactions, lastHash)

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

//find unspent transaction

func (chain *BlockChain) FindUnspentTransaction(address string) []Transaction {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spendOut := range spentTXOs[txID] {
						if spendOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}

			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}

		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs

}

func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UXTOs []TxOutput

	unspentTransactions := chain.FindUnspentTransaction(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UXTOs = append(UXTOs, out)
			}
		}
	}

	return UXTOs
}

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)

	unspendTxs := chain.FindUnspentTransaction(address)

	accumulated := 0

Work:
	for _, tx := range unspendTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}

			}
		}

	}

	return accumulated, unspentOuts

}
