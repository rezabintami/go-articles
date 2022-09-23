package postgre_driver

import (
	"database/sql"
	"fmt"

	"log"

	"go-articles/helpers"
	_config "go-articles/server/config"

	_ "github.com/lib/pq"
)


func GetConnection() string {
	if _config.GetConfiguration("app.env") == "DEV" {
	return fmt.Sprintf("user=%s host=%s dbname=%s sslmode=%s password=%s port=%s",
		_config.GetConfiguration("postgres.user"),
		_config.GetConfiguration("postgres.host"),
		_config.GetConfiguration("postgres.name"),
		_config.GetConfiguration("postgres.ssl"),
		_config.GetConfiguration("postgres.pass"),
		_config.GetConfiguration("postgres.port"))
	}
	return fmt.Sprintf("user=%s host=%s dbname=%s password=%s port=%s",
		_config.GetConfiguration("postgres.user"),
		_config.GetConfiguration("postgres.host"),
		_config.GetConfiguration("postgres.name"),
		_config.GetConfiguration("postgres.pass"),
		_config.GetConfiguration("postgres.port"))
}

func Init() *sql.DB {
	db, err := sql.Open("postgres",GetConnection())
	if err != nil {
		log.Fatal(err)
	}
	
	db.SetMaxOpenConns(helpers.ConvertStringtoInt(_config.GetConfiguration("postgres.maxOpenConnection")))
	db.SetMaxIdleConns(helpers.ConvertStringtoInt(_config.GetConfiguration("postgres.maxIdleConnection")))
	return db
}