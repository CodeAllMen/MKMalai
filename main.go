package main

import (
	"github.com/MobileCPX/PreBaseLib/databaseutil/redis"
	"github.com/MobileCPX/PreBaseLib/splib/click"
	_ "github.com/angui001/MKMalai/initial"
	"github.com/angui001/MKMalai/models"
	_ "github.com/angui001/MKMalai/routers"
	"github.com/astaxie/beego"
)

func init() {
	redis.Open("127.0.0.1", 6379, "")

	// 初始化服务配置
	models.InitServiceConfig()
	// SendClickDataToAdmin()
	// task.SendMtDaily()

	// 每日发送MT数据
	// sendMtTask()
}

func main() {
	beego.Run()
}

// func sendMtTask() {
//	cr := cron.New()
//	cr.AddFunc("0 7 */1 * * ?", SendClickDataToAdmin) // 一个小时存一次点击数据并且发送到Admin
//	cr.Start()
// }

func SendClickDataToAdmin() {
	models.InsertHourClick()

	for _, service := range models.ServiceData {
		click.SendHourData(service.CampID, click.PROD) // 发送有效点击数据
	}

}
