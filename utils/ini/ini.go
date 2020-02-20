package ini

import (
	"ds-yibasuo/utils"
	"fmt"
	"io/ioutil"
	"regexp"
	"runtime"
	"strings"
)

// 因为ansible的ini特殊性。
// 在go中不能使用常用的解析包，比如go-ini，viper都不行。
// ansible 的 inventory.ini 和 hosts.ini 都不能算是正常语法的ini配置文件。
// 只能采取正则方式读，拼字符串方式写，这种菜鸡行为。
// 第一版先这样把，后面迭代改。

const (
	SERVERS_REGEX   = "(?m)^\\[servers\\]([\\s\\S]+)^\\[db_servers\\]"
	DB_REGEX        = "(?m)^\\[db_servers\\]([\\s\\S]+)^\\[zookeeper_servers\\]"
	ZOOKEEPER_REGEX = "(?m)^\\[zookeeper_servers\\]([\\s\\S]+)^\\[master_servers\\]"
	MASTER_REGEX    = "(?m)^\\[master_servers\\]([\\s\\S]+)^\\[worker_servers\\]"
	WORKER_REGEX    = "(?m)^\\[worker_servers\\]([\\s\\S]+)^\\[api_servers\\]"
	API_REGEX       = "(?m)^\\[api_servers\\]([\\s\\S]+)^\\[alert_servers\\]"
	ALERT_REGEX     = "(?m)^\\[alert_servers\\]([\\s\\S]+)^\\[nginx_servers\\]"
	NGINX_REGEX     = "(?m)^\\[nginx_servers\\]([\\s\\S]+)^\\[all:vars\\]"

	VERSION_REGEX  = "dolphinscheduler_version.*"
	DIR_REGEX      = "deploy_dir.*"
	USER_REGEX     = "ansible_user.*"
	TYPE_REGEX     = "db_type.*"
	NAME_REGEX     = "db_name.*"
	USERNANE_REGEX = "db_username.*"
	PASSWORD_REGEX = "db_password.*"

	HOSTS_REGEX      = "(?m)^\\[servers\\]([\\s\\S]+)^\\[all:vars\\]"
	HOSTS_USER_REGEX = "username.*"

	UNWANTED_REGEX = "(?m)\\[.*\\]"
	IP_REGEX       = "(?m)\\d{0,3}\\.\\d{0,3}\\.\\d{0,3}\\.\\d{0,3}"
)

var (
	pathInventory = ""
	pathHosts     = ""
)

func init() {
	if strings.Contains(strings.ToLower(runtime.GOOS), "windows") {
		pathInventory = "..\\..\\devops\\inventory.ini"
		pathHosts = "..\\..\\devops\\hosts.ini"
	} else {
		pathInventory = "./../devops/inventory.ini"
		pathHosts = "./../devops/hosts.ini"
	}
}

type IniInventory struct {
	Servers                 []string
	DbServers               []string
	ZookeeperServers        []string
	MasterServers           []string
	WorkerServers           []string
	ApiServers              []string
	AlertServers            []string
	NginxServers            []string
	DolphinschedulerVersion string
	DeployDir               string
	AnsibleUser             string
	DbType                  string
	DbName                  string
	DbUsername              string
	DbPassword              string
}

// 读ansible ini配置文件
func (i *IniInventory) ReadInventory() error {
	dataByte, err := ioutil.ReadFile(pathInventory)
	if err != nil {
		return err
	}
	dataStr := utils.Byte2String(dataByte)

	servers := regexp.MustCompile(SERVERS_REGEX).FindString(dataStr)
	db := regexp.MustCompile(DB_REGEX).FindString(dataStr)
	zookeeper := regexp.MustCompile(ZOOKEEPER_REGEX).FindString(dataStr)
	master := regexp.MustCompile(MASTER_REGEX).FindString(dataStr)
	worker := regexp.MustCompile(WORKER_REGEX).FindString(dataStr)
	api := regexp.MustCompile(API_REGEX).FindString(dataStr)
	alert := regexp.MustCompile(ALERT_REGEX).FindString(dataStr)
	nginx := regexp.MustCompile(NGINX_REGEX).FindString(dataStr)

	i.Servers = parserList(servers)
	i.DbServers = parserList(db)
	i.ZookeeperServers = parserList(zookeeper)
	i.MasterServers = parserList(master)
	i.WorkerServers = parserList(worker)
	i.ApiServers = parserList(api)
	i.AlertServers = parserList(alert)
	i.NginxServers = parserList(nginx)

	version := regexp.MustCompile(VERSION_REGEX).FindString(dataStr)
	dir := regexp.MustCompile(DIR_REGEX).FindString(dataStr)
	user := regexp.MustCompile(USER_REGEX).FindString(dataStr)
	typ := regexp.MustCompile(TYPE_REGEX).FindString(dataStr)
	name := regexp.MustCompile(NAME_REGEX).FindString(dataStr)
	username := regexp.MustCompile(USERNANE_REGEX).FindString(dataStr)
	password := regexp.MustCompile(PASSWORD_REGEX).FindString(dataStr)

	i.DolphinschedulerVersion = parserString(version)
	i.DeployDir = parserString(dir)
	i.AnsibleUser = parserString(user)
	i.DbType = parserString(typ)
	i.DbName = parserString(name)
	i.DbUsername = parserString(username)
	i.DbPassword = parserString(password)

	return nil
}

// 写ansible ini配置文件
func (i IniInventory) WriteInventory() error {
	data := ""
	data += "[servers]\n"
	for k, v := range i.Servers {
		data += fmt.Sprintf(`%s host_ip=%s host_name=easy%d`, v, v, k+1) + "\n"
	}
	data += "\n"

	data += "[db_servers]\n"
	for _, value := range i.DbServers {
		data += value + "\n"
	}
	data += "\n"

	data += "[zookeeper_servers]\n"
	for k, v := range i.ZookeeperServers {
		data += fmt.Sprintf(`%s myid=%d`, v, k+1) + "\n"
	}
	data += "\n"

	data += "[master_servers]\n"
	for _, value := range i.MasterServers {
		data += value + "\n"
	}
	data += "\n"

	data += "[worker_servers]\n"
	for _, value := range i.WorkerServers {
		data += value + "\n"
	}
	data += "\n"

	data += "[api_servers]\n"
	for _, value := range i.ApiServers {
		data += value + "\n"
	}
	data += "\n"

	data += "[alert_servers]\n"
	for _, value := range i.AlertServers {
		data += value + "\n"
	}
	data += "\n"

	data += "[nginx_servers]\n"
	for _, value := range i.NginxServers {
		data += value + "\n"
	}
	data += "\n"

	data += "[all:vars]"

	data += fmt.Sprintf(`
dolphinscheduler_version = %s
deploy_dir = %s
ansible_user = %s
db_type = %s
db_name = %s
db_username = %s
db_password = %s
`, i.DolphinschedulerVersion, i.DeployDir, i.AnsibleUser,
		i.DbType, i.DbName, i.DbUsername, i.DbPassword)

	err := ioutil.WriteFile(pathInventory, utils.String2Byte(data), 0755)
	if err != nil {
		return err
	}
	return nil
}

type IniHosts struct {
	Servers     []string
	AnsibleUser string
}

// 读hosts ini配置文件
func (i *IniHosts) ReadHosts() error {
	dataByte, err := ioutil.ReadFile(pathInventory)
	if err != nil {
		return err
	}
	dataStr := utils.Byte2String(dataByte)

	hosts := regexp.MustCompile(HOSTS_REGEX).FindString(dataStr)
	user := regexp.MustCompile(HOSTS_USER_REGEX).FindString(dataStr)

	i.Servers = parserList(hosts)
	i.AnsibleUser = parserString(user)

	return nil
}

// 写hosts ini配置文件
func (i *IniHosts) WriteHosts() error {
	data := ""
	data += "[servers]\n"
	for _, v := range i.Servers {
		data += v + "\n"
	}
	data += "\n"
	data += "[all:vars]\n"
	data += "username = " + i.AnsibleUser

	err := ioutil.WriteFile(pathHosts, utils.String2Byte(data), 0755)
	if err != nil {
		return err
	}
	return nil
}

// 解析ini 数组私有方法
func parserList(in string) (out []string) {
	unwantedRegex, _ := regexp.Compile(UNWANTED_REGEX)
	filterHeadLast := utils.Byte2String(unwantedRegex.ReplaceAll([]byte(in), []byte("")))
	filterSpace := strings.TrimSpace(filterHeadLast)
	split := strings.Split(filterSpace, "\n")
	for _, value := range split {
		out = append(out, regexp.MustCompile(IP_REGEX).FindString(value))
	}
	return
}

// 解析ini 字符串私有方法
func parserString(in string) (out string) {
	split := strings.Split(in, "=")
	out = strings.TrimSpace(split[1])
	return
}
