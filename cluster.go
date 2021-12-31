package gmicro

import (
	"github.com/yuwnloyblog/gmicro/actorsystem"
	"google.golang.org/protobuf/proto"
)

type IRouteMsg interface {
	proto.Message
	TargetId() string
	Method() string
	Message() proto.Message
}

type Cluster struct {
	actorSystem *actorsystem.ActorSystem
}

func NewCluster(nodename, host string, port int) *Cluster {
	actorSystem := actorsystem.NewActorSystem(nodename, host, port)
	cluster := &Cluster{
		actorSystem: actorSystem,
	}
	return cluster
}

func (cluster *Cluster) RegisterActorProcessor(method string, newInput actorsystem.NewInput, processor actorsystem.Processor, currentCount int) {
	cluster.actorSystem.RegisterActorProcessor(method, newInput, processor, currentCount)
}

func (cluster *Cluster) Route(method string, host string, port int, obj proto.Message) {
	actor := cluster.actorSystem.ActerOf(host, port, method)
	actor.Tell(obj, actorsystem.NoSender)
}
