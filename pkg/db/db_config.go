package db

type dbConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func NewDbConfig(username, password, port, host, database string) *dbConfig {
	return &dbConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
	}
}
