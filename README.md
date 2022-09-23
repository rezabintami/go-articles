## Architecture
- Language [GoLang](https://golang.org/)
- Framework [Framework Echo](https://echo.labstack.com/)
- Object Relational Mapping [Gorm io](https://gorm.io/docs/index.html)
- [Postgres](https://www.postgresql.org/download/)
- [MongoDB](https://www.mongodb.com/try/download/community) (for logging)
- JWT [JWT](https://github.com/dgrijalva/jwt-go)
- Log [Logrus](https://github.com/sirupsen/logrus)

## Usage
-> Create database on your system, name is free.

1. Clone app from repo.
2. Navigate to project folder.
3. Execute go mod init.
4. Execute go mod vendor.
5. Adjust your config on `config.yaml` in `app/config/config.yaml`
6. Migrate your table using `go run cmd/migrations/migrations.go go-drop-logistik:migrate --up`

## How to Add Migration

1. Install [CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) in your computer.
2. Create file migration using `migrate create -ext sql -dir your_project_dir -seq your_file_name` 
3. Example cli command `migrate create -ext sql -dir drivers/postgres/files/migrations -seq add_foreign_key_tracks`

## How to Add Mock Usecase

1. Go to usecase folder.
2. run cli command `mockery --all`
3. Create Unit Test.