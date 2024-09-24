package requests

type Requests []*Request

func NewRequests() *Requests{
	Requests := make(Requests, 0);
	return &Requests
}

func (Requests Requests)Len() int{
	return len(Requests)
}

func (Requests Requests) Less(i, j int) bool{
	if(Requests[i].CreateTime.Equal(Requests[j].CreateTime)){
		return Requests[i].AppID < Requests[j].AppID
	}
	return Requests[i].CreateTime.Before(Requests[j].CreateTime)
}

func (Requests Requests) Swap(i, j int) {
	Requests[i], Requests[j] = Requests[j], Requests[i]
}

func (Requests *Requests) Push(x interface{}) {
    *Requests = append(*Requests, x.(*Request))
}

func (Requests *Requests) Pop() interface{} {
    old := *Requests
    n := len(old)
    x := old[n-1]
    *Requests = old[0 : n-1]
    return x
}
