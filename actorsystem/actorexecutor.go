package actorsystem

import (
	"reflect"
	"sync"

	"github.com/Jeffail/tunny"
	"github.com/yuwnloyblog/gmicro/actorsystem/rpc"
)

type IExecutor interface {
	Execute(req *rpc.RpcMessageRequest, msgSender *MsgSender)
}

type ActorExecutor struct {
	wraperChan  chan wraper
	executePool *tunny.Pool
	actorPool   sync.Pool
	actor       UntypedActor
}

func NewActorExecutor(concurrentCount int, actor UntypedActor) *ActorExecutor {
	pool := sync.Pool{
		New: func() interface{} {
			refObj := reflect.TypeOf(actor).Elem()
			objValue := reflect.New(refObj)
			obj := objValue.Interface()
			return obj
		},
	}
	executor := &ActorExecutor{
		wraperChan:  make(chan wraper, buffersize),
		executePool: tunny.NewCallback(concurrentCount),
		actorPool:   pool,
		actor:       actor,
	}
	go actorExecute(executor)
	return executor
}

func (executor *ActorExecutor) Execute(req *rpc.RpcMessageRequest, msgSender *MsgSender) {
	executor.wraperChan <- commonExecute(req, msgSender, executor.actor)
}

func actorExecute(executor *ActorExecutor) {
	for {
		wraper := <-executor.wraperChan
		go executor.executePool.Process(func() {
			actorObj := executor.actorPool.Get()
			actor := actorObj.(UntypedActor)
			actor.SetSender(wraper.sender)
			actor.OnReceive(wraper.msg)
			executor.actorPool.Put(actorObj)
		})
	}
}
