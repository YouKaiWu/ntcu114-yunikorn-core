package users

import (
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/requests"
	"github.com/apache/yunikorn-core/pkg/custom/parser"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"

	sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"

	"fmt"
	"sync"
	"time"
	"container/heap"
)

type User struct{
	totalWaitTime time.Duration // total wait time of all request
	completedRequestCnt int
	unScheduledRequests *requests.Requests 
	currConsumeResource *resources.Resource
	CurrAllocatedCnt int
	waitResource *resources.Resource
	appsRequestResource map[string]*resources.Resource
	sync.RWMutex
}

func NewUser() *User{
	return &User{
		completedRequestCnt: 0,
		unScheduledRequests: requests.NewRequests(),
		currConsumeResource: resources.NewResource(),
		CurrAllocatedCnt: 0,
		waitResource: resources.NewResource(),
		appsRequestResource: make(map[string]*resources.Resource, 0),
	}
}

func (user* User) AddRequest(request *objects.AllocationAsk){
	newRequest := requests.NewRequest(request.GetApplicationID(), request.GetCreateTime(), request.GetAllocationKey())
	heap.Push(user.unScheduledRequests, newRequest)
	_, requestResource := parser.ParseRequestInfo(request)
	user.waitResource.AddTo(requestResource)
	// log.Log(log.Custom).Info("waitresource update")
	// for _, rType := range parser.ResourceType{
	// 	log.Log(log.Custom).Info(fmt.Sprintf("resourceType: %v, quantity: %v", rType, user.waitResource.Resources[rType]))
	// }
}

func (user* User) GetunScheduledRequests() *requests.Requests {
	user.Lock()
	defer user.Unlock()
	return user.unScheduledRequests
}

func (user *User) GetCPUUsage(clusterResource *resources.Resource) float64{ 
	user.Lock()
	defer user.Unlock()
	CPUUsage := float64(user.currConsumeResource.Resources[sicommon.CPU]) / float64(clusterResource.Resources[sicommon.CPU])
	return CPUUsage
}


func (user *User) GetMemoryUsage(clusterResource *resources.Resource) float64{ 
	user.Lock()
	defer user.Unlock()
	memoryUsage := float64(user.currConsumeResource.Resources[sicommon.Memory]) / float64(clusterResource.Resources[sicommon.Memory])
	return memoryUsage
}

// get current used dominant resource share
func (user *User) GetDRS(clusterResource *resources.Resource) (dominantResourceShare float64, dominantResourcesType string){  
	user.Lock()
	defer user.Unlock()
	resourceTypes := []string{sicommon.CPU, sicommon.Memory}
	// compute dominant resource
	dominantResourceShare = 0.0
	dominantResourcesType = sicommon.CPU
	for _, resourceType := range resourceTypes{
		curResourceShare := float64(user.currConsumeResource.Resources[resourceType]) / float64(clusterResource.Resources[resourceType])
		if curResourceShare > dominantResourceShare {
			dominantResourcesType = resourceType
			dominantResourceShare = curResourceShare
		}
	}
	return 
}

// get current demand(waiting queue) dominant resource share
func (user *User) GetDDRS(clusterResource *resources.Resource) (demandDominantResourceShare float64, demandDominantResourcesType string){
	user.Lock()
	defer user.Unlock()
	resourceTypes := []string{sicommon.CPU, sicommon.Memory}
	// compute dominant resource
	demandDominantResourceShare = 0.0
	demandDominantResourcesType = sicommon.CPU
	for _, resourceType := range resourceTypes{
		curResourceShare := float64(user.waitResource.Resources[resourceType]) / float64(clusterResource.Resources[resourceType])
		if curResourceShare > demandDominantResourceShare {
			demandDominantResourcesType = resourceType
			demandDominantResourceShare = curResourceShare
		}
	}
	return 
}

func (user *User) AddWaitTime(createTime time.Time){
	user.totalWaitTime += time.Since(createTime)
}

func (user *User) GetWaitTime() time.Duration{
	return user.totalWaitTime
}

func (user *User) GetCompletedRequestCnt() int{
	return user.completedRequestCnt
}

func (user *User) Allocate(appID string, requestResource *resources.Resource){
	user.Lock()
	defer user.Unlock()
	log.Log(log.Custom).Info(fmt.Sprintf("request allocate, appID:%v", appID))
	if _, exist := user.appsRequestResource[appID]; !exist {
		user.appsRequestResource[appID] = requestResource.Clone()
	}
	user.waitResource.SubFrom(requestResource)
	user.currConsumeResource.AddTo(requestResource)
	user.CurrAllocatedCnt++
}

func (user *User) Release(appID string){
	user.Lock()
	defer user.Unlock()
	log.Log(log.Custom).Info(fmt.Sprintf("request release, appID:%v", appID))
	user.completedRequestCnt++
	user.currConsumeResource.SubFrom(user.appsRequestResource[appID])
	user.CurrAllocatedCnt--
}
