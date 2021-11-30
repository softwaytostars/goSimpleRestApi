package config

type DatabaseConfig struct {
	Uri         string `yaml:"uri"`
	DBName      string `yaml:"dbname"`
	MaxPoolSize uint64 `yaml:"maxPoolSize"`
}

type KafkaServerConfig struct {
	Uri string `yaml:"uri"`
}

type Config struct {
	ServerConfig struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	DbConfig        DatabaseConfig    `yaml:"database"`
	StorageInMemory bool              `yaml:"storageInMemory"`
	EmailConsumers  int               `yaml:"nEmailConsumers"`
	KafkaConfig     KafkaServerConfig `yaml:"kafkaServer"`
}
