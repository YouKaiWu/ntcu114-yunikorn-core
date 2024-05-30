package apps

type Apps []*App

func NewApps() *Apps{
	apps := make(Apps, 0);
	return &apps
}

func (apps Apps)Len() int{
	return len(apps)
}

func (apps Apps) Less(i, j int) bool{
	if(apps[i].SubmissionTime.Equal(apps[j].SubmissionTime)){
		return apps[i].Id < apps[j].Id
	}
	return apps[i].SubmissionTime.Before(apps[j].SubmissionTime)
}

func (apps Apps) Swap(i, j int) {
	apps[i], apps[j] = apps[j], apps[i]
}

func (apps *Apps) Push(x interface{}) {
    *apps = append(*apps, x.(*App))
}

func (apps *Apps) Pop() interface{} {
    old := *apps
    n := len(old)
    x := old[n-1]
    *apps = old[0 : n-1]
    return x
}
