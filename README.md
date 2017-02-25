Создаем таблицу в pg.

docker pull postgres
docker run --name redirect -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=redirect -e POSTGRES_USER=user -d -p 5434:5432 postgres

Наполняем данными.
go run insert.go
 
