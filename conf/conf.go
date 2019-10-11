package conf

import (
	"Athena/model"
	"Athena/tasks"
	"os"

	"github.com/joho/godotenv"
)

// Init 初始化项目
func Init() {
	// 获取环境变量
	godotenv.Load()

	// 连接数据库
	model.Database(os.Getenv("MYSQL_DSN"))

	//启动定时任务
	tasks.CronJob()
}
