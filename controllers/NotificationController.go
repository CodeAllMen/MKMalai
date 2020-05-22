/**
  create by yy on 2020/4/24
*/

package controllers

import (
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib"
	"github.com/MobileCPX/PreBaseLib/splib/admindata"
	"github.com/MobileCPX/PreBaseLib/splib/common"
	"github.com/MobileCPX/PreBaseLib/splib/mo"
	"github.com/MobileCPX/PreBaseLib/splib/notification"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/MobileCPX/PreBaseLib/util"
	"github.com/angui001/MKMalai/models"
	"github.com/angui001/MKMalai/service"
	"strconv"
	"strings"
)

type NotificationController struct {
	BaseController
}

// 只有订阅的回传
func (n *NotificationController) Mo() {
	msisdn := n.GetString("from")
	text := n.GetString("text")
	times := n.GetString("time")
	// 这两个 同时只会出现一个
	msgid := n.GetString("msgid")
	moid := n.GetString("moid")

	shortCode := n.GetString("shortcode")
	telcoid := n.GetString("telcoid")
	refid := n.GetString("refid")

	reqData := new(models.ReqData)

	fmt.Println("request url and param: ", n.Ctx.Input.URI())

	// 这里注意，要通过电话号码来判断唯一性
	reqData.Msisdn = msisdn
	reqData.Text = text
	reqData.Time = times
	reqData.Msgid = msgid
	reqData.Moid = moid
	reqData.ShortCode = shortCode
	reqData.Telcoid = telcoid
	reqData.Refid = refid

	fmt.Println()

	reqData.Insert()

	moT := new(mo.Mo)
	var moBase = common.MoBase{}
	moBase.Keyword = strings.ToUpper(reqData.Keyword)
	moBase.ShortCode = reqData.ShortCode

	// track id
	trackID, _ := strconv.Atoi(reqData.Refid)
	moBase.TrackID = int64(trackID)

	// 判断数据是否是重复数据
	subID := reqData.ShortCode + "-" + reqData.Keyword + "-" + reqData.Msisdn

	_, _ = moT.GetMoBySubscriptionID(subID)
	if moT.ID != 0 {
		// 重复数据
		return
	}

	track := new(models.AffTrack)
	track.TrackID = int64(trackID)
	_ = track.GetOne(tracking.ByTrackID)
	moBase.ServiceID = ""

	if track.TrackID != 0 {
		moBase.Track = track.Track
		moBase.ServiceID = track.ServiceID
		moBase.TrackID = track.TrackID
	} else {
		moBase.ServiceID = reqData.ShortCode + strings.ToUpper(reqData.Keyword)
	}

	var isExist bool
	serviceConf := models.ServiceInfo{} // 服务配置
	serviceConf, isExist = n.serviceConfig(moBase.ServiceID)

	if isExist {
		moBase.Operator = serviceConf.Operator
		moBase.Price = serviceConf.Price
		moBase.Country = serviceConf.Country
	}

	moBase.Msisdn = reqData.Msisdn
	moBase.SubscriptionID = subID
	// 存入MO数据
	notificationType := ""
	moT, notificationType = splib.InsertMO(moBase, false, true, "MKMalai")

	// 订阅成功之后要发送短信给对方
	// 开始发送短信
	fmt.Println("开始发送短信")

	service.SendMTSMS(reqData, &serviceConf)

	fmt.Println("发送短信完成")

	// 订阅成功后将数据传到Admin服务器
	if notificationType != "" {
		notify := new(notification.Notification)
		notify.NotificationType = notificationType
		notify.SubscriptionId = moBase.SubscriptionID
		notify.TransactionID = moBase.SubscriptionID + "-" + util.RandomString(4)
		notify.ServiceID = track.ServiceID
		notify.Insert()

		nowTime, _ := util.GetNowTime()

		sendNoti := new(admindata.Notification)

		sendNoti.PostbackPrice = moT.PostbackPrice

		sendNoti.OfferID = moT.OfferID
		sendNoti.SubscriptionID = moT.SubscriptionID
		sendNoti.ServiceID = moT.ServiceID
		sendNoti.ClickID = moT.ClickID
		sendNoti.Msisdn = moT.Msisdn
		sendNoti.CampID = serviceConf.CampID
		sendNoti.PubID = moT.PubID
		sendNoti.PostbackStatus = moT.PostbackStatus
		sendNoti.PostbackMessage = moT.PostbackMessage
		sendNoti.TransactionID = notify.TransactionID
		sendNoti.Keyword = moT.Keyword
		sendNoti.ShortCode = moT.ShortCode
		sendNoti.AffName = moT.AffName
		if sendNoti.AffName == "" {
			sendNoti.AffName = "未知"
		}
		sendNoti.Operator = moT.Operator

		sendNoti.Sendtime = nowTime
		sendNoti.NotificationType = notificationType
		go sendNoti.SendData(admindata.PROD)
	}

	n.Ctx.WriteString("OK")
}

func (n *NotificationController) MT() {
	// http://www.yourdomainDNurl/receive.asp?mtid=123296707&msisdn=60121234567&shortcod e=32300&telcoid=1&countryid=1&datetime=2010-06-15 10:10:10&status=DELIVRD&category=1&keyword=ABC

	var reqData = models.ReqData{}

	reqData.ShortCode = n.GetString("Shortcode")
	// 这里的电话号码很关键
	reqData.Msisdn = n.GetString("Msisdn")
	reqData.Keyword = n.GetString("Keyword")
	reqData.Mtid = n.GetString("mtid")
	reqData.Moid = n.GetString("moid")
	reqData.Datetime = n.GetString("datetime")
	reqData.Telcoid = n.GetString("Telcoid")
	reqData.Countryid = n.GetString("countryid")
	reqData.Status = n.GetString("status")
	reqData.Category = n.GetString("category")
	reqData.Keyword = n.GetString("Keyword")

	reqData.Insert()

	moT := new(mo.Mo)
	// 先根据subID 查询mo数据
	_, _ = moT.GetMoByMsisdnShortCodeAndKeywordID(reqData.Msisdn, reqData.ShortCode, strings.ToUpper(reqData.Keyword))

	notificationType := ""

	if moT.ID != 0 {
		if reqData.Category == "0" {
			notificationType, _ = moT.AddSuccessMTNum(moT.SubscriptionID, reqData.Mtid)
		} else {
			notificationType, _ = moT.AddFailedMTNum(moT.SubscriptionID, reqData.Mtid)
		}
	}

	if notificationType != "" {
		notify := new(notification.Notification)
		notify.NotificationType = notificationType
		notify.SubscriptionId = moT.SubscriptionID
		notify.TransactionID = reqData.Mtid
		notify.ServiceID = moT.ServiceID
		notify.Insert()

		nowTime, _ := util.GetNowTime()

		sendNoti := new(admindata.Notification)

		sendNoti.OfferID = moT.OfferID
		sendNoti.SubscriptionID = moT.SubscriptionID

		sendNoti.TransactionID = reqData.Mtid
		sendNoti.AffName = moT.AffName

		sendNoti.Sendtime = nowTime
		sendNoti.NotificationType = notificationType
		go sendNoti.SendData(admindata.PROD)
	}

	n.Ctx.WriteString("OK")
}
