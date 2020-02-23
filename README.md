# ds-yibasuo-web

ds 一把梭 安装，部署，运维 mini web 应用

## 原型地址

[点击跳转](http://121.36.110.148:22245/start_1.html)

## 后端接口文档

[点击跳转](./doc/backend)

## 运维文档

[点击跳转](./doc/devops)

---

开发TODO

- [x] **工具**
  - [x] yml配置文件读写 （剩1个资源管理相关TODO）
  - [x] ini配置文件读写
  - [x] 其他工具（类型转换，时间获取）
  - [x] blotdb 增删该查基础包
- [ ] **登陆、权限**
  - [ ] 用户权限，只读
  - [ ] 所有接口与权限对接
- [x] **主机管理**
  - [x] 主机管理文档
  - [x] 主机管理实体设计
    - [x] 主机管理实体成员设计
    - [ ] 主机管理实体方法设计
  - [ ] 主机接口逻辑
    - [ ] 创建主机
    - [ ] 删除主机
    - [ ] 更新主机
    - [ ] 查询具体主机
    - [ ] 查询主机列表
    - [ ] 主机对接ansible
- [x] **配置管理**
  - [x] 配置管理文档（剩4个alert相关TODO）
  - [x] 配置管理实体设计
    - [x] 主机管理实体成员设计
    - [ ] 主机管理实体方法设计
  - [ ] 配置接口逻辑
    - [ ] 创建配置
    - [ ] 删除配置
    - [ ] 更新配置
    - [ ] 查询具体配置
    - [ ] 查询配置列表
- [x] **集群管理**
  - [x] 集群管理文档 （剩1个创建相关TODO）
  - [x] 集群管理实体设计
    - [x] 主机管理实体成员设计
    - [ ] 主机管理实体方法设计
  - [ ] 集群接口逻辑
    - [ ] 集群增删改查
      - [ ] 创建/更新集群
      - [ ] 删除集群
      - [ ] 查询具体集群
      - [ ] 查询集群列表
    - [x] 对接ansible操作
      - [x] 对接ansible执行
      - [x] 对接ansible日志查看
      - [x] 对接ansible接收执行信号
- [ ] **系统设置**
  - [ ] 新手上路开关
  - [ ] 修改运维密码
  - [ ] 整体配置导出
