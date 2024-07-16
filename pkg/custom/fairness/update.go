package fairness

import (
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/apps"
	"github.com/apache/yunikorn-core/pkg/custom/parser"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
	
	"container/heap"
	"go.uber.org/zap"
)



func (fairnessManager *FairnessManager) UpdateScheduledApp(app *objects.Application) {
	log.Log(log.Custom).Info("update scheduled app")
	appID, username, requestResource := parser.ParseApp(app)
	fairnessManager.scheduledApps[appID] = true
	user := fairnessManager.GetTenants().GetUser(username)
	user.Allocate(appID, requestResource)

	if unScheduledApps := user.GetUnScheduledApps(); unScheduledApps.Len() == 0 {
		log.Log(log.Custom).Error("Non existed app update", zap.String("app", appID), zap.String("user", username))
	} else {
		remainingApps := make([]*apps.App, 0)
		for unScheduledApps.Len() > 0 {
			app := heap.Pop(unScheduledApps).(*apps.App)
			id := app.Id
			if _, exist := fairnessManager.scheduledApps[id]; !exist {
				log.Log(log.Custom).Info("Delete app is not in the heap", zap.String("appid", id))
				remainingApps = append(remainingApps, app)
			} else {
				delete(fairnessManager.scheduledApps, id)
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
	user := fairnessManager.GetTenants().GetUser(username)
	user.Release(appID)
}
