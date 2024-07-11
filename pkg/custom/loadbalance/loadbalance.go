package loadbalance

import (
	"github.com/apache/yunikorn-core/pkg/custom/loadbalance/formula"
	"github.com/apache/yunikorn-core/pkg/custom/loadbalance/nodes"
	"github.com/apache/yunikorn-core/pkg/custom/utilization"
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

func (loadbalanceManager *LoadbalanceManager) SelectNode(app *objects.Application) string{
	_, _, requestResource := utilization.ParseApp(app)
	fitInNodes := loadbalanceManager.nodes.GetFitInNodes(requestResource);
	selectedNode := formula.TOPSIS(requestResource, fitInNodes)
	return selectedNode
}