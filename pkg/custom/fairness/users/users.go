package users

import (
	"sync"
	// "fmt"
	"github.com/apache/yunikorn-core/pkg/common/resources"
	// "github.com/apache/yunikorn-core/pkg/log"
)

type Users struct {
	users map[string]*User
	sync.RWMutex
}

func NewUsers() *Users {
	return &Users{
		users: make(map[string]*User, 0),
	}
}

func (users *Users) GetUser(name string) *User {
	users.Lock()
	defer users.Unlock()
	return users.users[name]
}

func (users *Users) AddUser(name string) {
	users.Lock()
	defer users.Unlock()
	if _, exist := users.users[name]; !exist {
		users.users[name] = NewUser()
	}
}

func (users *Users) GetMinDRFUser(clusterResource *resources.Resource) string {
	users.Lock()
	defer users.Unlock()
	if len(users.users) == 0 {
		return ""
	}
	minDRF := 1.0 // default max
	minUser := ""
	for userName, user := range users.users {
		drf := user.GetDRF(clusterResource)
		if drf < minDRF {
			minDRF = drf
			minUser = userName
		}
	}
	// log.Log(log.Custom).Info(fmt.Sprintf("mindrf: %v", minDRF))
	return minUser
}
