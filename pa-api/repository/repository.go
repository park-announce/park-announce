package repository

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

type IRepository interface {
	GetById(id int64) (interface{}, error)
	GetAll(instanceType interface{}) (interface{}, error)
	Update(query string, args ...interface{}) (result sql.Result, err error)
	Delete(query string, args ...interface{}) (result sql.Result, err error)
	Insert(query string, args ...interface{}) (result sql.Result, err error)
}

type BaseRepository struct {
	dbClient *DBClient
}

type DBClient struct {
	pool *sqlx.DB
}

type DBClientFactory struct {
	driverName     string
	dataSourceName string
}

func NewDbClientFactory(driverName string, dataSourceName string) DBClientFactory {
	return DBClientFactory{driverName: driverName, dataSourceName: dataSourceName}
}

func (dbCLientFactory DBClientFactory) NewDBClient() *DBClient {
	client := &DBClient{}

	pool, err := Connect(dbCLientFactory.driverName, dbCLientFactory.dataSourceName)

	if err != nil {
		log.Println("error :", err)
		panic(err)
	}

	client.pool = pool

	return client
}
