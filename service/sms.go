/**
  create by yy on 2020/5/14
*/

package service

import (
	"fmt"
	"github.com/angui001/MKMalai/libs"
	"github.com/angui001/MKMalai/models"
	"github.com/astaxie/beego/httplib"
)

// 发送短信函数
func SendMTSMS(reqData *models.ReqData, serviceConfig *models.ServiceInfo) {
	var (
		err    error
		result []byte
	)

	sendUrl := "http://mis.etracker.cc/mcppush/mcppush.aspx"

	postData := make(map[string]interface{})
	// 构造 要 post的数据
	postData["User"] = "foxseek"
	postData["Pass"] = "seek0804"
	// 0是发送文本
	postData["Type"] = 0
	postData["To"] = reqData.Msisdn
	postData["Text"] = libs.EscapeQueryParam("This param should be url encode")
	postData["From"] = reqData.ShortCode
	postData["Telcoid"] = reqData.Telcoid
	postData["Keyword"] = serviceConfig.Keyword
	postData["SecKey"] = ""
	postData["Charge"] = "Subscription"
	postData["Moid"] = reqData.Moid
	postData["Price"] = 2

	req := httplib.Post(sendUrl)

	if _, err = req.JSONBody(postData); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
	}

	if result, err = req.Bytes(); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
	}

	fmt.Println(fmt.Sprintf("send text result: %v", string(result)))
}
