package fairness

import (
	"container/heap"
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/apps"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/monitor"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/users"
	"sync"
)

type FairnessManager struct {
	tenants          *users.Users
	clusterResources *resources.Resource
	nodesCapacity    map[string]*resources.Resource
	tenantsMonitor   *monitor.TenantsMonitor

	sync.RWMutex
}

func NewFairnessManager() *FairnessManager {
	return &FairnessManager{
		tenants:          users.NewUsers(),
		clusterResources: resources.NewResource(),
		nodesCapacity:    make(map[string]*resources.Resource, 0),
		tenantsMonitor: monitor.NewTenantsMonitor(),
	}
}

func (fairnessManager *FairnessManager) GetTenants() *users.Users {
	return fairnessManager.tenants
}

func (fairnessManager *FairnessManager) NextAppToSchedule() (appId string, allocationKey string) {
	fairnessManager.Lock()
	defer fairnessManager.Unlock()
	tenants := fairnessManager.GetTenants()
	username := tenants.GetMinDRSUser(fairnessManager.clusterResources.Clone())
	if username == "" {
		return "", ""
	}
	user := tenants.GetUser(username)
	unScheduledApps := user.GetUnScheduledApps()
	if unScheduledApps.Len() == 0 {
		return "", ""
	}
	targetApp := heap.Pop(unScheduledApps).(*apps.App)
	heap.Push(unScheduledApps, targetApp)
	return targetApp.Id, targetApp.AllocationKey
}

func (fairnessManager *FairnessManager) SaveExcelFile(){
	fairnessManager.tenantsMonitor.Save()
}