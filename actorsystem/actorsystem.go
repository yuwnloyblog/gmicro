package actorsystem

import (
	"github.com/yuwnloyblog/gmicro/utils"
)

const (
	NoRpcHost string = "-"
	NoRpcPort int    = 0
)

type ActorSystem struct {
	Name        string
	Host        string
	Prot        int
	sender      *MsgSender
	receiver    *MsgReceiver
	RecvDecoder func([]byte, interface{})
	dispatcher  *ActorDispatcher
}

func NewActorSystemNoRpc(name string) *ActorSystem {
	return NewActorSystem(name, NoRpcHost, NoRpcPort)
}

func NewActorSystem(name, host string, port int) *ActorSystem {
	sender := NewMsgSender()
	dispatcher := NewActorDispatcher(sender)
	receiver := NewMsgReceiver(host, port, dispatcher)

	sender.SetMsgReceiver(receiver)
	system := &ActorSystem{
		Name:       name,
		Host:       host,
		Prot:       port,
		sender:     sender,
		receiver:   receiver,
		dispatcher: dispatcher,
	}
	return system
}

func (system *ActorSystem) LocalActorOf(method string) ActorRef {
	return system.ActerOf(system.Host, system.Prot, method)
}

func (system *ActorSystem) ActerOf(host string, port int, method string) ActorRef {
	uid := utils.GenerateUUIDBytes()
	ref := NewActorRef(host, port, method, uid, system.sender)
	return ref
}

func (system *ActorSystem) CallbackActerOf(ttl int64, newInput NewInput, processor Processor) ActorRef {
	uid := utils.GenerateUUIDBytes()
	ref := NewActorRef(system.Host, system.Prot, "method", uid, system.sender)
	return ref
}

// func (system *ActorSystem) RegisterActorProcessor(method string, newInput NewInput, processor Processor, currentCount int) {

// 	system.dispatcher.RegisterActorProcessor(method, NewExecutor(currentCount, newInput, processor))
// }
func (system *ActorSystem) RegisterActor(method string, actor UntypedActor, concurrentCount int) {
	system.dispatcher.RegisterActor(method, actor, concurrentCount)
}
