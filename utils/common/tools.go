package common

import (
	"github.com/google/uuid"
	"time"
)

func Now() string {
	//return time.Now().Format("2006-01-02 15:04:05")
	return time.Now().Format("20060102150405")
}

func MakeUuid(str string) string {
	rawUuid := uuid.NewMD5(uuid.NameSpaceDNS, []byte(str)).String()
	return rawUuid[0:8]
}
