package utils

import (
	"testing"
)

func TestIniInventory_ReadInventory(t *testing.T) {
	i := IniInventory{}
	err := i.ReadInventory()
	if err != nil {
		t.Error(err)
	}
	if len(i.Servers) == 0 {
		t.Errorf("read ini error")
	}
}

func TestIniInventory_WriteInventory(t *testing.T) {
	i := IniInventory{}
	err := i.ReadInventory()
	if err != nil {
		t.Error(err)
	}

	i.MasterServers = []string{"192.167.8.141"}

	err = i.WriteInventory()
	if err != nil {
		t.Error(err)
	}

	if i.MasterServers[0] != "192.167.8.141" {
		t.Errorf("write ini error")
	}
}

func TestIniHosts_ReadHosts(t *testing.T) {
	i := IniHosts{}
	err := i.ReadHosts()
	if err != nil {
		t.Error(err)
	}
	if len(i.Servers) == 0 {
		t.Errorf("read ini error")
	}
}

func TestIniHosts_WriteHosts(t *testing.T) {
	i := IniHosts{}
	err := i.ReadHosts()
	if err != nil {
		t.Error(err)
	}

	i.Servers = []string{"192.167.8.141"}

	err = i.WriteHosts()
	if err != nil {
		t.Error(err)
	}

	if i.Servers[0] != "192.167.8.141" {
		t.Errorf("write ini error")
	}
}
