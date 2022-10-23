package main

type Backend struct {
	Metrics bool              `koanf:"metrics"`
	Port    int               `koanf:"port"`
	Users   map[string]string `koanf:"users"`
}

type MQTT struct {
	Protocol string `koanf:"protocol"`
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Secret   string `koanf:"secret"`
}

type BackendConfig struct {
	Backend `koanf:"backend"`
	MQTT    `koanf:"mqtt"`
}

var DefaultConfig = BackendConfig{
	Backend: Backend{
		Metrics: true,
		Port:    8081,
	},
	MQTT: MQTT{
		Protocol: "tcp",
		Port:     8080,
	},
}
