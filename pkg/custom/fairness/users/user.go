package users

import (
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/apps"
	"sync"
)

type User struct{
	UnscheduledApps *apps.Apps 
	CompletedApps map[string]bool 
	AppRequestResource map[string]resources.Resource
	sync.RWMutex
}

func NewUser() *User{
	return &User{
		UnscheduledApps: apps.NewApps(),
		CompletedApps: make(map[string]bool),
		AppRequestResource: make(map[string]resources.Resource, 0),
	}
}