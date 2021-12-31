package actorsystem

import (
	"bot/pkg/gmicro/actorsystem/rpc"
	"sync"
)

type ActorDispatcher struct {
	dispatchMap sync.Map
	reqDecoder  ReqDecoder
	tellEncoder func(interface{}) []byte
	msgSender   *MsgSender
}

type ReqDecoder func([]byte, interface{})

func NewActorDispatcher(decoder ReqDecoder, encoder func(interface{}) []byte, sender *MsgSender) *ActorDispatcher {
	dispatcher := &ActorDispatcher{
		reqDecoder:  decoder,
		tellEncoder: encoder,
		msgSender:   sender,
	}
	return dispatcher
}

func (dispatcher *ActorDispatcher) Dispatch(req *rpc.RpcMessageRequest) {
	targetMethod := req.GetTarMethod()
	obj, ok := dispatcher.dispatchMap.Load(targetMethod)
	if ok {
		executor := obj.(Executor)
		executor.Execute(req, dispatcher.reqDecoder, dispatcher.tellEncoder, dispatcher.msgSender)
	}

}
func (dispatcher *ActorDispatcher) RegisterActorProcessor(method string, executor Executor) {
	dispatcher.dispatchMap.Store(method, executor)
}
