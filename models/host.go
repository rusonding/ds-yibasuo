package models

import (
	"ds-yibasuo/utils/black"
	"ds-yibasuo/utils/blotdb"
	"ds-yibasuo/utils/common"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
)

type HostInfo struct {
	Id     string `json:"id"`
	Ip     string `json:"ip"`
	Name   string `json:"name"`
	Port   int    `json:"port"`
	Root   string `json:"root"`
	Remark string `json:"remark"`
}

// model层
// 创建主机
func (m *HostInfo) CreateHost() error {
	m.Id = common.MakeUuid(m.Name + m.Ip + common.Now())
	hostBody, _ := json.Marshal(m)
	return blotdb.Db.Add("host", black.String2Byte(m.Id), hostBody)
}

// model层
// 根据id删除指定host
// 删除前将id从uint64 转成了 byte
func (m *HostInfo) DeleteHost() error {
	return blotdb.Db.RemoveID("host", black.String2Byte(m.Id))
}

// model层
// 将实体里的内容覆盖插入原id内，达到更新的效果
func (m *HostInfo) UpdateHost() error {
	hostBody, _ := json.Marshal(m)
	return blotdb.Db.Update("host", black.String2Byte(m.Id), hostBody)
}

// model层
// 根据id来查询指定host，并返回string
// 因为查出来是一个string list 这里默认取第一个，可能是个隐患
func (m *HostInfo) SelectHost() (*HostInfo, error) {
	res, err := blotdb.Db.SelectVal("host", black.String2Byte(m.Id))
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, errors.New("null")
	}

	h := HostInfo{}
	json.Unmarshal(black.String2Byte(res[0]), &h)
	return &h, err
}

func (m *HostInfo) CheckName() (bool, error) {
	res, err := blotdb.Db.SelectValues("host")
	if err != nil || len(res) < 1 {
		return false, errors.New("查询错误 或者 没有内容！")
	}

	for _, value := range res {
		h := HostInfo{}
		if err := json.Unmarshal(value, &h); err != nil {
			logs.Error(err)
			return false, err
		} else {
			if h.Name == m.Name {
				return true, nil
			}
		}
	}

	return false, nil
}

// model层
// 查询host列表
// 得到map结果，再序列化成string返回
type HostInfoResult struct {
	CurrentPage int         `json:"currentPage"`
	Total       int         `json:"total"`
	Data        []*HostInfo `json:"data"`
}

func SelectHostList(page int) (*HostInfoResult, error) {
	res, err := blotdb.Db.SelectValues("host")
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("null")
	}

	var fuck []*HostInfo
	for _, value := range res {
		h := HostInfo{}
		err := json.Unmarshal(value, &h)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		fuck = append(fuck, &h)
	}

	fucks := slidingHost(fuck, 10)

	var fuckOff []*HostInfo
	if len(fucks) <= page {
		fuckOff = fucks[len(fucks)-1]
	} else {
		fuckOff = fucks[page-1]
	}

	result := &HostInfoResult{
		CurrentPage: page,
		Total:       len(fuck),
		Data:        fuckOff,
	}

	return result, nil
}

func slidingHost(list []*HostInfo, step int) (res [][]*HostInfo) {
	start, end := 0, 0
	for {
		if len(list) <= 0 {
			break
		}
		if (start + step) > len(list) {
			end = len(list)
		} else {
			end += step
		}
		res = append(res, list[start:end])
		start += step
		if start > len(list) {
			break
		}
	}
	return
}
