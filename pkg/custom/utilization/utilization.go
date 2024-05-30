package utilization

import (
	"fmt"

	"github.com/apache/yunikorn-core/pkg/common/configs"
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/common/security"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
	sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"

	"go.uber.org/zap"
)

var (
	ResourceType        = [2]string{sicommon.CPU, sicommon.Memory}
)


// Parse the vcore and memory in node
func ParseNode(n *objects.Node) (nodeID string, avialble *resources.Resource, capacity *resources.Resource) {
	nodeID = n.NodeID
	avialble = resources.NewResource()
	capacity = resources.NewResource()

	resources := n.GetAvailableResource().Resources
	for _, targetType := range ResourceType {
		avialble.Resources[targetType] = resources[targetType]
	}

	resources = n.GetCapacity().Resources
	for _, targetType := range ResourceType {
		capacity.Resources[targetType] = resources[targetType]
	}
	return
}

func ParseApp(app *objects.Application) (appID string, username string, resourceQuantity *resources.Resource) {
	appID = app.ApplicationID
	username = app.GetUser().User
	resourceQuantity = resources.NewResource()
	totalResource := resources.NewResource()
	if app.GetAllRequests() == nil {
		log.Log(log.Custom).Info("request is nil")
	}
	for _, request := range app.GetAllRequests(){
		log.Log(log.Custom).Info("request exist")
		curResource := request.GetAllocatedResource()
		totalResource.AddTo(curResource)
	}
	log.Log(log.Custom).Info(fmt.Sprintf("appId:%v, the following resource of this app", appID))
	for key, val := range totalResource.Resources{
		log.Log(log.Custom).Info(fmt.Sprintf("key: %v; val: %v", key, val))
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