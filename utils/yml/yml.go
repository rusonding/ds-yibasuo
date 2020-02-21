package yml

import (
	. "ds-yibasuo/models"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// 只能用这种方式来实现多态，感觉代码好冗余
// 后面再想办法迭代

// 读yml配置文件，返回配置接口
func ReadYml(path string, conf ConfigType) (ConfigBody, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	switch conf {
	case Frontend:
		config := ConfigFrontend{}
		if err := yaml.Unmarshal(data, &config); err == nil {
			return &config, nil
		}
	case Backend:
		config := ConfigBackend{}
		if err := yaml.Unmarshal(data, &config); err == nil {
			return &config, nil
		}
	case Alert:
		config := ConfigAlert{}
		if err := yaml.Unmarshal(data, &config); err == nil {
			return &config, nil
		}
	case Master:
		config := ConfigMaster{}
		if err := yaml.Unmarshal(data, &config); err == nil {
			return &config, nil
		}
	case Worker:
		config := ConfigWorker{}
		if err := yaml.Unmarshal(data, &config); err == nil {
			return &config, nil
		}
	case Resources: // TODO 后面再写
	default:
		return nil, errors.New("yaml unmarshal err or no fuck type.")
	}
	return nil, err
}

// 写yml配置文件
func WriteYml(path string, conf ConfigBody) error {
	var data []byte
	switch n := conf.(type) {
	case *ConfigFrontend:
		data, _ = yaml.Marshal(&n)
	case *ConfigBackend:
		data, _ = yaml.Marshal(&n)
	case *ConfigAlert:
		data, _ = yaml.Marshal(&n)
	case *ConfigMaster:
		data, _ = yaml.Marshal(&n)
	case *ConfigWorker:
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
