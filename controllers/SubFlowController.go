package controllers

import (
	"fmt"
	"github.com/angui001/MKMalai/libs"
	"github.com/angui001/MKMalai/models"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/guaidashu/go_helper/crypto_tool"
	"strconv"
	"strings"
	"time"
)

type SubFlowController struct {
	BaseController

	TrackID int

	track *models.AffTrack

	serviceConf models.ServiceInfo
}

func (c *SubFlowController) Prepare() {
	trackID := c.GetString("tid")
	trackIDInt, _ := strconv.Atoi(trackID)
	if trackIDInt == 0 {
		logs.Error("SubFlowController Prepare 将trackID转为INT类型失败", trackID)
		c.RedirectURL("https://google.com")
	}

	c.TrackID = trackIDInt

	c.track = c.getTrackData(trackIDInt)

	// 配置信息
	c.serviceConf = c.getServiceConfig(c.track.Track.ServiceID)
	fmt.Println(c.serviceConf)
}

// /sub/req?tid={track_id}
func (c *SubFlowController) SubReq() {
	// 首先 获取 auth token
	authToken := c.GetAuthToken()

	fmt.Println(authToken)

	// 然后进行 跳转到 AOC页面

	fmt.Println("ok")
}

func (s *SubFlowController) GetAuthToken() (token string) {
	// url: http://mis.etracker.cc/MYPSMSToken/GetToken  GET方式
	// example: http://mis.etracker.cc/MYPSMSToken/GetToken?AccessToken=aa15856d7b7fbb796debe4b6fb3aa715&Keyw ord=ring&ShortCode=4541889&DateTime=20180914155512
	var (
		err         error
		accessToken string
		userName    string
		keyWord     string
		shortCode   string
		dateTime    string
		password    string
		result      []byte
	)

	defer func() {
		// token = ""
		if token == "" {
			fmt.Println("出错了")
		}
	}()

	requestUrl := "http://mis.etracker.cc/MYPSMSToken/GetToken?AccessToken=%v&Keyword=%v&ShortCode=%v&DateTime=%v"

	// 生成accessToken
	password = "seek0804"
	userName = "foxseek"

	var cstZone = time.FixedZone("CST", 8*3600)
	now := time.Now().In(cstZone)

	year, mon, day := now.Date()
	hour, min, sec := now.Clock()

	dateTime = fmt.Sprintf("%d%02d%02d%02d%02d%02d", year, mon, day, hour, min, sec)

	password = crypto_tool.Md5(password)
	keyWord = s.serviceConf.Keyword
	shortCode = s.serviceConf.ShortCode

	accessToken = "%v%v%v%v%v"
	accessToken = fmt.Sprintf(accessToken, strings.ToUpper(userName), strings.ToUpper(keyWord), strings.ToUpper(shortCode), dateTime, strings.ToUpper(password))

	fmt.Println("accessToken: ", accessToken)

	accessToken = crypto_tool.MD5(accessToken)

	requestUrl = fmt.Sprintf(requestUrl, accessToken, keyWord, shortCode, dateTime)

	fmt.Println("request url is: ", requestUrl)

	if result, err = httplib.Get(requestUrl).Bytes(); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
	}

	fmt.Println(string(result))

	resultArr := strings.Split(string(result), ",")

	// 如果出错了，数组就只有一个数据
	if len(resultArr) == 1 {
		token = ""
	} else {
		token = resultArr[0]
	}

	fmt.Println(token)

	return
}

// 拿到authToken之后 去请求 url，就会自动跳转到aoc页面了
func (s *SubFlowController) SendMoMessageViaWAP() {
	// example:
	// http://mis.etracker.cc/MYWAP/WAPMORequest.aspx?Telcoid=3&Shortcode=32300&Keyword=ring&Refid=12345 6789abc&AuthToken=10100000000
	sendUrl := "http://mis.etracker.cc/MYWAP/WAPMORequest.aspx?Telcoid=%v&Shortcode=%v&Keyword=%v&Refid=%v&AuthToken=%v"

	// 我们只有digi运营商，所以telcoid直接就是3
	telcoid := 3
	shortcode := s.serviceConf.ShortCode
	keyword := s.serviceConf.Keyword
	// refid 是track_id
	refid := s.GetString("tid")
	authToken := s.GetAuthToken()

	sendUrl = fmt.Sprintf(sendUrl, telcoid, shortcode, keyword, refid, authToken)

	fmt.Println("request to aoc url is: ", sendUrl)

	// 访问链接后自动 跳转 到 aoc
	s.RedirectURL(sendUrl)

}
