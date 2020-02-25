package controllers

import (
	"ds-yibasuo/models"
	"ds-yibasuo/utils/black"
	"ds-yibasuo/utils/common"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
)

type ClusterController struct {
	beego.Controller
}

// controller层
// 创建 或 修改集群
func (c *ClusterController) CreateUpdateCluster() {
	logs.Info("controller create update cluster")
	var req models.ClusterInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		// 基本判断
		for _, value := range req.Hosts {
			if value == "" {
				c.Data["json"] = models.Response{Code: 500, Message: "主机信息错误", Result: nil}
				c.ServeJSON()
				return
			}
		}
		if req.Base.DeployPath == "" || req.Base.DeployUser == "" || req.Base.Version == "" {
			c.Data["json"] = models.Response{Code: 500, Message: "基础信息错误", Result: nil}
			c.ServeJSON()
			return
		}
		for _, value := range req.Roles {
			if value.DependHost == "" || value.RoleName == "" {
				c.Data["json"] = models.Response{Code: 500, Message: "角色信息错误", Result: nil}
				c.ServeJSON()
				return
			}
		}
		// 开始创建或修改
		if err := req.CreateUpdateCluster(); err != nil {
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
	logs.Info("controller delete host")
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
// 查询具体集群信息
func (c *ClusterController) SelectCluster() {
	logs.Info("controller select cluster")
	var req models.ClusterInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		res, err := req.SelectCluster()
		if err != nil {
			logs.Error(err)
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("查询发生错误！%s", err), Result: nil}
		} else {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: res}
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
		c.Data["json"] = models.Response{Code: 500, Message: "查询错误", Result: nil}
	} else {
		res, err := models.SelectClusterList(page)
		if err == nil {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: res}
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
func (c *ClusterController) ExecuteCluster() {
	now := common.Now()
	logs.Info("execute cluster, %s: \n", now, black.Byte2String(c.Ctx.Input.RequestBody))
	var ansible models.DevopsInfo
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ansible); err == nil {
		ansible.ExecTime = now
		if singnal, err := ansible.GetSignal(); singnal == false && err == nil {
			c.Data["json"] = models.Response{Code: 500, Message: "上一次执行未结束，请等待"}
			return
		}

		switch ansible.ExecuteType {
		case models.Start:
			ansible.BackupLog()
			ansible.Start()
		case models.Stop:
			ansible.BackupLog()
			ansible.Stop()
		case models.DeployUpdate:
			ansible.BackupLog()
			// TODO 只能允许集群中某一个集群工作，不能同时有两个，启动前需要判断
			// TODO 将当前集群的状态设置为启动，并将当前集群的信息写入到ini配置文件中
			ansible.DeployUpdate()
		default:
			c.Data["json"] = models.Response{Code: 200, Message: "请输入正确的参数"}
			return
		}
	} else {
		c.Data["json"] = models.Response{Code: 200, Message: "请输入正确的参数"}
	}

	c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: nil}
	c.ServeJSON()
}

// controller层
// 查看上一次执行日志
// 上一次的动作可能是：执行、停止、部署/升级
func (c *ClusterController) ReadLog() {
	now := common.Now()
	logs.Info("read log, %s: \n", now, black.Byte2String(c.Ctx.Input.RequestBody))
	ansible := models.DevopsInfo{ExecTime: now}

	// 分页读日志，一页10条
	page, _ := strconv.Atoi(c.Input().Get("page"))
	log, err := ansible.ReadLog((page*10)-9, page*10)
	if err != nil {
		logs.Error(err)
	}

	// 获得日志总行数
	rows, err := ansible.GetLogRows()
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
	logs.Info("get singal, %s: \n", now, black.Byte2String(c.Ctx.Input.RequestBody))
	ansible := models.DevopsInfo{ExecTime: now}
	singnal, err := ansible.GetSignal()
	if err != nil {
		logs.Error(err)
	}
	c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: singnal}
	c.ServeJSON()
}
