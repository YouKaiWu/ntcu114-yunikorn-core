package fairness

import (
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/apps"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/users"

	"container/heap"
	"sync"
)

type FairnessManager struct {
	tenants         *users.Users
	scheduledApps   map[string]bool 
	clusterResources *resources.Resource
	nodesCapacity map[string]*resources.Resource
	sync.RWMutex
}

func NewFairnessManager() *FairnessManager {
	return &FairnessManager{
		tenants:         users.NewUsers(),
		scheduledApps:   make(map[string]bool, 0),
		clusterResources: resources.NewResource(),
		nodesCapacity: make(map[string]*resources.Resource, 0),
	}
}

func (fairnessManager *FairnessManager) GetTenants() *users.Users{
	return fairnessManager.tenants
}

func(fairnessManager *FairnessManager) NextAppToSchedule() (scheduled bool, appId string){
	fairnessManager.Lock()
	defer fairnessManager.Unlock()
	tenants:= fairnessManager.GetTenants()
	username := tenants.GetMinDRSUser(fairnessManager.clusterResources.Clone())
	if username == ""{
		return false, ""
	}
	unScheduledApps := tenants.GetUser(username).GetUnScheduledApps()
	if unScheduledApps.Len() == 0{
		return false, ""
	}
	targetApp := heap.Pop(unScheduledApps).(*apps.App)
	if _, exist := fairnessManager.scheduledApps[targetApp.Id]; exist{
		delete(fairnessManager.scheduledApps, targetApp.Id)
		if unScheduledApps.Len() > 0{
			targetApp = heap.Pop(unScheduledApps).(*apps.App)
			heap.Push(unScheduledApps, targetApp)
		}else{
			return false, ""
		}
	}else {
		heap.Push(unScheduledApps, targetApp)
	}
	return true, targetApp.Id
}