package models

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	ANSIBLE_LOG = "./devops/log/ansible.log"
)

type DevopsInfo struct {
	ExecTime string
}

func (m *DevopsInfo) BackupLog() {
	cmd := exec.Command("mv", ANSIBLE_LOG, ANSIBLE_LOG+"."+m.ExecTime)
	cmd.Start()
}

func (m *DevopsInfo) Deploy() {
	cmd := exec.Command("ansible-playbook", "music.yml")
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
	return int(binary.BigEndian.Uint16(out)), nil
}

type DevopsLogResult struct {
	CurrentPage int
	Rows        int
	Data        []string
}
