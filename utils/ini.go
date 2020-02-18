package utils

import (
	. "ds-yibasuo/models"
	"io/ioutil"
)

// 因为ansible的ini特殊性。
// 在go中不能使用常用的解析包，比如go-ini，viper都不行。
// ansible 的 inventory.ini 和 hosts.ini 都不能算是正常语法的ini配置文件。
// 只能采取正则方式读，拼字符串方式写，这种菜鸡行为。
// 第一版先这样把，后面迭代改。

// 读ansible ini配置文件
func ReadInventory() {

}

// 写ansible ini配置文件
func WriteInventory(cluster ClusterInfo) error {
	data := "## DolphinScheduler Part\n"

	// TODO 中间拼接

	err := ioutil.WriteFile("./devops/inventory.ini", []byte(data), 0755)
	if err != nil {
		return err
	}
	return nil
}

// 读hosts ini配置文件
func ReadHosts() {

}

// 写hosts ini配置文件
func WriteHosts() {

}
