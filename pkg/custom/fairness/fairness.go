package fairness

import (
	"container/heap"
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/requests"
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

func (fairnessManager *FairnessManager) GetClusterResources() *resources.Resource {
	return fairnessManager.clusterResources.Clone()
}

func (fairnessManager *FairnessManager) NextAppToSchedule() (username string, appId string, allocationKey string) {
	fairnessManager.Lock()
	defer fairnessManager.Unlock()
	tenants := fairnessManager.GetTenants()
	username = tenants.GetMinEval(fairnessManager.clusterResources.Clone())
	if username == "" {
		return "", "", ""
	}
	user := tenants.GetUser(username)
	unScheduledApps := user.GetunScheduledRequests()
	if unScheduledApps.Len() == 0 {
		return username, "", ""
	}
	targetRequest := heap.Pop(unScheduledApps).(*requests.Request)
	heap.Push(unScheduledApps, targetRequest)
	appId = targetRequest.AppID
	allocationKey = targetRequest.AllocationKey
	return 
}
