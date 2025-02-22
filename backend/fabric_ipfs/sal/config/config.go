package config

type AppConfig struct {
	Mysql MysqlConfig `yaml:"mysql"`
}

type MysqlConfig struct {
	Dsn string `yaml:"dsn"`
}
