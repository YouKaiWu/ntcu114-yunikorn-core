package users

import (
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/apps"
	sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"
	"github.com/apache/yunikorn-core/pkg/log"
	"fmt"
	"sync"
)

type User struct{
	unScheduledApps *apps.Apps 
	CompletedApps map[string]bool 
	AppsRequestResource map[string]*resources.Resource
	sync.RWMutex
}

func NewUser() *User{
	return &User{
		unScheduledApps: apps.NewApps(),
		CompletedApps: make(map[string]bool),
		AppsRequestResource: make(map[string]*resources.Resource, 0),
	}
}


func (user* User) GetUnScheduledApps() *apps.Apps {
	user.Lock()
	defer user.Unlock()
	return user.unScheduledApps
}



func (user *User) GetDRF(clusterResource *resources.Resource) float64{
	user.Lock()
	defer user.Unlock()
	resourceTypes := []string{sicommon.CPU, sicommon.Memory}
	userResources :=  make(map[string]float64) 
	// init
	for _, resourceType := range resourceTypes{
		userResources[resourceType] = 0.0
	}
	for appID, completed := range user.CompletedApps{
		if !completed {
			if _, exist := user.AppsRequestResource[appID]; exist{
				user.Release(appID)
			}
		}
	}
	// sum up
	for _, resource := range user.AppsRequestResource{
		for _, resourceType := range resourceTypes{
			userResources[resourceType] += float64(resource.Resources[resourceType]) 
		}
	}
	// compute dominant resource
	drf := 0.0
	for _, resourceType := range resourceTypes{
		tmp := userResources[resourceType] / float64(clusterResource.Resources[resourceType])
		if tmp > drf {
			drf = tmp
		}
	}
	return drf
}

func (user *User) Allocate(appID string, requestResource *resources.Resource){
	user.Lock()
	defer user.Unlock()
	log.Log(log.Custom).Info(fmt.Sprintf("app allocate, id:%v", appID))
	if _, exist := user.AppsRequestResource[appID]; !exist {
		user.AppsRequestResource[appID] = requestResource.Clone()
	}
}

func (user *User) Release(appID string){
	user.Lock()
	defer user.Unlock()
	log.Log(log.Custom).Info(fmt.Sprintf("app release, id:%v", appID))
	if _, exist := user.AppsRequestResource[appID]; exist {
		delete(user.AppsRequestResource, appID)
		user.CompletedApps[appID] = true
	} else {
		user.CompletedApps[appID] = false
	}
	// if _, exist := user.AppsRequestResource[appID]; exist {
		
	// 	// f.GetDRFsWhenComplete(f.GetTenants().GetDRFs(f.clusterResource.Clone()))
	// }
}