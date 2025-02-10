package config

type Config struct {
	Port       string `yaml:"PORT"`
	MonogDBUri string `yaml:"MONGO_DB_URI"`
}

