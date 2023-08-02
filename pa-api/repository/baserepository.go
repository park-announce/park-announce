package repository

import (
	"database/sql"
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/park-announce/pa-api/factory"
)

func (repository BaseRepository) GetConnection() *sqlx.DB {
	return repository.dbClient.pool
}

func (repository BaseRepository) CloseConnection() error {

	if repository.dbClient != nil && repository.dbClient.pool != nil {
		return repository.dbClient.pool.Close()
	}

	return nil
}

func (repository BaseRepository) Update(query string, args ...interface{}) (result sql.Result, err error) {

	res, err := Update(repository.GetConnection(), query, args...)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repository BaseRepository) Delete(query string, args ...interface{}) (result sql.Result, err error) {

	res, err := Delete(repository.GetConnection(), query, args...)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repository BaseRepository) Insert(query string, args ...interface{}) (result sql.Result, err error) {

	res, err := Insert(repository.GetConnection(), query, args...)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repository BaseRepository) GetById(instanceType string, id int64, query string) (interface{}, error) {

	result, err := Query(repository.GetConnection(), query, id)
	if err != nil {
		return nil, err
	}

	instance := factory.GetEntityInstance(instanceType)

	convertMapToStruct(result[0], instance)

	return instance, err
}

func (repository BaseRepository) GetByIdList(instanceType string, query string, id ...interface{}) ([]interface{}, error) {
	result, err := Query(repository.GetConnection(), query, id...)
	if err != nil {
		return nil, err
	}

	var items []interface{}
	if len(result) > 0 {
		items = make([]interface{}, 0, len(result))
		for _, item := range result {
			instance := factory.GetEntityInstance(instanceType)
			convertMapToStruct(item, instance)
			items = append(items, instance)
		}
	}

	return items, err
}

func (repository BaseRepository) GetAll(instanceType string, query string) (interface{}, error) {

	result, err := Query(repository.GetConnection(), query)
	if err != nil {
		return nil, err
	}

	var items []interface{}
	if len(result) > 0 {
		items = make([]interface{}, 0, len(result))
		for _, item := range result {
			instance := factory.GetEntityInstance(instanceType)
			convertMapToStruct(item, instance)
			items = append(items, instance)
		}
	}

	return items, err
}

func (repository BaseRepository) GetAllX(instanceType string, query string) (interface{}, error) {

	result, err := QueryX(repository.GetConnection(), query, instanceType)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Insert(db *sqlx.DB, query string, args ...interface{}) (result sql.Result, err error) {
	return exec(db, query, args...)
}

func Delete(db *sqlx.DB, query string, args ...interface{}) (result sql.Result, err error) {
	return exec(db, query, args...)
}

func Update(db *sqlx.DB, query string, args ...interface{}) (result sql.Result, err error) {
	return exec(db, query, args...)
}

func UpdateTransaction(tx *sql.Tx, query string, args ...interface{}) (result sql.Result, err error) {
	return execTransaction(tx, query, args...)
}

func Connect(driverName, dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driverName, dataSourceName)
	return db, err
}

func Query(db *sqlx.DB, query string, args ...interface{}) (result []map[string]interface{}, err error) {
	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		columnsData := make([]interface{}, len(columns))
		columnsPointersData := make([]interface{}, len(columns))

		for i := range columnsData {
			columnsPointersData[i] = &columnsData[i]
		}

		if err = rows.Scan(columnsPointersData...); err != nil {
			return nil, err
		}

		item := make(map[string]interface{})
		for i, columnName := range columns {
			val := columnsPointersData[i].(*interface{})
			item[columnName] = *val
		}

		result = append(result, item)

	}
	rows.Close()
	return
}

func QueryX(db *sqlx.DB, query string, instanceType string, args ...interface{}) (result []interface{}, err error) {
	rows, err := db.Queryx(query, args...)

	if err != nil {
		return nil, err
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	result = make([]interface{}, 0)
	for rows.Next() {
		instance := factory.GetEntityInstance(instanceType)
		rows.StructScan(instance)
		result = append(result, instance)
	}
	rows.Close()
	return
}

func exec(db *sqlx.DB, query string, args ...interface{}) (result sql.Result, err error) {

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(args...)

	if err != nil {
		return nil, err
	}

	return res, err
}

func execTransaction(tx *sql.Tx, query string, args ...interface{}) (result sql.Result, err error) {

	stmt, err := tx.Prepare(query)
	defer stmt.Close()

	if err != nil {
		return nil, err
	}

	res, err := tx.Stmt(stmt).Exec(args...)

	if err != nil {
		return nil, err
	} else {
		return res, err
	}
}

func convertMapToStruct(data map[string]interface{}, targetStruct interface{}) {
	t := reflect.ValueOf(targetStruct).Elem()
	st := reflect.TypeOf(targetStruct).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			continue
		}
		value := data[tag]
		if value != nil {
			val := t.FieldByName(field.Name)
			val.Set(reflect.ValueOf(data[tag]).Convert(val.Type()))
		}
	}
}
