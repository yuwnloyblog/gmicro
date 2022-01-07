package gmicro

import (
	"github.com/yuwnloyblog/gmicro/actorsystem"
	"google.golang.org/protobuf/proto"
)

type Cluster struct {
	actorSystem  *actorsystem.ActorSystem
	nodesManager *NodesManager
}

func NewSingleCluster(nodename string) *Cluster {
	actorsystem := actorsystem.NewActorSystemNoRpc(nodename)
	cluster := &Cluster{
		actorSystem:  actorsystem,
		nodesManager: nil,
	}
	return cluster
}

func NewCluster(nodename, host string, port int, zkAddress []string) *Cluster {
	actorSystem := actorsystem.NewActorSystem(nodename, host, port)
	//add self to server

	//start nodesmanager
	nodesManager, _ := NewNodesManager("/gmicro/clusters/"+nodename, zkAddress)
	cluster := &Cluster{
		actorSystem:  actorSystem,
		nodesManager: nodesManager,
	}

	return cluster
}

func (cluster *Cluster) RegisterActorProcessor(method string, newInput actorsystem.NewInput, processor actorsystem.Processor, currentCount int) {
	cluster.actorSystem.RegisterActorProcessor(method, newInput, processor, currentCount)
}

func (cluster *Cluster) UnicastRouteWithNoSender(method, targetId string, obj proto.Message) {
	nod := cluster.nodesManager.GetTargetNode(method, targetId)
	if nod != nil {
		cluster.baseRouteWithNoSender(method, nod.Ip, nod.Port, obj)
	}
}

func (cluster *Cluster) baseRouteWithNoSender(method string, host string, port int, obj proto.Message) {
	actor := cluster.actorSystem.ActerOf(host, port, method)
	actor.Tell(obj, actorsystem.NoSender)
}

func (cluster *Cluster) RouteOnlyMethod(method string, obj proto.Message) {

}

func (cluster *Cluster) BroadcastWithNoSender(method string, obj proto.Message) {

}
