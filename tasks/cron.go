package tasks

import (
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/robfig/cron"
)

var Cron *cron.Cron

func Run(job func() error) {
	from := time.Now().UnixNano()
	err := job()
	to := time.Now().UnixNano()
	jobName := runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name()
	if err != nil {
		fmt.Printf("%s error: %dms\n", jobName, (to-from)/int64(time.Millisecond))
	} else {
		fmt.Printf("%s success: %dms\n", jobName, (to-from)/int64(time.Millisecond))
	}
}

func CronJob() {
	if Cron == nil {
		Cron = cron.New()
	}

	Cron.AddFunc("0 */10 * * * ?", func() { Run(RefreshWeibo) })
	Cron.AddFunc("0 0 20 ? * 0,3", func() { Run(WarnAbyss) })
	//Cron.AddFunc("0 0 0 * * *", func() { Run(RestartDailyRank) })
	Cron.Start()

	fmt.Println("Cronjob start.....")
}
