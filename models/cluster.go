package models

import (
	"ds-yibasuo/utils/black"
	"ds-yibasuo/utils/blotdb"
	"ds-yibasuo/utils/common"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
)

// 基础信息
type BaseInfo struct {
	Version    string `json:"version"`    // ds版本
	DeployUser string `json:"deployUser"` //ds部署用户
	DeployDir  string `json:"deployDir"`  //ds部署路径
}

// 角色信息
type RoleInfo struct {
	RoleName       string      `json:"roleName"`
	RoleBody       interface{} `json:"roleBody"` // ConfigInfoPuppet
	RoleDependHost []string    `json:"roleDependHost"`
}

// 集群信息
// 启动成功 startsuccess
// 启动失败 startfailure
// 停止成功 stopsuccess
// 停止失败 stopfailure
// 部署成功 deploysuccess
// 部署失败 deployfailure
// 准备     prepare
type ClusterInfo struct {
	Id     string      `json:"id"`     // 集群id
	Name   string      `json:"name"`   // 集群名称
	Status string      `json:"status"` // 集群状态
	Hosts  []string    `json:"hosts"`  // 主机信息
	Base   *BaseInfo   `json:"base"`   // 基本信息
	Roles  []*RoleInfo `json:"roles"`  // 角色信息
	Remark string      `json:"remark"` // 备注
}

// model 层
// 创建或更新集群
func (m *ClusterInfo) CreateCluster() error {
	m.Id = common.MakeUuid(m.Name + common.Now())
	hostBody, _ := json.Marshal(m)
	return blotdb.Db.Add("cluster", black.String2Byte(m.Id), hostBody)
}

// model 层
// 删除集群
func (m *ClusterInfo) DeleteCluster() error {
	return blotdb.Db.RemoveID("cluster", black.String2Byte(m.Id))
}

// model 层
// 修改集群
func (m *ClusterInfo) UpdateCluster() error {
	hostBody, _ := json.Marshal(m)
	return blotdb.Db.Update("cluster", black.String2Byte(m.Id), hostBody)
}

// model 层
// 查询指定集群
func (m *ClusterInfo) SelectCluster() (*ClusterInfo, error) {
	res, err := blotdb.Db.SelectVal("cluster", black.String2Byte(m.Id))
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, errors.New("没有查到！")
	}

	c := ClusterInfo{}
	json.Unmarshal(black.String2Byte(res[0]), &c)
	return &c, err
}

// model 层
// 查询集群列表
type ClusterInfoResult struct {
	CurrentPage int            `json:"currentPage"`
	Total       int            `json:"total"`
	Data        []*ClusterInfo `json:"data"`
}

// 判断重名
func (m *ClusterInfo) CheckName() (bool, error) {
	res, err := blotdb.Db.SelectValues("cluster")
	if err != nil || len(res) < 1 {
		return false, errors.New("查询错误 或者 没有内容！")
	}

	for _, value := range res {
		h := ClusterInfo{}
		if err := json.Unmarshal(value, &h); err != nil {
			logs.Error(err)
			return false, err
		}
		if h.Name == m.Name {
			return true, nil
		}
	}

	return false, nil
}

// page 传入 -1  代表不分页,查询所有
func SelectClusterList(page int) (*ClusterInfoResult, error) {
	res, err := blotdb.Db.SelectValues("cluster")
	if err != nil || len(res) < 1 {
		return nil, errors.New("查询错误 或者 没有内容！")
	}

	var fuck []*ClusterInfo
	for _, value := range res {
		h := ClusterInfo{}
		err := json.Unmarshal(value, &h)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		fuck = append(fuck, &h)
	}

	var result *ClusterInfoResult
	if page == -1 {
		result = &ClusterInfoResult{
			CurrentPage: page,
			Total:       len(fuck),
			Data:        fuck,
		}
	} else {
		fucks := slidingCluster(fuck, 10)
		var fuckOff []*ClusterInfo
		if len(fucks) <= page {
			fuckOff = fucks[len(fucks)-1]
		} else {
			fuckOff = fucks[page-1]
		}

		result = &ClusterInfoResult{
			CurrentPage: page,
			Total:       len(fuck),
			Data:        fuckOff,
		}
	}

	return result, nil
}

func slidingCluster(list []*ClusterInfo, step int) (res [][]*ClusterInfo) {
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

// 执行类型
type ExecuteType int

const (
	Stop = iota
	Start
	DeployUpdate
)
