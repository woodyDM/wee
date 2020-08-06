package support

type Configuration struct {
	Web   WebConfiguration   `json:"web"`
	Redis RedisConfiguration `json:"redis"`
	Mysql MysqlConfiguration `json:"mysql"`
}

type WebConfiguration struct {
	Port                int `json:"port"`
	SessionSeconds      int `json:"sessionSeconds"`
	CacheTimeoutSeconds int `json:"cacheTimeoutSeconds"`
}

type MysqlConfiguration struct {
	Url string `json:"url"`
}
type RedisConfiguration struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Auth string `json:"auth"`
}

var DefaultConfiguration = &Configuration{
	Web: WebConfiguration{
		Port:                8080,
		SessionSeconds:      3600,
		CacheTimeoutSeconds: 1800,
	},
	Redis: RedisConfiguration{
		Host: "localhost",
		Port: 6379,
		Auth: "123456",
	},
	Mysql: MysqlConfiguration{
		Url: "root:123456@tcp(127.0.0.1:3306)/blog?charset=utf8mb4",
	},
}
