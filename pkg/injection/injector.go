package injection

import (
	"os"

	"github.com/adi-kmt/bitespeed_backend/pkg/db"
	"github.com/adi-kmt/bitespeed_backend/pkg/repositories"
	"github.com/adi-kmt/bitespeed_backend/pkg/services"
)

func InjectDependencies() *services.Service {
	db_user, isErr := os.LookupEnv("DB_USER")
	if !isErr {
		db_user = "postgres"
	}

	db_password, isErr := os.LookupEnv("DB_PASSWORD")
	if !isErr {
		db_password = "password"
	}

	db_name, isErr := os.LookupEnv("DB_NAME")
	if !isErr {
		db_name = "bitespeed_db"
	}

	db_host, isErr := os.LookupEnv("DB_HOST")
	if !isErr {
		db_host = "127.0.0.1"
	}

	db_port, isErr := os.LookupEnv("DB_PORT")
	if !isErr {
		db_port = "5432"
	}

	dbConfig := db.NewDbConfig(
		db_user,
		db_password,
		db_name,
		db_host,
		db_port,
	)
	connConfig := db.InitPool(dbConfig)

	repository := repositories.NewRepository(connConfig.DbQueries, connConfig.DbPool)

	return services.NewService(repository)

}
