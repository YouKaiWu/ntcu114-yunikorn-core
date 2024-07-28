package parser

import (
	"github.com/apache/yunikorn-core/pkg/common/configs"
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/common/security"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"

	sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"

	// "fmt"
	"go.uber.org/zap"
)

var (
	ResourceType        = [2]string{sicommon.CPU, sicommon.Memory}
)


// Parse the vcore and memory in node
func ParseNode(n *objects.Node) (nodeID string, available *resources.Resource, capacity *resources.Resource) {
	nodeID = n.NodeID
	available = resources.NewResource()
	capacity = resources.NewResource()

	resources := n.GetAvailableResource().Resources
	for _, targetType := range ResourceType {
		available.Resources[targetType] = resources[targetType]
	}

	resources = n.GetCapacity().Resources
	for _, targetType := range ResourceType {
		capacity.Resources[targetType] = resources[targetType]
	}
	return
}

func ParseApp(app *objects.Application) (appID string, username string, totalResource *resources.Resource) {
	appID = app.ApplicationID
	username = app.GetUser().User
	totalResource = resources.NewResource()
	// if app.GetAllRequests() == nil {
	// 	log.Log(log.Custom).Info("request is nil")
	// }else{
	// 	log.Log(log.Custom).Info("request exist")
	// }
	for _, request := range app.GetAllRequests(){
		curResource := request.GetAllocatedResource()
		totalResource.AddTo(curResource)
	}
	return
}


func ParseUsersInPartitionConfig(conf configs.PartitionConfig) []string {
	users := make([]string, 0)
	for _, q := range conf.Queues {
		acl, err := security.NewACL(q.SubmitACL)
		if err != nil {
			log.Log(log.Custom).Info("Parsing ACL in fair manager is failed", zap.String("error", err.Error()))
		}
		for user := range acl.GetUsers() {
			log.Log(log.Custom).Info("User in config", zap.String("user", user))
			users = append(users, user)
		}
	}
	return users
}