package users

type UserInfo struct{
	username string
	dominantResourceShare float64 
	dominantResourceType string
}

func NewUserInfo(username string, dominantResourceShare float64, dominantResourceType string) *UserInfo{
	return &UserInfo{
		username: username,
		dominantResourceShare: dominantResourceShare,
		dominantResourceType: dominantResourceType,
	}
}

type UsersHeap []*UserInfo

func NewUsersHeap() *UsersHeap{
	usersHeap := make(UsersHeap, 0);
	return &usersHeap
}

func (usersHeap UsersHeap)Len() int{
	return len(usersHeap)
}

func (usersHeap UsersHeap) Less(i, j int) bool{
	if(usersHeap[i].dominantResourceShare == usersHeap[j].dominantResourceShare){
		return usersHeap[i].username < usersHeap[j].username
	}
	return usersHeap[i].dominantResourceShare < usersHeap[j].dominantResourceShare
}

func (usersHeap UsersHeap) Swap(i, j int) {
	usersHeap[i], usersHeap[j] = usersHeap[j], usersHeap[i]
}

func (usersHeap *UsersHeap) Push(element interface{}) {
    *usersHeap = append(*usersHeap, element.(*UserInfo))
}

func (usersHeap *UsersHeap) Pop() interface{} {
    old := *usersHeap
    n := len(old)
    element := old[n-1]
    *usersHeap = old[0 : n-1]
    return element
}
