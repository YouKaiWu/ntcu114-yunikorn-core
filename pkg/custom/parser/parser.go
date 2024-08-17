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

func ParseApp(app *objects.Application, allocationKey string) (appID string, username string, requestResource *resources.Resource) {
	appID = app.ApplicationID
	username = app.GetUser().User
	requestResource = app.GetAllocationAsk(allocationKey).GetAllocatedResource()
	// log.Log(log.Custom).Info(fmt.Sprintf("appID:%v consume resource", appID))
	// for _, rType := range ResourceType{
	// 	log.Log(log.Custom).Info(fmt.Sprintf("resourceType: %v, quantity: %v", rType, requestResource.Resources[rType]))
	// }
	return
}

func ParseAppWithoutResource(app *objects.Application) (appID string, username string) {
	appID = app.ApplicationID
	username = app.GetUser().User
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