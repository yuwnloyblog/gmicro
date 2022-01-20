package actorsystem

import (
	"google.golang.org/protobuf/proto"
)

type IUntypedActor interface {
	OnReceive(msg proto.Message)
	SetSender(sender ActorRef)
	GetSender() ActorRef
	SetSelf(self ActorRef)
	GetSelf() ActorRef
	OnTimeout()
	CreateInputObj() proto.Message
}

type UntypedActor struct {
	sender ActorRef
	self   ActorRef
}

func (act *UntypedActor) OnReceive(input proto.Message) {
}
func (act *UntypedActor) SetSender(sender ActorRef) {
	act.sender = sender
}
func (act *UntypedActor) GetSender() ActorRef {
	return act.sender
}
func (act *UntypedActor) SetSelf(self ActorRef) {
	act.self = self
}
func (act *UntypedActor) GetSelf() ActorRef {
	return act.self
}
func (act *UntypedActor) OnTimeout() {

}
func (act *UntypedActor) CreateInputObj() proto.Message {
	return nil
}
