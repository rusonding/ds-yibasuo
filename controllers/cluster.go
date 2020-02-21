package controllers

import (
	. "ds-yibasuo/models"
	"ds-yibasuo/utils"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
)

type ClusterController struct {
	beego.Controller
}

// 创建/更新集群
func (c *ClusterController) CreateUpdateCluster() {

}

// 删除集群
func (c *ClusterController) DeleteCluster() {

}

// 查询集群列表
func (c *ClusterController) SelectClusterList() {

}

// 执行集群
func (c *ClusterController) ExecuteCluster() {
	now := utils.Now()
	logs.Info("execute cluster, %s: \n", now, utils.Byte2String(c.Ctx.Input.RequestBody))
	var ansible DevopsInfo
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ansible); err == nil {
		ansible.ExecTime = now
		if singnal, err := ansible.GetSignal(); singnal == false && err == nil {
			c.Data["json"] = Response{Code: 500, Message: "上一次执行未结束，请等待"}
			return
		}

		switch ansible.ExecuteType {
		case Start:
			ansible.BackupLog(Start)
			ansible.Start()
		case Stop:
			ansible.BackupLog(Stop)
			ansible.Stop()
		case DeployUpdate:
			ansible.BackupLog(DeployUpdate)
			ansible.DeployUpdate()
		default:
			c.Data["json"] = Response{Code: 200, Message: "请输入正确的参数"}
			return
		}
	} else {
		c.Data["json"] = Response{Code: 200, Message: "请输入正确的参数"}
	}

	c.Data["json"] = Response{Code: 200, Message: "ok", Result: nil}
	c.ServeJSON()
}

// 查看上一次执行日志
func (c *ClusterController) ReadLog() {
	now := utils.Now()
	logs.Info("read log, %s: \n", now, utils.Byte2String(c.Ctx.Input.RequestBody))
	ansible := DevopsInfo{ExecTime: now}

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
	res := DevopsLogResult{
		CurrentPage: page,
		Rows:        rows,
		Data:        log,
	}

	c.Data["json"] = Response{Code: 200, Message: "ok", Result: res}
	c.ServeJSON()
}

// 执行结果信号
func (c *ClusterController) ExecuteResultSignal() {
	now := utils.Now()
	logs.Info("get singal, %s: \n", now, utils.Byte2String(c.Ctx.Input.RequestBody))
	ansible := DevopsInfo{ExecTime: now}
	singnal, err := ansible.GetSignal()
	if err != nil {
		logs.Error(err)
	}
	c.Data["json"] = Response{Code: 200, Message: "ok", Result: singnal}
	c.ServeJSON()
}
