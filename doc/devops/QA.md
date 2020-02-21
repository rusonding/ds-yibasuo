# 可能遇到的问题

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
  msg: async task did not complete within the requested time
```

### 解决方案：

安装 db 或者 nginx的时候会使用yum下载一些基础包。这里给yum下载配置了90秒的超时。如果出现该错误就是120秒内没有安装完毕所导致的，请检查你的网络，并保证yum的正常。  
一般都是修改yum源成国内源，或者光盘源。


[返回目录](./README.md)

### 错误4

```
TASK [dancer : exists mysql client] ****************************************************************************************************************************************************************************
failed: [localhost] (item={u'url': u'https://repo1.maven.org/maven2/mysql/mysql-connector-java/5.1.48/mysql-connector-java-5.1.48.jar', u'version': u'5.1.48', u'name': u'mysql-connector-java'}) => changed=true 
  cmd: ls /home/easy/ds-yibasuo-web/devops/resources/pkg | grep mysql-connector-java-5.1.48.jar
  delta: '0:00:00.022663'
  end: '2020-02-21 16:15:18.230829'
  item:
    name: mysql-connector-java
    url: https://repo1.maven.org/maven2/mysql/mysql-connector-java/5.1.48/mysql-connector-java-5.1.48.jar
    version: 5.1.48
  msg: non-zero return code
  rc: 1
  start: '2020-02-21 16:15:18.208166'
  stderr: ''
  stderr_lines: []
  stdout: ''
  stdout_lines: <omitted>
...ignoring
```
### 解决方案：

注意最后有个 ignoring 的字样，这个错误不用管，是作者没处理好【错误类型】输出。不影响使用，用户可以认为这是正常的。