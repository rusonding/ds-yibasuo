package yml

import (
	"ds-yibasuo/models"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// 只能用这种方式来实现多态，感觉代码好冗余
// 后面再想办法迭代

// 读yml配置文件，返回配置接口
func ReadYml(path string, conf models.ConfigType) (models.ConfigBody, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	switch conf {
	case models.Frontend:
		config := models.ConfigFrontend{}
		if err := yaml.Unmarshal(data, &config); err == nil {
			return &config, nil
		}
	case models.Backend:
		config := models.ConfigBackend{}
		if err := yaml.Unmarshal(data, &config); err == nil {
			return &config, nil
		}
	case models.Alert:
		config := models.ConfigAlert{}
		if err := yaml.Unmarshal(data, &config); err == nil {
			return &config, nil
		}
	case models.Master:
		config := models.ConfigMaster{}
		if err := yaml.Unmarshal(data, &config); err == nil {
			return &config, nil
		}
	case models.Worker:
		config := models.ConfigWorker{}
		if err := yaml.Unmarshal(data, &config); err == nil {
			return &config, nil
		}
	case models.Hadoop:
		config := models.ConfigHadoop{}
		if err := yaml.Unmarshal(data, &config); err == nil {
			return &config, nil
		}
	case models.Common:
		config := models.ConfigCommon{}
		if err := yaml.Unmarshal(data, &config); err == nil {
			return &config, nil
		}
	default:
		return nil, errors.New("yaml unmarshal err or no fuck type.")
	}
	return nil, err
}

// 写yml配置文件
func WriteYml(path string, conf models.ConfigBody) error {
	var data []byte
	switch n := conf.(type) {
	case *models.ConfigFrontend:
		data, _ = yaml.Marshal(&n)
	case *models.ConfigBackend:
		data, _ = yaml.Marshal(&n)
	case *models.ConfigAlert:
		data, _ = yaml.Marshal(&n)
	case *models.ConfigMaster:
		data, _ = yaml.Marshal(&n)
	case *models.ConfigWorker:
		data, _ = yaml.Marshal(&n)
	case *models.ConfigHadoop:
		data, _ = yaml.Marshal(&n)
	case *models.ConfigCommon:
		data, _ = yaml.Marshal(&n)
	default:
		return errors.New("yaml marshal err or no fuck type.")
	}
	err := ioutil.WriteFile(path, data, 0755)
	if err != nil {
		return err
	}
	return nil
}
