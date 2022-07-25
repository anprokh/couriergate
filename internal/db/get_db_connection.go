package db

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"time"
)

var DB_COURIER *sql.DB

// получить соединение с бд MS SQL
// 01
func GetDBConnection(ConnectionString string) (*sql.DB, error) {

	// создание соединения с MS SQL
	db, err := sql.Open("sqlserver", ConnectionString)
	if err != nil {
		return nil, err
	}
	//defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// проверка соединения
	err = db.Ping()
	if err != nil {
		log.Println("DB PING ERROR")
		return nil, err
	} else {
		log.Println("DB PING OK")
	}

	return db, nil
}

// TODO: обновить запись таблицы Document.Out
// 02
func UpdateDocumentOutRecord(FileID string, Attribute string) {

	// дата/время изменения записи
	DATE_TIME := time.Now().String()
	DATE_TIME = fmt.Sprintf("%sT%s", DATE_TIME[:10], DATE_TIME[11:23])

}
