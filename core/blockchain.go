package core

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Blockchain struct {
	store     Storage
	headers   []*Header
	validator Validator
}

func NewBlockchain(genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		headers: []*Header{},
		store:   NewMemoryStore(),
	}
	bc.validator = NewBlockValidator(bc)

	err := bc.addBlockWithNoValidation(genesis)
	return bc, err
}

func (bc *Blockchain) addBlockWithNoValidation(b *Block) error {
	bc.headers = append(bc.headers, b.Header)

	logrus.WithFields(logrus.Fields{
		"height": b.Height,
		"hash":   b.Hash(BlockHasher{}),
	}).Info("adding new block")

	return bc.store.Put(b)
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}
	return bc.addBlockWithNoValidation(b)
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

func (bc *Blockchain) Height() uint32 {
	// subtract the genesis block to obtain real height
	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}
	return bc.headers[height], nil
}
