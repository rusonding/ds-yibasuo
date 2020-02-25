package models

import "fmt"

// 配置内容接口
type ConfigBody interface{ ConfigToString() string }

// 前端配置
type ConfigFrontend struct {
	Frontend Frontends `yaml:"nginx"`
}

func (m *ConfigFrontend) ConfigToString() string { return "" }

type Frontends struct {
	EscProxyPort int `yaml:"esc_proxy_port"`
}

func (m *Frontends) ConfigToString() string { return "" }

// 后端配置
type ConfigBackend struct {
	Backend Backends `yaml:"api"`
}

func (m *ConfigBackend) ConfigToString() string { return "" }

type Backends struct {
	ServerPort int `yaml:"server.port"`
}

func (m *Backends) ConfigToString() string { return "" }

// 报警配置
type AlertType int

const (
	Mail = iota
	Weixin
)

type AlertBody interface{ AlertToString() string }

type AlertMail struct {
	// TODO 后面再写
}

func (m *AlertMail) AlertToString() string {
	return fmt.Sprintf(`AlertMail:{}`)
}

type AlertWeixin struct {
}

func (m *AlertWeixin) AlertToString() string {
	return fmt.Sprintf(`AlertWeixin:{}`)
}

type ConfigAlert struct {
	AlertType AlertType `json:"alertType"`
	Alert     AlertBody `json:"alert"`
}

func (m *ConfigAlert) ConfigToString() string { return "" }

// Master配置
type ConfigMaster struct {
	Master Masters `yaml:"master"`
}

func (m *ConfigMaster) ConfigToString() string { return fmt.Sprintf(`ConfigMaster:{}`) }

type Masters struct {
	MasterReservedMemory float32 `yaml:"master.reserved.memory"`
}

func (m *Masters) ConfigToString() string { return "" }

// Worker配置
type ConfigWorker struct {
	Worker Workers `yaml:"worker"`
}

func (m *ConfigWorker) ConfigToString() string { return fmt.Sprintf(`ConfigWorker:{}`) }

type Workers struct {
	WorkerReservedMemory float32 `yaml:"worker.reserved.memory"`
}

func (m *Workers) ConfigToString() string { return "" }

// zookeeper配置
type ConfigZookeeper struct {
	Zookeeper []string `json:"zookeeper"`
}

func (m *ConfigZookeeper) ConfigToString() string { return "" }

// 资源中心配置
// TODO 后面再写

// database配置
type ConfigDatabase struct {
	DatabaseType string `json:"databaseType"`
	DatabaseName string `json:"databaseName"`
	Account      string `json:"account"`
	Password     string `json:"password"`
}

func (m *ConfigDatabase) ConfigToString() string {
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
