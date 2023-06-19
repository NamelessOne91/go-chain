package network

import (
	"fmt"
	"sync"
)

type LocalTransport struct {
	m           sync.RWMutex
	addr        NetAddr
	peers       map[NetAddr]*LocalTransport
	consumeChan chan RPC
}

func NewLocalTransport(addr NetAddr) *LocalTransport {
	return &LocalTransport{
		addr:        addr,
		consumeChan: make(chan RPC, 1024),
		peers:       make(map[NetAddr]*LocalTransport),
	}
}

func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeChan
}

func (t *LocalTransport) Connect(tr *LocalTransport) error {
	t.m.Lock()
	defer t.m.Unlock()

	t.peers[tr.Addr()] = tr
	return nil
}

func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.m.Lock()
	defer t.m.Unlock()

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not send message to %s", t.addr, to)
	}

	peer.consumeChan <- RPC{
		From:    t.addr,
		Payload: payload,
	}

	return nil
}
