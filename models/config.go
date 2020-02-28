package models

import (
	"ds-yibasuo/utils/black"
	"ds-yibasuo/utils/blotdb"
	"ds-yibasuo/utils/common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
)

var (
	ConfigList  = []string{"frontend", "backend", "alert", "master", "worker", "resources", "database", "zookeeper"}
	FrontendYml = "./devops/conf/nginx.yml"
	BackendYml  = "./devops/conf/api.yml"
	AlertYml    = "./devops/conf/alert.yml"
	MasterYml   = "./devops/conf/master.yml"
	WorkerYml   = "./devops/conf/worker.yml"
	HadoopYml   = "./devops/conf/hadoop.yml"
	CommonYml   = "./devops/conf/common.yml"
)

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
//type AlertType int
//
//const (
//	Mail = iota
//	Weixin
//)

type AlertBody interface{ AlertToString() string }

type AlertMail struct {
	AlertType    string `yaml:"alert.type"`
	MailProtocol string `yaml:"mail.protocol"`
	MailSender   string `yaml:"mail.sender"`
	MailUser     string `yaml:"mail.user"`
	MailPasswd   string `yaml:"mail.passwd"`
	MailSmtpHost string `yaml:"mail.server.host"`
	MailSmtpPort int    `yaml:"mail.server.port"`
	MailSsl      bool   `yaml:"mail.smtp.ssl.enable"`
	MailTls      bool   `yaml:"mail.smtp.starttls.enable"`
}

func (m *AlertMail) AlertToString() string {
	return fmt.Sprintf(`AlertMail:{}`)
}

type AlertWeixin struct {
	//TODO 后面在写
}

func (m *AlertWeixin) AlertToString() string {
	return fmt.Sprintf(`AlertWeixin:{}`)
}

//序列化的时候为什么不认识这个类型呢AlertBody
//错误 * 'Alert' expected type 'models.AlertBody', got 'map[string]interface {}'
type ConfigAlert struct {
	Alert interface{} `yaml:"alert"`
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
// TODO 第一版不用配zk参数
type ConfigZookeeper struct {
	Zookeeper []string `json:"zookeeper"`
}

func (m *ConfigZookeeper) ConfigToString() string { return "" }

// 资源中心配置
type Hadoops struct {
	FsDefaultFS string `yaml:"fs.defaultFS"`
}

func (m *Hadoops) ConfigToString() string { return "" }

type ConfigHadoop struct {
	Hadoop Hadoops `yaml:"hadoop"`
}

func (m *ConfigHadoop) ConfigToString() string { return "" }

type Commons struct {
	DataStore2hdfsBasepath string `yaml:"data.store2hdfs.basepath"`
}

func (m *Commons) ConfigToString() string { return "" }

type ConfigCommon struct {
	Common Commons `yaml:"common"`
}

func (m *ConfigCommon) ConfigToString() string { return "" }

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
	Database
	Hadoop
	Common
)

type ConfigInfo struct {
	Id   string      `json:"id"`   // 配置id
	Name string      `json:"name"` // 配置别名
	Typ  string      `json:"typ"`  // 配置类型
	Conf interface{} `json:"conf"` // 配置内容 //TODO interface 隐患
}

func (m *ConfigInfo) CreateConfig() error {
	m.Id = common.MakeUuid(m.Name + m.Typ)
	hostBody, _ := json.Marshal(m)
	return blotdb.Db.Add("config", black.String2Byte(m.Id), hostBody)
}

func (m *ConfigInfo) DeleteConfig() error {
	return blotdb.Db.RemoveID("config", black.String2Byte(m.Id))
}

func (m *ConfigInfo) UpdateConfig() error {
	hostBody, _ := json.Marshal(m)
	return blotdb.Db.Update("config", black.String2Byte(m.Id), hostBody)
}

func (m *ConfigInfo) SelectConfig() (*ConfigInfo, error) {
	res, err := blotdb.Db.SelectVal("config", black.String2Byte(m.Id))
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, errors.New("没有查到！")
	}

	c := ConfigInfo{}
	json.Unmarshal(black.String2Byte(res[0]), &c)
	return &c, err
}

type ConfigInfoResult struct {
	CurrentPage int           `json:"currentPage"`
	Total       int           `json:"total"`
	Data        []*ConfigInfo `json:"data"`
}

func SelectConfigList(page int, typ string) (*ConfigInfoResult, error) {
	res, err := blotdb.Db.SelectValues("config")
	if err != nil || len(res) < 1 {
		return nil, errors.New("查询错误 或者 没有内容！")
	}

	var fuck []*ConfigInfo
	for _, value := range res {
		h := ConfigInfo{}
		err := json.Unmarshal(value, &h)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		if h.Typ == typ {
			fuck = append(fuck, &h)
		}
	}
	if len(fuck) == 0 {
		return nil, errors.New("查询错误 或者 没有内容！")
	}

	fucks := slidingConfig(fuck, 10)

	var fuckOff []*ConfigInfo
	if len(fucks) <= page {
		fuckOff = fucks[len(fucks)-1]
	} else {
		fuckOff = fucks[page-1]
	}

	result := &ConfigInfoResult{
		CurrentPage: page,
		Total:       len(fuck),
		Data:        fuckOff,
	}

	return result, nil
}

func slidingConfig(list []*ConfigInfo, step int) (res [][]*ConfigInfo) {
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
