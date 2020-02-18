package utils

import (
	. "ds-yibasuo/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// 读yml配置文件，返回配置接口
func ReadYml(path string, conf *ConfigBody) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return err
	}
	return nil
}

// 写yml配置文件
func WriteYml(path string, confData []byte) error {
	err := ioutil.WriteFile(path, confData, 0755)
	if err != nil {
		return err
	}
	return nil
}
