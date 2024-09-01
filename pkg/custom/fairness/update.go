package fairness

import (
	"fmt"

	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/apps"
	"github.com/apache/yunikorn-core/pkg/custom/parser"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"

	// "fmt"
	"container/heap"
	"time"

	"go.uber.org/zap"
)

func (fairnessManager *FairnessManager) UpdateScheduledApp(app *objects.Application, allocationKey string) {
	fairnessManager.Lock()
	defer fairnessManager.Unlock()
	appID, username, requestResource := parser.ParseApp(app, allocationKey)
	user := fairnessManager.GetTenants().GetUser(username)
	dominantResourceShare, dominantResourceType := user.GetDRS(fairnessManager.clusterResources.Clone())
	log.Log(log.Custom).Info(fmt.Sprintf("updated application:[appID: %v, username: %v, dominantResourceShare: %v, dominantResourceType: %v]", appID, username, dominantResourceShare, dominantResourceType))
	user.Allocate(appID, requestResource)
	fairnessManager.tenantsMonitor.Record(time.Now(), fairnessManager.tenants, fairnessManager.clusterResources.Clone())

	if unScheduledApps := user.GetUnScheduledApps(); unScheduledApps.Len() == 0 {
		log.Log(log.Custom).Error("Non existed app update", zap.String("app", appID), zap.String("user", username))
	} else {
		remainingApps := make([]*apps.App, 0)
		for unScheduledApps.Len() > 0 {
			app := heap.Pop(unScheduledApps).(*apps.App)
			id := app.Id
			if app.AllocationKey != allocationKey {
				remainingApps = append(remainingApps, app)
			} else {
				log.Log(log.Custom).Info("Delete app", zap.String("appid", id), zap.Int("heap", unScheduledApps.Len()))
			}
		}
		for _, element := range remainingApps {
			heap.Push(unScheduledApps, element)
		}
	}
}

func (fairnessManager *FairnessManager) AddNode(nodeID string, capacity *resources.Resource) {
	fairnessManager.Lock()
	defer fairnessManager.Unlock()
	tmp := fairnessManager.clusterResources.Clone()
	if curCapacity, nodeExist := fairnessManager.nodesCapacity[nodeID]; nodeExist {
		if !resources.StrictlyGreaterThanOrEquals(curCapacity, capacity) {
			tmp = resources.Sub(tmp, curCapacity)
			tmp = resources.Add(tmp, capacity)
		}
	} else {
		fairnessManager.nodesCapacity[nodeID] = capacity
		tmp = resources.Add(tmp, capacity)
	}
	fairnessManager.clusterResources = tmp
}

func (fairnessManager *FairnessManager) RemoveNode(nodeID string) {
	fairnessManager.Lock()
	defer fairnessManager.Unlock()
	if curCapacity, nodeExist := fairnessManager.nodesCapacity[nodeID]; nodeExist {
		fairnessManager.clusterResources = resources.Sub(fairnessManager.clusterResources, curCapacity)
		delete(fairnessManager.nodesCapacity, nodeID)
	}
}

func (fairnessManager *FairnessManager) AddCompletedApp(appID string, username string) {
	fairnessManager.Lock()
	defer fairnessManager.Unlock()
	// log.Log(log.Custom).Info(fmt.Sprintf("app complete, appId:%v", appID))
	user := fairnessManager.GetTenants().GetUser(username)
	user.Release(appID)
	fairnessManager.tenantsMonitor.Record(time.Now(), fairnessManager.tenants, fairnessManager.clusterResources.Clone())
}

func (fairnessManager *FairnessManager) SaveExcelFile(){
	fairnessManager.tenantsMonitor.Save()
}