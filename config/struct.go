package config

// ExchangeConfig 交换机配置
type ExchangeConfig struct {
	Port uint16 `json:"port"`
	Key  string `json:"key"`
}

// ServiceConfig 服务端配置
type ServiceConfig struct {
	Exchange string `json:"exchange"`
	Key      string `json:"key"`
	ID       string `json:"id"`
	Password string `json:"password"`
}

// ClientConfig 客户端配置
type ClientConfig struct {
	Exchange string `json:"exchange"`
	Port     uint16 `json:"port"`
	Key      string `json:"key"`
	ID       string `json:"id"`
	Password string `json:"password"`
}

// Config 配置
type Config struct {
	Exchange *ExchangeConfig `json:"exchange"`
	Service  *ServiceConfig  `json:"service"`
	Client   *ClientConfig   `json:"client"`
}
