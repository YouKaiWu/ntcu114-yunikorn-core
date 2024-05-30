package apps

import(
	"time"
)

type App struct{
	Id string;
	SubmissionTime time.Time
}

func NewApp(id string, sub_time time.Time) *App{
	return &App{
		Id: id,
		SubmissionTime: sub_time,
	}
}