package config

import (
	"encoding/json"
	"io/ioutil"
)

// FromJSONFile 从JSON文件中读取配置
func FromJSONFile(file string) (*Config, error) {
	configBytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	var cfg Config
	err = json.Unmarshal(configBytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
