package config

type Config struct {
	Addr_PORT string `toml:"addr_port"`
	Host      string `toml:"host"`
	Port      string `toml:"port"`
	Username  string `toml:"username"`
	Password  string `toml:"pass"`
	SSLmode   string `toml:"sslmode"`
	DBname    string `toml:"dbname"`
}

func CreateConfig() *Config {
	return &Config{}
}
