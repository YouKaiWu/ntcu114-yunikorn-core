package nodes 

import(
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"

	"fmt"
)

type Nodes []*objects.Node;

func NewNodes() *Nodes{
	tmp := make(Nodes, 0);
	return &tmp;
}

func (nodes *Nodes) AddNode(node *objects.Node){
	*nodes = append(*nodes, node);
	log.Log(log.Custom).Info(fmt.Sprintf("custom add node: %v", node.NodeID));
	// log.Log(log.Custom).Info(fmt.Sprintf("cur node: %v", node.NodeID));
	// for key, val := range node.GetAttributes(){
	// 	log.Log(log.Custom).Info(fmt.Sprintf("key:%v, val: %v", key, val))
	// }
}

func (nodes *Nodes) GetFitInNodes(requestResource *resources.Resource) Nodes{
	var fitInNodes Nodes 
	for _, node := range *nodes{
		if node.IsSchedulable() && node.CanAllocate(requestResource) {
			fitInNodes = append(fitInNodes, node)
		}
	}
	return fitInNodes
}