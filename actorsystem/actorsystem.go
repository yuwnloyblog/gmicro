package actorsystem

import (
	"github.com/yuwnloyblog/gmicro/utils"
)

type ActorSystem struct {
	Name        string
	Host        string
	Port        int
	sender      *MsgSender
	receiver    *MsgReceiver
	TellEncoder func(interface{}) []byte
	RecvDecoder func([]byte, interface{})
	dispatcher  *ActorDispatcher
}

func NewActorSystem(name, host string, port int, encoder func(interface{}) []byte, decoder func([]byte, interface{})) *ActorSystem {
	sender := NewMsgSender()
	dispatcher := NewActorDispatcher(decoder, encoder, sender)
	receiver := NewMsgReceiver(host, port, dispatcher)

	sender.SetMsgReceiver(receiver)
	system := &ActorSystem{
		Name:        name,
		Host:        host,
		Port:        port,
		sender:      sender,
		receiver:    receiver,
		TellEncoder: encoder,
		RecvDecoder: decoder,
		dispatcher:  dispatcher,
	}
	return system
}

func (system *ActorSystem) ActerOf(host string, port int, method string) ActorRef {
	uid := utils.GenerateUUIDBytes()
	ref := NewActorRef(host, port, method, uid, system.TellEncoder, system.sender)
	return ref
}

func (system *ActorSystem) RegisterActorProcessor(method string, newInput NewInput, processor Processor, currentCount int) {
	system.dispatcher.RegisterActorProcessor(method, Executor{
		CurrentCount: currentCount,
		NewInputObj:  newInput,
		Proc:         processor,
	})
}
