package users

import (
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/apps"
	"github.com/apache/yunikorn-core/pkg/log"

	sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"
	
	"fmt"
	"sync"
)

type User struct{
	unScheduledApps *apps.Apps 
	completedApps map[string]bool 
	appsRequestResource map[string]*resources.Resource
	sync.RWMutex
}

func NewUser() *User{
	return &User{
		unScheduledApps: apps.NewApps(),
		completedApps: make(map[string]bool),
		appsRequestResource: make(map[string]*resources.Resource, 0),
	}
}


func (user* User) GetUnScheduledApps() *apps.Apps {
	user.Lock()
	defer user.Unlock()
	return user.unScheduledApps
}



func (user *User) GetDRS(clusterResource *resources.Resource) (dominantResourceShare float64, dominantResourcesType string){  // DRS stand for dominant resource share
	user.Lock()
	defer user.Unlock()
	resourceTypes := []string{sicommon.CPU, sicommon.Memory}
	userResources :=  make(map[string]float64) 
	// init
	for _, resourceType := range resourceTypes{
		userResources[resourceType] = 0.0
	}
	for appID, completed := range user.completedApps{
		if !completed {
			if _, exist := user.appsRequestResource[appID]; exist{
				user.Release(appID)
			}
		}
	}
	// sum up
	for _, resource := range user.appsRequestResource{
		for _, resourceType := range resourceTypes{
			userResources[resourceType] += float64(resource.Resources[resourceType]) 
		}
	}
	// compute dominant resource
	dominantResourceShare = 0.0
	dominantResourcesType = sicommon.CPU
	for _, resourceType := range resourceTypes{
		curResourceShare := userResources[resourceType] / float64(clusterResource.Resources[resourceType])
		if curResourceShare > dominantResourceShare {
			dominantResourcesType = resourceType
			dominantResourceShare = curResourceShare
		}
	}
	return 
}

func (user *User) Allocate(appID string, requestResource *resources.Resource){
	user.Lock()
	defer user.Unlock()
	log.Log(log.Custom).Info(fmt.Sprintf("app allocate, id:%v", appID))
	if _, exist := user.appsRequestResource[appID]; !exist {
		user.appsRequestResource[appID] = requestResource.Clone()
	}
}

func (user *User) Release(appID string){
	user.Lock()
	defer user.Unlock()
	log.Log(log.Custom).Info(fmt.Sprintf("app release, id:%v", appID))
	if _, exist := user.appsRequestResource[appID]; exist {
		delete(user.appsRequestResource, appID)
		user.completedApps[appID] = true
	} else {
		user.completedApps[appID] = false
	}
}