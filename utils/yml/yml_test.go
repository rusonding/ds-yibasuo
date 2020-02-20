package yml

import (
	"ds-yibasuo/models"
	"testing"
)

func TestReadYml(t *testing.T) {
	path := "C:\\iceblue\\go\\src\\ds-yibasuo\\devops\\conf\\master.yml"

	conf, err := ReadYml(path, models.Master)
	if err != nil {
		t.Error(err)
	}

	if conf.(*models.ConfigMaster).Master.MasterReservedMemory != 0.1 {
		t.Errorf("read yml error")
	}
}

func TestWriteYml(t *testing.T) {
	path := "C:\\iceblue\\go\\src\\ds-yibasuo\\devops\\conf\\master.yml"

	conf, err := ReadYml(path, models.Master)
	if err != nil {
		t.Error(err)
	}

	conf.(*models.ConfigMaster).Master.MasterReservedMemory = 0.8

	err = WriteYml(path, conf)
	if err != nil {
		t.Error(err)
	}
}
