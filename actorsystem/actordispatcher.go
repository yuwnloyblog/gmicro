package actorsystem

import (
	"sync"
	"time"

	"github.com/Jeffail/tunny"
	timewheel "github.com/rfyiamcool/go-timewheel"
	"github.com/yuwnloyblog/gmicro/actorsystem/rpc"
	"github.com/yuwnloyblog/gmicro/logs"
	"github.com/yuwnloyblog/gmicro/utils"
	"google.golang.org/protobuf/proto"
)

const buffersize int = 1024

type ActorDispatcher struct {
	dispatchMap        sync.Map
	callbackMap        sync.Map
	msgSender          *MsgSender
	timer              *timewheel.TimeWheel
	callbackPool       *tunny.Pool
	callbackWraperChan chan wraper
}

func NewActorDispatcher(sender *MsgSender) *ActorDispatcher {
	timer, err := timewheel.NewTimeWheel(1*time.Second, 360)
	if err != nil {
		logs.Error("error when start timewheel of dispatcher")
	}
	dispatcher := &ActorDispatcher{
		msgSender:          sender,
		timer:              timer,
		callbackPool:       tunny.NewCallback(64),
		callbackWraperChan: make(chan wraper, buffersize),
	}
	timer.Start()
	go callbackActorExecute(dispatcher.callbackPool, dispatcher.callbackWraperChan)
	return dispatcher
}

func (dispatcher *ActorDispatcher) Dispatch(req *rpc.RpcMessageRequest) {
	targetMethod := req.GetTarMethod()
	var executor IExecutor

	if targetMethod == "" { //callback actor
		key, err := utils.UUIDStringByBytes(req.Session)
		if err == nil {
			obj, ok := dispatcher.callbackMap.LoadAndDelete(key)
			if ok {
				callbackExecutor := obj.(*CallbackActorExecutor)
				//remove from timer task
				task := callbackExecutor.Task
				if task != nil {
					dispatcher.timer.Remove(task)
				}
				executor = callbackExecutor
			}
		}
	} else {
		obj, ok := dispatcher.dispatchMap.Load(targetMethod)
		if ok {
			executor = obj.(IExecutor)
		}
	}
	if executor != nil {
		executor.Execute(req, dispatcher.msgSender)
	}
}
func (dispatcher *ActorDispatcher) Destroy() {
	if dispatcher.timer != nil {
		dispatcher.timer.Stop()
	}
}
func (dispatcher *ActorDispatcher) RegisterActor(method string, actorCreateFun func() IUntypedActor, concurrentCount int) {
	executor := NewActorExecutor(concurrentCount, actorCreateFun)
	dispatcher.dispatchMap.Store(method, executor)
}

func (dispatcher *ActorDispatcher) AddCallbackActor(session []byte, actor IUntypedActor, ttl int) {
	executor := NewCallbackActorExecutor(dispatcher.callbackPool, dispatcher.callbackWraperChan, actor)
	key, err := utils.UUIDStringByBytes(session)
	if err == nil {
		dispatcher.callbackMap.Store(key, executor)
		task := dispatcher.timer.Add(time.Duration(ttl)*time.Second, func() {
			obj, ok := dispatcher.callbackMap.LoadAndDelete(key)
			if ok {
				executor := obj.(*CallbackActorExecutor)
				executor.doTimeout()
			}
		})
		executor.Task = task
	}
}

func commonExecute(req *rpc.RpcMessageRequest, msgSender *MsgSender, actor IUntypedActor) wraper {
	var sender ActorRef

	srcHost := req.SrcHost
	srcPort := req.SrcPort
	srcMethod := req.SrcMethod
	srcSession := req.Session

	if IsNoSender(req) {
		sender = NoSender
	} else {
		sender = &DefaultActorRef{
			Host:    srcHost,
			Port:    int(srcPort),
			Method:  srcMethod,
			Session: srcSession,
			Sender:  msgSender,
		}
	}

	bytes := req.Data

	input := actor.CreateInputObj()
	proto.Unmarshal(bytes, input)
	return wraper{
		sender: sender,
		msg:    input,
		actor:  actor,
	}
}

type wraper struct {
	sender ActorRef
	msg    proto.Message
	actor  IUntypedActor
}

func callbackActorExecute(pool *tunny.Pool, callbackWraperChan chan wraper) {
	for {
		wrapper := <-callbackWraperChan
		pool.Process(func() {
			actor := wrapper.actor
			actor.SetSender(wrapper.sender)
			actor.OnReceive(wrapper.msg)
		})
	}
}
