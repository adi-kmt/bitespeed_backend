package db

type DbConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func NewDbConfig(username, password, port, host, database string) *DbConfig {
	return &DbConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
	}
}
