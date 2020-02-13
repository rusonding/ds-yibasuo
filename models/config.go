package models

import "fmt"

// 配置内容接口
type ConfigBody interface {
	String()
}

// 前端配置
type ConfigFrontend struct {
	FrontendPort int `json:"frontendPort"`
}

func (m *ConfigFrontend) String() string {
	return fmt.Sprintf(`"configFrontend":{"frontendPort":%d}`, m.FrontendPort)
}

// 后端配置
type ConfigBackend struct {
	BackendPort int `json:"backendPort"` // server.port
}

func (m *ConfigBackend) string() string {
	return fmt.Sprintf(`ConfigBackend:{}`)
}

// 报警配置
type AlertType int

const (
	Mail = iota
	Weixin
)

type AlertBody interface{ String() string }

type AlertMail struct {
}

func (m *AlertMail) String() string {
	return fmt.Sprintf(`AlertMail:{}`)
}

type AlertWeixin struct {
}

func (m *AlertWeixin) String() string {
	return fmt.Sprintf(`AlertWeixin:{}`)
}

type ConfigAlert struct {
	AlertType AlertType `json:"alertType"`
	Alert     AlertBody `json:"alert"`
}

// Master配置
type ConfigMaster struct {
	MasterMem float32 `json:"masterMem"`
}

func (m *ConfigMaster) String() string {
	return fmt.Sprintf(`ConfigMaster:{}`)
}

// Worker配置
type ConfigWorker struct {
	WorkerMem float32 `json:"workerMem"`
}

func (m *ConfigWorker) String() string {
	return fmt.Sprintf(`ConfigWorker:{}`)
}

// db配置
type DatabaseType int

const (
	Mysql = iota
	PostgreSQL
)

type ConfigDatabase struct {
	DatabaseType DatabaseType `json:"databaseType"`
	DbName       string       `json:"name"`
	Account      string       `json:"account"`
	Password     string       `json:"password"`
}

func (m *ConfigDatabase) String() string {
	return fmt.Sprintf(`ConfigDB:{}`)
}

// 配置信息
type ConfigType int

const (
	Frontend = iota
	Backend
	Alert
	Master
	Worker
	Zookeeper
	Resources
	Database
)

type ConfigInfo struct {
	Id   int        `json:"id"`   // 配置id
	Name string     `json:"name"` // 配置别名
	Typ  ConfigType `json:"typ"`  // 配置类型
	Conf ConfigBody `json:"conf"` // 配置内容
}
