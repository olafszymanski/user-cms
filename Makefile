include .env

migrateup:
	docker run -v C:\Users\olafs\Desktop\user-cms\migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@localhost:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" up

migratedown:
	docker run -v C:\Users\olafs\Desktop\user-cms\migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@localhost:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" down -all