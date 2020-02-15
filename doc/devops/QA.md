# 可能遇到的问题

[返回首页](../../)

[返回目录](./README.md)

### 错误1

```
You are using pip version 8.1.2, however version 19.3.1 is available.
You should consider upgrading via the 'pip install --upgrade pip' command.
```

### 解决方案：

执行 sudo pip install --upgrade pip

---

### 错误2

```
NameError: name 'platform_system' is not defined
```

### 解决方案：
执行 sudo pip install --upgrade setuptools

---

### 错误3

```
fatal: [192.167.8.131]: FAILED! => changed=false 
  msg: async task did not complete within the requested time - 60s
```

### 解决方案：
安装 db 或者 nginx的时候会使用yum下载一些基础包。这里给yum下载配置了120秒的超时。如果出现该错误就是120秒内没有安装完毕所导致的，请检查你的网络，并保证yum的正常。  
一般都是修改yum源成国内源，或者光盘源。