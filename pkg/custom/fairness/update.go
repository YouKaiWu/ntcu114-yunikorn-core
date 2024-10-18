package fairness

import (
	"fmt"

	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/requests"
	"github.com/apache/yunikorn-core/pkg/custom/parser"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"

	// "fmt"
	"container/heap"
	// "time"

	"go.uber.org/zap"
)

// process request

// parse new request info
func (fairnessManager *FairnessManager) ParseRequest(request *objects.AllocationAsk, username string) {
	fairnessManager.Lock()
	defer fairnessManager.Unlock()
	appID, _ := parser.ParseRequestInfo(request)
	tenants := fairnessManager.GetTenants()
	tenants.AddUser(username)
	fairnessManager.tenantsMonitor.AddUser(username)
	user := tenants.GetUser(username)
	user.AddRequest(request)
	log.Log(log.Custom).Info("Add new request in fair manager", zap.String("user", username), zap.String("applicationID", appID), zap.Int("total reqeust:", user.GetunScheduledRequests().Len()))
}

func (fairnessManager *FairnessManager) UpdateScheduledRequest(request *objects.AllocationAsk, username string) {
	fairnessManager.Lock()
	defer fairnessManager.Unlock()
	appID, requestResource := parser.ParseRequestInfo(request)
	user := fairnessManager.GetTenants().GetUser(username)
	user.AddWaitTime(request.GetCreateTime())
	dominantResourceShare, dominantResourceType := user.GetDRS(fairnessManager.GetClusterResources())
	log.Log(log.Custom).Info(fmt.Sprintf("updated application:[appID: %v, username: %v, dominantResourceShare: %v, dominantResourceType: %v]", appID, username, dominantResourceShare, dominantResourceType))
	user.Allocate(appID, requestResource)
	fairnessManager.tenantsMonitor.Record(fairnessManager.tenants, fairnessManager.GetClusterResources())
	fairnessManager.tenantsMonitor.RecordScheduleInterval(fairnessManager.tenants, appID)
	if unScheduledRequests := user.GetunScheduledRequests(); unScheduledRequests.Len() == 0 {
		log.Log(log.Custom).Error("Non existed request update", zap.String("appID: ", appID), zap.String("user: ", username))
	} else {
		remainingRequests := make([]*requests.Request, 0)
		for unScheduledRequests.Len() > 0 {
			curRequest := heap.Pop(unScheduledRequests).(*requests.Request)
			id := curRequest.AppID
			if curRequest.AllocationKey != request.GetAllocationKey() {
				remainingRequests = append(remainingRequests, curRequest)
			} else {
				log.Log(log.Custom).Info("Delete request", zap.String("appId", id), zap.Int("remaining request:", unScheduledRequests.Len()))
			}
		}
		for _, element := range remainingRequests {
			heap.Push(unScheduledRequests, element)
		}
	}
}

func (fairnessManager *FairnessManager) AddCompletedRequest(appID string, username string) {
	fairnessManager.Lock()
	defer fairnessManager.Unlock()
	// log.Log(log.Custom).Info(fmt.Sprintf("app complete, appId:%v", appID))
	user := fairnessManager.GetTenants().GetUser(username)
	user.Release(appID)
	fairnessManager.tenantsMonitor.Record(fairnessManager.tenants, fairnessManager.GetClusterResources())
}

// process node

func (fairnessManager *FairnessManager) AddNode(nodeID string, capacity *resources.Resource) {
	fairnessManager.Lock()
	defer fairnessManager.Unlock()
	tmp := fairnessManager.GetClusterResources()
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

// process excel file

func (fairnessManager *FairnessManager) SaveExcelFile(){
	fairnessManager.tenantsMonitor.Save(fairnessManager.GetTenants())
}