package loadbalance

import(
	"github.com/apache/yunikorn-core/pkg/custom/loadbalance/nodes"
)

type LoadbalanceManager struct{
	nodes *nodes.Nodes
}

func NewLoadBalanceManager() *LoadbalanceManager{
	return &LoadbalanceManager{
		nodes : nodes.NewNodes(),
	}
}	

func (lb *LoadbalanceManager)GetNodes() *nodes.Nodes{
	return lb.nodes;
}