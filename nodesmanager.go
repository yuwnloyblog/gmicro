package gmicro

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	//"github.com/yuwnloyblog/commonutils"
	"github.com/yuwnloyblog/commonutils"
	"github.com/yuwnloyblog/gmicro/utils"
)

type Node struct {
	Name    string   `json:"name"`
	Ip      string   `json:"ip"`
	Port    int      `json:"port"`
	Methods []string `json:"methods"`
}

type NodesManager struct {
	basePath          string
	zkConn            *zk.Conn
	nodeMap           map[string]*Node
	method2RingMapStr map[string]string
	ringMap           map[string]*commonutils.ConsistentHash
}

func NewSingleNodeManager(name string, ip string, port int) *NodesManager {
	manager := &NodesManager{}
	node := &Node{
		Name: name,
		Ip:   ip,
		Port: port,
	}
	manager.initHashRingByNodes([]*Node{node})
	return manager
}

func NewNodesManager(basePath string, zkAddress []string) (*NodesManager, error) {
	conn, _, err := zk.Connect(zkAddress, time.Second*5)

	if err != nil {
		return nil, err
	}
	manager := &NodesManager{
		basePath: basePath,
		zkConn:   conn,
	}
	go manager.WatchChildrensChange(basePath+"/nodes", NodesChangeListener)
	return manager, nil
}

func NodesChangeListener(children []string, manager *NodesManager) {
	for _, child := range children {
		fmt.Println("Node:", child)
	}
	//reInit
	manager.initHashRingByZk()
}

func (manager *NodesManager) WatchChildrensChange(path string, listener func(children []string, manager *NodesManager)) {
	for {
		children, _, child_ch, err := manager.zkConn.ChildrenW(path)
		if err != nil {
			fmt.Println("Can not get the list of path [" + path + "]")
			time.Sleep(30 * time.Second)
		} else {
			listener(children, manager)
			<-child_ch
		}
	}
}

func (manager *NodesManager) initHashRingByNodes(nodes []*Node) {
	tmpNodeMap := make(map[string]*Node)
	tmpMethod2Nodes := make(map[string]map[string]bool)
	tmpMethod2RingMapStr := make(map[string]string)
	tmpRingMap := make(map[string]*commonutils.ConsistentHash)

	for _, node := range nodes {
		tmpNodeMap[node.Name] = node

		for _, method := range node.Methods {
			nodeNameMap := tmpMethod2Nodes[method]
			if nodeNameMap == nil {
				nodeNameMap = make(map[string]bool)
				tmpMethod2Nodes[method] = nodeNameMap
			}
			isExist := nodeNameMap[node.Name]
			if !isExist {
				nodeNameMap[node.Name] = true
			}
		}

	}

	for method, nodeNameSet := range tmpMethod2Nodes {
		nodeNameSlice := make([]string, 0, len(nodeNameSet))
		for name, _ := range nodeNameSet {
			nodeNameSlice = append(nodeNameSlice, name)
		}
		sort.Strings(nodeNameSlice)
		namesStr := strings.Join(nodeNameSlice, ":")
		tmpMethod2RingMapStr[method] = namesStr

		hashRing := generateConsistentHashRing(tmpNodeMap, nodeNameSlice)
		tmpRingMap[namesStr] = hashRing
	}
	manager.nodeMap = tmpNodeMap
	manager.method2RingMapStr = tmpMethod2RingMapStr
	manager.ringMap = tmpRingMap
}

func (manager *NodesManager) initHashRingByZk() {
	nodesPath := manager.basePath + "/nodes"
	children, _, err := manager.zkConn.Children(nodesPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	var nodes []*Node
	for _, childName := range children {
		node, err := manager.getNodeWithZk(childName)
		if err == nil {
			nodes = append(nodes, node)
		}
	}
	manager.initHashRingByNodes(nodes)
}

func generateConsistentHashRing(nodeMap map[string]*Node, names []string) *commonutils.ConsistentHash {
	hashRing := commonutils.NewConsistentHash(true)

	for _, name := range names {
		nod := nodeMap[name]
		if nod != nil {
			hashRing.Add(nod.Name, nod.Name, 1)
		}
	}
	return hashRing
}

func (manager *NodesManager) getNodeWithZk(name string) (*Node, error) {
	nodePath := manager.basePath + "/nodes/" + name
	bytes, _, err := manager.zkConn.Get(nodePath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	node := &Node{}
	err = utils.JsonUnMarshal(bytes, node)
	return node, err
}

func (manager *NodesManager) GetTargetNode(method, targetId string) *Node {
	nodeNameStr := manager.method2RingMapStr[method]
	hashRing := manager.ringMap[nodeNameStr]
	if hashRing != nil {
		nod := hashRing.Get(targetId)
		if nod != nil {
			return manager.nodeMap[nod.Name]
		}
	}
	return nil
}

/**
* The path of node data:   /gmicro/clusters/{clusterName}/nodes/{node.Name}
*
**/
func (manager *NodesManager) RegisterSelf(node Node, path string) {
	nodePath := manager.basePath + "/nodes/" + node.Name
	data, _ := utils.JsonMarshal(node)
	isExist, state, _ := manager.zkConn.Exists(nodePath)
	var version int32 = 0
	var err error
	if isExist {
		version = state.Version
		_, err = manager.zkConn.Set(nodePath, data, version)
	} else {
		// flags有4种取值：
		// 0:永久，除非手动删除
		// zk.FlagEphemeral = 1:短暂，session断开则该节点也被删除
		// zk.FlagSequence  = 2:会自动在节点后面添加序号
		// 3:Ephemeral和Sequence，即，短暂且自动添加序号
		_, err = manager.zkConn.Create(nodePath, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	}
	if err != nil {
		log.Fatal(err)
	}
}

/**
* destroy
**/
func (manager *NodesManager) Destroy() {
	manager.zkConn.Close()
}
