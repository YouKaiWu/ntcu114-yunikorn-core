package requests

import(
	"time"
)

type Request struct{
	AppID string;
	CreateTime time.Time
	AllocationKey string
}

func NewRequest(id string, createTime time.Time, allocationKey string) *Request{
	return &Request{
		AppID: id,
		CreateTime: createTime,
		AllocationKey : allocationKey,
	}
}