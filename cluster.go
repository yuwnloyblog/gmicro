package gmicro

import (
	"github.com/yuwnloyblog/gmicro/actorsystem"
	"google.golang.org/protobuf/proto"
)

type Cluster struct {
	Name         string
	currentNode  *Node
	actorSystem  *actorsystem.ActorSystem
	nodesManager *NodesManager
	isSingle     bool
}

func NewSingleCluster(clustername, nodename string) *Cluster {
	actorSystem := actorsystem.NewActorSystemNoRpc(nodename)
	//current Node
	curNode := NewNode(nodename, actorSystem.Host, actorSystem.Prot)
	cluster := &Cluster{
		Name:        clustername,
		currentNode: curNode,
		actorSystem: actorSystem,
		isSingle:    true,
	}
	return cluster
}

func NewCluster(clustername, nodename, host string, port int, zkAddress []string) *Cluster {
	actorSystem := actorsystem.NewActorSystem(nodename, host, port)
	//add self to server

	//current Node
	curNode := NewNode(nodename, host, port)
	//start nodesmanager
	nodesManager, _ := NewNodesManager("/gmicro/clusters/"+clustername, zkAddress)
	cluster := &Cluster{
		Name:         clustername,
		currentNode:  curNode,
		actorSystem:  actorSystem,
		nodesManager: nodesManager,
	}
	return cluster
}

func (cluster *Cluster) RegisterActor(method string, actorCreateFun func() actorsystem.IUntypedActor, concurrentCount int) {
	cluster.actorSystem.RegisterActor(method, actorCreateFun, concurrentCount)
	cluster.currentNode.AddMethod(method)
}

func (cluster *Cluster) StartUp() {
	if !cluster.isSingle {
		cluster.nodesManager.RegisterSelf2ZK(*cluster.currentNode)
	}
}

func (cluster *Cluster) getTargetNode(method, targetId string) *Node {
	if cluster.isSingle {
		return cluster.currentNode
	} else {
		return cluster.nodesManager.GetTargetNode(method, targetId)
	}
}

func (cluster *Cluster) getNodeList(method string) []*Node {
	if cluster.isSingle {
		return []*Node{
			cluster.currentNode,
		}
	} else {
		return []*Node{}
	}
}

func (cluster *Cluster) UnicastRouteWithNoSender(method, targetId string, obj proto.Message) {
	nod := cluster.getTargetNode(method, targetId)
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
	nodes := cluster.getNodeList(method)
	for _, node := range nodes {
		cluster.baseRouteWithNoSender(method, node.Ip, node.Port, obj)
	}
}
