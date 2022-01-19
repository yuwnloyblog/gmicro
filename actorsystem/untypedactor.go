package actorsystem

import "google.golang.org/protobuf/proto"

type UntypedActor interface {
	OnReceive(msg proto.Message)
	SetSender(sender ActorRef)
	GetSender() ActorRef
	SetSelf(self ActorRef)
	GetSelf() ActorRef
	OnTimeout()
	CreateInputObj() proto.Message
}
