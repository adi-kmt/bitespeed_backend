package db

type dbConfig struct {
	host     string
	port     string
	username string
	password string
	database string
}

func NewDbConfig(username, password, port, host, database string) *dbConfig {
	return &dbConfig{
		host:     host,
		port:     port,
		username: username,
		password: password,
		database: database,
	}
}
