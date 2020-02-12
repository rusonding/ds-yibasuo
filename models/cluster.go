package models

type ClusterInfo struct {
  Name      string `json:"name"`  // 集群名称
  Model     string `json:"model"` // 集群模式
  Status    interface{}           // 当前状态
  HostInfos []HostInfo            // 主机信息
  BaseInfo  interface{}           // 基本信息
  RoleInfos interface{}           // 角色信息
}
