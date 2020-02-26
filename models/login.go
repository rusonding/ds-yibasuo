package models

import (
	"ds-yibasuo/utils/black"
	"ds-yibasuo/utils/blotdb"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
)

var (
	userId  = black.String2Byte("yibasuo")
	userPwd = black.String2Byte("yibasuo")
)

type User struct {
	Password string `json:"password"`
}

func (m *User) UserChange(originPwd, newPwd string) (bool, error) {
	pwd, err := blotdb.Db.SelectVal("user", userId)
	if err != nil {
		return false, errors.New("数据库文件可能损坏，请下载最新的数据库文件。")
	}
	// 如果原始密码 等于 库里存的密码，则将新的密码赋值进去
	if pwd[0] == originPwd {
		err = blotdb.Db.Add("user", userId, black.String2Byte(newPwd))
		if err != nil {
			return false, err
		}
	} else {
		return false, errors.New("密码错误！")
	}
	return true, nil
}

func (m *User) UserCheck(password string) (bool, error) {
	pwd, err := blotdb.Db.SelectVal("user", userId)
	if err != nil {
		return false, errors.New("数据库文件可能损坏，请下载最新的数据库文件。")
	}
	if pwd[0] == password {
		return true, err
	} else {
		return false, errors.New("密码错误！")
	}
}

func UserInit() {
	logs.Info("Init user password")
	pwd, _ := blotdb.Db.SelectVal("user", userId)
	if len(pwd) == 0 {
		err := blotdb.Db.Add("user", userId, userPwd)
		if err != nil {
			logs.Error("init user err: ", err)
		}
	}
}
