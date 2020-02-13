package models

import (
	. "ds-yibasuo/utils"
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

type HostInfo struct {
	Id     uint64 `json:"id"`
	Ip     string `json:"ip"`
	Name   string `json:"name"`
	Port   string `json:"port"`
	Root   string `json:"root"`
	Remark string `json:"remark"`
}

func (m *HostInfo) HostInsert() (uint64, error) {
	hostBody, _ := json.Marshal(m)
	id, err := Db.Add("host", hostBody)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *HostInfo) HostDelete() {

}

func (m *HostInfo) HostUpdate() {

}

func (m *HostInfo) HostSelect() {

}

func QueryHostList() (string, error) {
	res, err := Db.SelectAll2Map("host")
	if err != nil {
		return "", err
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		logs.Error(err)
	}
	return string(jsonRes), nil
}
