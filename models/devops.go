package models

import (
	"bufio"
	. "ds-yibasuo/utils/black"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	ANSIBLE_LOG = "./devops/log/ansible.log"
)

type DevopsInfo struct {
	ExecuteType ExecuteType `json:"executeType"`
	ClusterId   int         `json:"clusterId"`
	ExecTime    string      `json:"execTime"`
}

func (m *DevopsInfo) BackupLog() {
	logs.Info("ansible backup log ")
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("mv %s %s.%s", ANSIBLE_LOG, ANSIBLE_LOG, m.ExecTime))
	cmd.Start()
}

func (m *DevopsInfo) DeployUpdate() {
	cmd := exec.Command("ansible-playbook", "deploy.yml")
	cmd.Dir = "./devops"
	cmd.Start()
}

func (m *DevopsInfo) ReadLog(start, end int) ([]string, error) {
	filePath := "./devops/log/ansible.log"
	//filePath := "C:\\iceblue\\go\\src\\ds-yibasuo\\test\\ansible.log"
	var result []string

	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return nil, err
	}
	defer file.Close()

	buf := bufio.NewReader(file)

	i := 1
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if i >= start && i <= end {
			result = append(result, line)
		}

		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return nil, err
			}
		}

		if i == end {
			return result, nil
		}
		i += 1
	}
	return result, nil
}

func (m *DevopsInfo) GetLogRows() (int, error) {
	fuck := fmt.Sprintf("wc -l %s | awk '{print $1}'", ANSIBLE_LOG)

	cmd := exec.Command("/bin/bash", "-c", fuck)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	//return int(binary.BigEndian.Uint16(out)), nil
	fuckRows, _ := strconv.Atoi(strings.Replace(Byte2String(out), "\n", "", -1))
	return fuckRows, nil
}

type DevopsLogResult struct {
	CurrentPage int
	Rows        int
	Data        []string
}

func (m *DevopsInfo) GetSignal() (bool, error) {
	cmd := exec.Command("/bin/bash", "-c", "ps -ef | grep ansible-playbook | grep -v grep | wc -l")
	bytes, err := cmd.Output()
	if err != nil {
		return false, err
	}
	ansibleCount, err := strconv.Atoi(strings.Replace(Byte2String(bytes), "\n", "", -1))
	if err != nil {
		return false, err
	}
	if ansibleCount == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (m *DevopsInfo) Start() {
	cmd := exec.Command("ansible-playbook", "start.yml")
	cmd.Dir = "./devops"
	cmd.Start()
}

func (m *DevopsInfo) Stop() {
	cmd := exec.Command("ansible-playbook", "stop.yml")
	cmd.Dir = "./devops"
	cmd.Start()
}

func (m *DevopsInfo) RefreshHost(pwd string) error {
	logs.Info("change password")
	changePwd1 := fmt.Sprintf("sed -i 's/$pwd/%s/g' ./devops/hosts.ini", pwd)
	changePwd2 := fmt.Sprintf("sed -i 's/%s/$pwd/g' ./devops/hosts.ini", pwd)
	err := exec.Command("/bin/bash", "-c", changePwd1).Run()
	if err != nil {
		return err
	}

	logs.Info("ansible refresh host")
	refresh := "ansible-playbook -i hosts.ini create_users.yml -u root"
	cmd := exec.Command("/bin/bash", "-c", refresh)
	cmd.Dir = "./devops"
	out, err := cmd.Output()

	exec.Command("/bin/bash", "-c", changePwd2).Run()

	if strings.Contains(Byte2String(out), "Permission denied") {
		return errors.New("当前版本仅支持所有服务器root密码一致！！")
	}
	if err != nil {
		return err
	}
	return nil
}
