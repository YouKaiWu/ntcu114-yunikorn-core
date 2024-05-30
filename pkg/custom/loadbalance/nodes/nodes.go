package nodes 

import(
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
	"github.com/apache/yunikorn-core/pkg/log"
	"fmt"
)

type Nodes []*objects.Node;

func NewNodes() *Nodes{
	tmp := make(Nodes, 0);
	return &tmp;
}

func (nodes *Nodes) AddNode(node *objects.Node){
	*nodes = append(*nodes, node);
	log.Log(log.Custom).Info(fmt.Sprintf("custom added node: %v", node.NodeID));
}