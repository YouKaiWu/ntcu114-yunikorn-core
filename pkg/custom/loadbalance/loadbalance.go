package loadbalance

import (
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/loadbalance/formula"
	"github.com/apache/yunikorn-core/pkg/custom/loadbalance/monitor"
	"github.com/apache/yunikorn-core/pkg/custom/loadbalance/nodes"
	"github.com/apache/yunikorn-core/pkg/custom/parser"
	// "github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
	
	"sync"
	// "fmt"
)

type LoadbalanceManager struct{
	nodes *nodes.Nodes
	nodesMonitor   *monitor.NodesMonitor
	sync.RWMutex
}

func NewLoadBalanceManager() *LoadbalanceManager{
	return &LoadbalanceManager{
		nodes : nodes.NewNodes(),
		nodesMonitor: monitor.NewNodesMonitor(),
	}
}	

func (loadbalanceManager *LoadbalanceManager) SelectNode(request *objects.AllocationAsk) string{
	_, requestResource := parser.ParseRequestInfo(request)
	tmp := resources.MultiplyBy(requestResource, 1.5) // multiply 1.5 times to avoid oomkilled
	fitInNodes := loadbalanceManager.nodes.GetFitInNodes(tmp)
	if len(fitInNodes) == 0 {
		return ""
	}
	selectedNode := formula.TOPSIS(requestResource, fitInNodes)
	return selectedNode
}