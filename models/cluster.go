package models

// 基础信息
type BaseInfo struct {
	Version    string `json:"version"`    // ds版本
	DeployUser string `json:"deployUser"` //ds部署用户
	DeployPath string `json:"deployPath"` //ds部署路径
}

// 角色信息
type RoleInfo struct {
	RoleName     string     `json:"name"`
	RoleBody     ConfigBody `json:"roleBody"`
	DependHostId int        `json:"dependHostId"`
}

// 集群状态
type ClusterType int

const (
	Stopped = iota //停止
	Started        //启动
)

// 集群信息
type ClusterInfo struct {
	Id     int         `json:"id"`     // 集群id
	Name   string      `json:"name"`   // 集群名称
	Status ClusterType `json:"status"` // 当前状态
	Hosts  []int       `json:"hosts"`  // 主机信息
	Base   BaseInfo    `json:"base"`   // 基本信息
	Roles  []RoleInfo  `json:"roles"`  // 角色信息
}

// 执行类型
type ExecuteType int

const (
	Stop = iota
	Start
	DeployUpdate
)

// 执行结果
type ExecuteResultType int

const (
	Success = iota
	Failed
)

// 执行请求
type ExecuteInfo struct {
	ExecuteType   ExecuteType       `json:"executeType"`
	ClusterId     int               `json:"clusterId"`
	ExecuteResult ExecuteResultType `json:"executeResult"`
}
