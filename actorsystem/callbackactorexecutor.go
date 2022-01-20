package actorsystem

import (
	"github.com/Jeffail/tunny"
	"github.com/rfyiamcool/go-timewheel"
	"github.com/yuwnloyblog/gmicro/actorsystem/rpc"
)

type CallbackActorExecutor struct {
	Task         *timewheel.Task
	wraperChan   chan wraper
	callbackPool *tunny.Pool
	actor        IUntypedActor
}

func NewCallbackActorExecutor(callbackPool *tunny.Pool, wraperChan chan wraper, actor IUntypedActor) *CallbackActorExecutor {
	executor := &CallbackActorExecutor{
		wraperChan:   wraperChan,
		callbackPool: callbackPool,
		actor:        actor,
	}
	return executor
}

func (executor *CallbackActorExecutor) Execute(req *rpc.RpcMessageRequest, msgSender *MsgSender) {
	executor.wraperChan <- commonExecute(req, msgSender, executor.actor)
}

func (executor *CallbackActorExecutor) doTimeout() {
	if executor.actor != nil {
		executor.actor.OnTimeout()
	}
}
