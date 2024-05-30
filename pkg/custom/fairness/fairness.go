package fairness

import (
	"container/heap"

	"github.com/apache/yunikorn-core/pkg/custom/fairness/apps"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/users"
)

type FairnessManager struct {
	tenants         *users.Users
	scheduledApps   map[string]bool 
}

func NewFairnessManager() *FairnessManager {
	return &FairnessManager{
		tenants:         users.NewUsers(),
		scheduledApps:   make(map[string]bool, 0),
	}
}

func (fairnessManager *FairnessManager) GetTenants() *users.Users{
	return fairnessManager.tenants
}

func(fairnessManager *FairnessManager) NextAppToSchedule() (scheduled bool, appId string){
	tenants:= fairnessManager.GetTenants()
	username := tenants.GetMinDRFUser()
	if username == ""{
		return false, ""
	}
	unScheduledApps := tenants.GetUser(username).UnscheduledApps
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