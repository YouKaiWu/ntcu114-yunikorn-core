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
	currConsumeResource *resources.Resource
	appsRequestResource map[string]*resources.Resource
	sync.RWMutex
}

func NewUser() *User{
	return &User{
		unScheduledApps: apps.NewApps(),
		currConsumeResource: resources.NewResource(),
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

func (user *User) Allocate(appID string, requestResource *resources.Resource){
	user.Lock()
	defer user.Unlock()
	log.Log(log.Custom).Info(fmt.Sprintf("app allocate, id:%v", appID))
	if _, exist := user.appsRequestResource[appID]; !exist {
		user.appsRequestResource[appID] = requestResource.Clone()
	}
	user.currConsumeResource.AddTo(requestResource)
}

func (user *User) Release(appID string){
	user.Lock()
	defer user.Unlock()
	log.Log(log.Custom).Info(fmt.Sprintf("app release, id:%v", appID))
	user.currConsumeResource.SubFrom(user.appsRequestResource[appID])
}
