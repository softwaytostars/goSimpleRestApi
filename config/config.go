package config

type DatabaseConfig struct {
	Uri         string `yaml:"uri"`
	DBName      string `yaml:"dbname"`
	MaxPoolSize uint64 `yaml:"maxPoolSize"`
}

type EmailServerConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	UseStartTLS bool   `yaml:"useStartTLS"`
}

type KafkaServerConfig struct {
	Uri string `yaml:"uri"`
}

type Config struct {
	ServerConfig struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	DbConfig          DatabaseConfig    `yaml:"database"`
	StorageInMemory   bool              `yaml:"storageInMemory"`
	EmailConsumers    int               `yaml:"nEmailConsumers"`
	KafkaConfig       KafkaServerConfig `yaml:"kafkaServer"`
	EmailServerConfig EmailServerConfig `yaml:"emailServer"`
}
