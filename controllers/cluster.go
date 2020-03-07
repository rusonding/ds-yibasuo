package controllers

import (
	"ds-yibasuo/models"
	"ds-yibasuo/utils/black"
	"ds-yibasuo/utils/common"
	"ds-yibasuo/utils/ini"
	"ds-yibasuo/utils/yml"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/mitchellh/mapstructure"
	"strconv"
)

type ClusterController struct {
	beego.Controller
}

// controller层
// 创建集群
func (c *ClusterController) CreateCluster() {
	logs.Info("controller create update cluster")
	var req models.ClusterInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		// 判断重名
		if ok, err := req.CheckName(); err == nil {
			if ok {
				c.Data["json"] = models.Response{Code: 500, Message: "名称冲突", Result: nil}
				c.ServeJSON()
				return
			}
		}

		// 基本判断
		for _, value := range req.Hosts {
			if value == "" {
				c.Data["json"] = models.Response{Code: 500, Message: "主机信息错误", Result: nil}
				c.ServeJSON()
				return
			}
		}
		if req.Base.DeployDir == "" || req.Base.DeployUser == "" || req.Base.Version == "" {
			c.Data["json"] = models.Response{Code: 500, Message: "基础信息错误", Result: nil}
			c.ServeJSON()
			return
		}
		// TODO 集群信息请求体的基本验证
		// 开始创建或修改
		req.Status = "prepare"
		if err := req.CreateCluster(); err != nil {
			c.Data["json"] = models.Response{Code: 500, Message: err.Error(), Result: nil}
		} else {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: nil}
		}
	} else {
		logs.Error(err)
		c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
	}

	c.ServeJSON()
}

// controller层
// 删除集群
func (c *ClusterController) DeleteCluster() {
	logs.Info("controller delete cluster")
	var req models.ClusterInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		err := req.DeleteCluster()
		if err != nil {
			logs.Error(err)
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("删除发生错误！%s", err), Result: nil}
		} else {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: nil}
		}
	} else {
		c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
	}

	c.ServeJSON()
}

// controller层
// 修改集群
// 取消状态是从修改集群中传进来的
func (c *ClusterController) UpdateCluster() {
	logs.Info("controller update cluster")
	var req models.ClusterInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		// 判断重名
		if ok, err := req.CheckName(); err == nil {
			if ok {
				c.Data["json"] = models.Response{Code: 500, Message: "名称冲突", Result: nil}
				c.ServeJSON()
				return
			}
		}

		if req.Id == "" {
			c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
			c.ServeJSON()
			return
		}
		if err := req.UpdateCluster(); err != nil {
			logs.Error(err)
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("修改发生错误！%s", err), Result: nil}
		} else {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: nil}
		}
	} else {
		c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
	}

	c.ServeJSON()
}

// controller层
// 查询具体集群信息
func (c *ClusterController) SelectCluster() {
	logs.Info("controller select cluster")
	var req models.ClusterInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		res, err := req.SelectCluster()
		if err == nil {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: res}
		} else if err.Error() == "null" {
			c.Data["json"] = models.Response{Code: 200, Message: "没有查到", Result: nil}
		} else {
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("查询发生错误！%s", err), Result: nil}
		}
	} else {
		c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
	}

	c.ServeJSON()
}

// controller层
// 查询集群列表
func (c *ClusterController) SelectClusterList() {
	logs.Info("controller select cluster list")
	page, err := strconv.Atoi(c.Input().Get("page"))
	if err != nil {
		c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
	} else {
		res, err := models.SelectClusterList(page)
		if err == nil {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: res}
		} else if err.Error() == "null" {
			c.Data["json"] = models.Response{Code: 200, Message: "没有查到", Result: nil}
		} else {
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("查询发生错误！%s", err), Result: nil}
		}
	}

	c.ServeJSON()
}

// controller层
// 执行集群
// 0 停止
// 1 启动
// 2 部署/升级
// TODO 代码太多 隐患
func (c *ClusterController) ExecuteCluster() {
	now := common.Now()
	logs.Info("controller execute cluster: \n", black.Byte2String(c.Ctx.Input.RequestBody))
	var dev models.DevopsInfo
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dev); err == nil {
		dev.ExecTime = now
		// 检查上一次执行是否成功结束
		if singnal, err := dev.GetSignal(); err == nil {
			if singnal.Over == false {
				c.Data["json"] = models.Response{Code: 500, Message: "上一次执行未结束，请等待！！！"}
				c.ServeJSON()
				return
			}
		}
		// 查询是否有已经启动的集群
		// TODO 这样的跳出方式,我有点无奈
		if res, err := models.SelectClusterList(-1); err == nil && dev.ExecuteType != models.Stop {
			for _, value := range res.Data {
				if value.Status == "startsuccess" {
					c.Data["json"] = models.Response{Code: 500, Message: "同一时间只能启动一个集群,请停止再操作其他集群!!"}
					c.ServeJSON()
					return
				}
			}
		}

		// 读取集群信息
		cluster := models.ClusterInfo{}
		cluster.Id = dev.ClusterId
		clusterInfo, err := cluster.SelectCluster()
		if err != nil {
			logs.Error(err)
			c.Data["json"] = models.Response{Code: 500, Message: "查询主机出错", Result: nil}
			c.ServeJSON()
			return
		}
		switch dev.ExecuteType {
		case models.Start:
			dev.BackupLog(models.Start)
			dev.Start()
			go dev.UpdateStatus(models.Start, clusterInfo)
		case models.Stop:
			dev.BackupLog(models.Stop)
			dev.Stop()
			go dev.UpdateStatus(models.Stop, clusterInfo)
		case models.DeployUpdate:
			dev.BackupLog(models.DeployUpdate)
			// 写入ini配置信息
			i := ini.IniInventory{}
			i.Servers = clusterInfo.Hosts
			for _, role := range clusterInfo.Roles {
				switch role.RoleName {
				case "database":
					var db models.ConfigDatabase
					if err := mapstructure.Decode(role.RoleBody.(map[string]interface{}), &db); err != nil {
						logs.Error("db map to struct err: ", err)
					}
					i.DbServers = role.RoleDependHost
					i.DbType = db.DatabaseType
					i.DbName = db.DatabaseName
					i.DbUsername = db.Account
					i.DbPassword = db.Password
				case "zookeeper":
					i.ZookeeperServers = role.RoleDependHost
				case "master":
					var m models.ConfigMaster
					if err := mapstructure.Decode(role.RoleBody.(map[string]interface{}), &m); err != nil {
						logs.Error("master map to struct err: ", err)
					}
					err = yml.WriteYml(models.MasterYml, &m)
					if err != nil {
						logs.Error("master yml err: ", err)
					}
					i.MasterServers = role.RoleDependHost
				case "worker":
					var m models.ConfigWorker
					if err := mapstructure.Decode(role.RoleBody.(map[string]interface{}), &m); err != nil {
						logs.Error("worker map to struct err: ", err)
					}
					err = yml.WriteYml(models.WorkerYml, &m)
					if err != nil {
						logs.Error("worker yml err: ", err)
					}
					i.WorkerServers = role.RoleDependHost
				case "backend":
					var b models.ConfigBackend
					if err := mapstructure.Decode(role.RoleBody.(map[string]interface{}), &b); err != nil {
						logs.Error("backend map to struct err: ", err)
					}
					err := yml.WriteYml(models.BackendYml, &b)
					if err != nil {
						logs.Error("api yml err: ", err)
					}
					i.ApiServers = role.RoleDependHost
				case "alert":
					var a models.ConfigAlert
					if err := mapstructure.Decode(role.RoleBody.(map[string]interface{}), &a); err != nil {
						logs.Error("alert map to struct err: ", err)
					}
					if err := yml.WriteYml(models.AlertYml, &a); err != nil {
						logs.Error("alert yml err: ", err)
					}
					i.AlertServers = role.RoleDependHost
				case "frontend":
					var f models.ConfigFrontend
					if err := mapstructure.Decode(role.RoleBody.(map[string]interface{}), &f); err != nil {
						logs.Error("frontend map to struct err: ", err)
					}
					if err := yml.WriteYml(models.FrontendYml, &f); err != nil {
						logs.Error("frontend yml err: ", err)
					}
					i.NginxServers = role.RoleDependHost
				case "hadoop":
					var h models.ConfigHadoop
					if err := mapstructure.Decode(role.RoleBody.(map[string]interface{}), &h); err != nil {
						logs.Error("hadoop map to struct err: ", err)
					}
					if err := yml.WriteYml(models.HadoopYml, &h); err != nil {
						logs.Error("hadoop yml err: ", err)
					}
				case "common":
					var c models.ConfigCommon
					if err := mapstructure.Decode(role.RoleBody.(map[string]interface{}), &c); err != nil {
						logs.Error("common map to struct err: ", err)
					}
					if err := yml.WriteYml(models.CommonYml, &c); err != nil {
						logs.Error("common yml err: ", err)
					}
				}
			}
			i.DolphinschedulerVersion = clusterInfo.Base.Version
			i.DeployDir = clusterInfo.Base.DeployDir
			i.AnsibleUser = clusterInfo.Base.DeployUser
			i.WriteInventory()
			// 开始ansible 部署
			dev.DeployUpdate()
			// 异步去修改集群的状态
			go dev.UpdateStatus(models.DeployUpdate, clusterInfo)
		default:
			c.Data["json"] = models.Response{Code: 200, Message: "请输入正确的参数"}
			c.ServeJSON()
			return
		}
	} else {
		c.Data["json"] = models.Response{Code: 200, Message: "参数错误"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: nil}
	c.ServeJSON()
}

// controller层
// 查看上一次执行日志
// 上一次的动作可能是：执行、停止、部署/升级
func (c *ClusterController) ReadLog() {
	now := common.Now()
	logs.Info("controller read log: \n", now, black.Byte2String(c.Ctx.Input.RequestBody))
	dev := models.DevopsInfo{ExecTime: now}
	// 分页读日志，一页10条
	page, _ := strconv.Atoi(c.Input().Get("page"))
	log, err := dev.ReadLog((page*10)-9, page*10)
	if err != nil {
		logs.Error(err)
	}
	// 获得日志总行数
	rows, err := dev.GetLogRows()
	if err != nil {
		logs.Error(err)
	}
	// 组装响应
	res := models.DevopsLogResult{
		CurrentPage: page,
		Rows:        rows,
		Data:        log,
	}

	c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: res}
	c.ServeJSON()
}

// controller层
// 执行结果信号
// 因为ansible执行需要一定时间，接口层面将判断是否成功改为异步
// 只要成功执行是会立刻返回200
// 前端需要每隔一段时间来调用该接口来判断执行是否结束
func (c *ClusterController) ExecuteResultSignal() {
	now := common.Now()
	logs.Info("controller get signal: \n", black.Byte2String(c.Ctx.Input.RequestBody))
	dev := models.DevopsInfo{ExecTime: now}
	singnal, err := dev.GetSignal()
	if err != nil {
		logs.Error(err)
	}
	c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: singnal}
	c.ServeJSON()
}
