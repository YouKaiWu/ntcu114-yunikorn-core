package apps

import(
	"time"
)

type App struct{
	Id string;
	SubmissionTime time.Time
	AllocationKey string
}

func NewApp(id string, sub_time time.Time, allocationKey string) *App{
	return &App{
		Id: id,
		SubmissionTime: sub_time,
		AllocationKey : allocationKey,
	}
}