package fairness

import (

	"github.com/apache/yunikorn-core/pkg/common/configs"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/apps"
	"github.com/apache/yunikorn-core/pkg/custom/parser"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
	
	"container/heap"
	"go.uber.org/zap"
)

func (fairnessManager *FairnessManager) ParseUsersInPartitionConfig(conf configs.PartitionConfig) {
	tenants := fairnessManager.GetTenants()
	users := parser.ParseUsersInPartitionConfig(conf)
	for _, user := range users {
		tenants.AddUser(user)
	}
}

// If there is a new tenant's name in the new submitted application, add the username to the fairnessmanager
func (fairnessManager *FairnessManager) ParseUserInApp(app *objects.Application, allocationKey string) {
	fairnessManager.Lock()
	defer fairnessManager.Unlock()
	appID, username := parser.ParseAppWithoutResource(app)
	tenants := fairnessManager.GetTenants()
	tenants.AddUser(username)
	user := tenants.GetUser(username)
	unscheduledApps := user.GetUnScheduledApps()
	newApp := apps.NewApp(appID, app.SubmissionTime, allocationKey)
	heap.Push(unscheduledApps, newApp)
	log.Log(log.Custom).Info("Add application in fair manager", zap.String("user", username), zap.String("applicationID", appID), zap.Int("heap", unscheduledApps.Len()))
}