package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type ReqData struct {
	ID        int    `orm:"pk;auto;column(id)"`
	Text      string `json:"text"`
	Time      string `json:"time"`
	Msgid     string `json:"msgid"`
	Moid      string `json:"moid"`
	ShortCode string `json:"short_code"`
	Telcoid   string `json:"telcoid"`
	Refid     string `json:"refid"`
	Mtid      string `json:"mtid"`
	Datetime  string `json:"datetime"`
	Msisdn    string `json:"msisdn"`
	Countryid string `json:"countryid"`
	Status    string `json:"status"`
	Category  string `json:"category"`
	Keyword   string `json:"keyword"`
}

func (reqData *ReqData) Insert() {
	o := orm.NewOrm()
	_, err := o.Insert(reqData)
	if err != nil {
		logs.Error("ReqData  Insert  ERROR: ", err.Error())
	}
}
