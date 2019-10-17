package config

type (
	// Config will holds mapped key value for service configuration
	Config struct {
		Server   ServerConfig   `yaml:"server"`
		Database DatabaseConfig `yaml:"database"`
		Redis    RedisConfig    `yaml:"redis"`
	}

	// ServerConfig server config
	ServerConfig struct {
		Port string `yaml:"port"`
	}

	// DatabaseConfig db config
	DatabaseConfig struct {
		Master   string `yaml:"master"`
		Follower string `yaml:"follower"`
	}

	// RedisConfig redis config
	RedisConfig struct {
		MaxIdle   int    `yaml:"maxidle"`
		MaxActive int    `yaml:"maxactive"`
		TimeOut   int    `yaml:"timeout"`
		Wait      bool   `yaml:"wait"`
		Address   string `yaml:"address"`
	}
)
