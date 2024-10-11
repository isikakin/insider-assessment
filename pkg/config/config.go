package config

type Configuration struct {
	ApiKey string
	Redis  RedisSettings
}

type RedisSettings struct {
	Url      string
	Password string
}
