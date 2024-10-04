package users

import (
	"github.com/apache/yunikorn-core/pkg/common/resources"
	// "github.com/apache/yunikorn-core/pkg/log"

	"container/heap"
	// "fmt"
	"sync"
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

func (users *Users) GetUser(username string) *User {
	users.Lock()
	defer users.Unlock()
	return users.users[username]
}

func (users *Users) AddUser(username string) {
	users.Lock()
	defer users.Unlock()
	if _, exist := users.users[username]; !exist {
		users.users[username] = NewUser()
	}
}

// DRS stand for dominant resource share
func (users *Users) GetMinDRSUser(clusterResource *resources.Resource) string {  
	users.Lock()
	defer users.Unlock()
	if len(users.users) == 0 {
		return ""
	}
	usersHeap := NewUsersHeap()
	for username, user := range users.users {
		dominantResourceShare, dominantResourceType := user.GetDRS(clusterResource)
		userInfo := NewUserInfo(username, dominantResourceShare, dominantResourceType)
		heap.Push(usersHeap, userInfo)
	}
	for usersHeap.Len() > 0 {
		curUserInfo := heap.Pop(usersHeap).(*UserInfo)
		curUsername := curUserInfo.username
		// curDRS := curUserInfo.dominantResourceShare
		// curDominantResourceType := curUserInfo.dominantResourceType
		curUser := users.users[curUsername]
		if curUser.unScheduledRequests.Len() > 0 {
			// log.Log(log.Custom).Info(fmt.Sprintf("minUser:%v, dominantResourceShare: %v, resourceType: %v", curUsername, curDRS, curDominantResourceType))
			return curUsername
		}
	}
	return ""
}

// cosidering DRS and DDRS
func (users *Users) GetMinEval(clusterResource *resources.Resource) string {  
	users.Lock()
	defer users.Unlock()
	if len(users.users) == 0 {
		return ""
	}
	usersHeap := NewUsersHeap()
	for username, user := range users.users {
		DRS, DRType := user.GetDRS(clusterResource)
		DDRS, _ := user.GetDDRS(clusterResource)
		eval := 0.9 * DRS - 0.1 * DDRS
		userInfo := NewUserInfo(username, eval, DRType)
		heap.Push(usersHeap, userInfo)
	}
	for usersHeap.Len() > 0 {
		curUserInfo := heap.Pop(usersHeap).(*UserInfo)
		curUsername := curUserInfo.username
		curUser := users.users[curUsername]
		if curUser.unScheduledRequests.Len() > 0 {
			return curUsername
		}
	}
	return ""
}