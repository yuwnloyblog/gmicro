package actorsystem

import (
	"github.com/yuwnloyblog/gmicro/actorsystem/rpc"
)

type MsgReceiver struct {
	host       string
	port       int
	recQueue   chan *rpc.RpcMessageRequest
	dispatcher *ActorDispatcher
}

func NewMsgReceiver(host string, port int, dispatcher *ActorDispatcher) *MsgReceiver {
	rec := &MsgReceiver{
		host:       host,
		port:       port,
		recQueue:   make(chan *rpc.RpcMessageRequest, 10000),
		dispatcher: dispatcher,
	}
	rpcServer := NewRpcServer(host, port, rec)
	go rec.start()
	go rpcServer.Start()
	return rec
}

func (rec *MsgReceiver) Receive(req *rpc.RpcMessageRequest) {
	if req != nil {
		rec.recQueue <- req
	}
}

func (rec *MsgReceiver) isMatch(host string, port int) bool {
	if rec.host == host && rec.port == port {
		return true
	}
	return false
}

func (rec *MsgReceiver) start() {
	for {
		req := <-rec.recQueue
		rec.dispatcher.Dispatch(req)
	}
}
