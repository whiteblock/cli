package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	db *sql.DB
)

func getDB() {
	dataLoc := conf.StoreDirectory + ".store"
	needsInit := false
	// Check if the file exists, if not create the file dir
	_, err := os.Stat(dataLoc)
	if os.IsNotExist(err) {
		os.MkdirAll(conf.StoreDirectory, 0755)
		_, err := os.Create(dataLoc)
		if err != nil {
			log.Error(err)
		}
		needsInit = true
	}

	// Create a new DB object
	db, err = sql.Open("sqlite3", dataLoc)
	if err != nil {
		log.Fatal(err)
	}
	if needsInit {
		err := CreateTable()
		if err != nil {
			log.Panic(err)
		}
	}

}

func init() {
	getDB()
	log.Info("Finished initializing the db")
}

func CreateTable() error {
	_, err := db.Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS meta (key TEXT,value TEXT);"))
	return err
}

//Set stores a key value pair in the sql-lite database as json
func Set(key string, value interface{}) error {
	Delete(key)
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO meta (key,value) VALUES (?,?)"))
	if err != nil {
		return err
	}

	defer stmt.Close()

	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(key, string(v))
	if err != nil {
		return err
	}
	return tx.Commit()
}

//GetP fetches the value of key and returns it to v, v should be a pointer
func GetP(key string, v interface{}) error {
	row := db.QueryRow(fmt.Sprintf("SELECT value FROM meta WHERE key = \"%s\"", key))
	var data []byte
	err := row.Scan(&data)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &v)
}

func Exists(key string) bool {
	var out interface{}
	return (GetP(key, &out) == nil)
}

//Delete deletes the value stored at key
func Delete(key string) error {
	_, err := db.Exec(fmt.Sprintf("DELETE FROM meta WHERE key = \"%s\"", key))
	return err
}
