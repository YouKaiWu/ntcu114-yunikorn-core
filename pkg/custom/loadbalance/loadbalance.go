package loadbalance

import (
	// "github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/loadbalance/formula"
	"github.com/apache/yunikorn-core/pkg/custom/loadbalance/nodes"
	"github.com/apache/yunikorn-core/pkg/custom/parser"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
)

type LoadbalanceManager struct{
	nodes *nodes.Nodes
}

func NewLoadBalanceManager() *LoadbalanceManager{
	return &LoadbalanceManager{
		nodes : nodes.NewNodes(),
	}
}	

func (loadbalanceManager *LoadbalanceManager)GetNodes() *nodes.Nodes{
	return loadbalanceManager.nodes;
}

func (loadbalanceManager *LoadbalanceManager) SelectNode(app *objects.Application, allocationKey string) string{
	_, _, requestResource := parser.ParseApp(app, allocationKey)
	fitInNodes := loadbalanceManager.nodes.GetFitInNodes(requestResource)
	if len(fitInNodes) == 0 {
		return ""
	}
	selectedNode := formula.TOPSIS(requestResource, fitInNodes)
	return selectedNode
}