package actorsystem

import (
	"sync"

	"github.com/yuwnloyblog/gmicro/actorsystem/rpc"
)

type ActorDispatcher struct {
	dispatchMap sync.Map
	reqDecoder  ReqDecoder
	msgSender   *MsgSender
}

type ReqDecoder func([]byte, interface{})

func NewActorDispatcher(sender *MsgSender) *ActorDispatcher {
	dispatcher := &ActorDispatcher{
		msgSender: sender,
	}
	return dispatcher
}

func (dispatcher *ActorDispatcher) Dispatch(req *rpc.RpcMessageRequest) {
	targetMethod := req.GetTarMethod()
	obj, ok := dispatcher.dispatchMap.Load(targetMethod)
	if ok {
		executor := obj.(Executor)
		executor.Execute(req, dispatcher.msgSender)
	}

}
func (dispatcher *ActorDispatcher) RegisterActorProcessor(method string, executor Executor) {
	dispatcher.dispatchMap.Store(method, executor)
}
