db-migrate:
	migrate -database postgres://postgres:password@db:5432/bitespeed_db?sslmode=disable -path=db/migrations/ up 1

inspect-db:
	 docker exec -it bitespeed_postgres psql -U postgres -W bitespeed_db
	#  password is password