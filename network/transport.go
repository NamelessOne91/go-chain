package network

type NetAddr string

type RPC struct {
	From    NetAddr
	Payload []byte
}

type Transport interface {
	Addr() NetAddr
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(NetAddr, []byte) error
}
