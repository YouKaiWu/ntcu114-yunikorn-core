package fairness

import (
	"container/heap"

	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/apps"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/users"
	// "github.com/apache/yunikorn-core/pkg/log"
	"sync"

	// "fmt"

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
	username := tenants.GetMinDRFUser(fairnessManager.clusterResources.Clone())
	// log.Log(log.Custom).Info(fmt.Sprintf("cluster cpu res:%v\n cluster memory res:%v", fairnessManager.clusterResources.Resources["vcore"], fairnessManager.clusterResources.Resources["memory"]))
	if username == ""{
		return false, ""
	}
	unScheduledApps := tenants.GetUser(username).GetUnScheduledApps()
	if unScheduledApps.Len() == 0{
		return false, ""
	}
	// log.Log(log.Custom).Info(fmt.Sprintf("username:%v", username))
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