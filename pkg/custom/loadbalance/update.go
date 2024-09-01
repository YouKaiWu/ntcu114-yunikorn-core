package loadbalance

import (
	"time"

	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
)

func (loadbalanceManager *LoadbalanceManager)AddNode(node *objects.Node) {
	if node.GetAttribute("Taints") == "false"{
		loadbalanceManager.nodes.AddNode(node)
		loadbalanceManager.nodesMonitor.AddNode(node.NodeID)
	}
}

func (loadbalanceManager *LoadbalanceManager)SaveExcelFile() {
	loadbalanceManager.nodesMonitor.Save()
}

func (loadbalanceManager *LoadbalanceManager) UpdateNodes(){
	loadbalanceManager.Lock()
	defer loadbalanceManager.Unlock()
	loadbalanceManager.nodesMonitor.Record(time.Now(), loadbalanceManager.nodes)
}