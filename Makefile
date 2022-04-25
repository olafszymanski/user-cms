include .env

migrateup:
	docker run -v C:\Users\olafs\Desktop\user-cms\migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "${DATABASE_URL}?sslmode=disable" up

migratedown:
	docker run -v C:\Users\olafs\Desktop\user-cms\migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "${DATABASE_URL}?sslmode=disable" down -all