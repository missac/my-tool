package dbmgm

import (
	"database/sql"
	//"fmt"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"muslog"
	"os"
)

var dbName = "mustang.db"

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func createCon(name string) *(sql.DB) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		muslog.Error(err)
		return nil
	}
	db.Exec("PRAGMA foreign_keys = ON")
	return db
}

func getJSON_0(rows *sql.Rows) (string, error) {
	columns, err := rows.Columns()
	if err != nil {
		muslog.Error(err)
		return "", err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		muslog.Error(err)
		return "", err
	}
	fmt.Println(string(jsonData))
	return string(jsonData), nil
}

func getJSON(db *sql.DB, sqlString string) (string, error) {
	rows, err := db.Query(sqlString)
	if err != nil {
		muslog.Error(err)
		return "", err
	}
	defer rows.Close()
	return getJSON_0(rows)
}

func InitDB(dbname string) {
	if isExist(dbname) {
		return
	}
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		muslog.Error(err)
		return
	}
	defer db.Close()

	db.Exec("PRAGMA foreign_keys = ON")
	sqlCreate := `
	create table task_streaming ( 
		id  integer not null primary key AUTOINCREMENT, 
		task_id  text not null, 
		task_type  text not null,
        task_name text not null,
        task_des text,
        language text not null,
		data_source text,
		submit_time TIMESTAMP default(datetime('now', 'localtime')),
		last_exec_time TIMESTAMP default(datetime('now', 'localtime')),
		exec_times integer default 1,
		last_modify_time TIMESTAMP default(datetime('now', 'localtime')),
		modify_times integer default 0,
        schedule_strategy text default 'now',
		deploy_mode text default 'client',
		driver_id text,
		task_status text default 'stopped',
		app_id text,
		process_id text,
		exec_seq text,
		source_file text);`

	muslog.Trace("create table use:" + sqlCreate)
	_, err = db.Exec(sqlCreate)
	if err != nil {
		muslog.Error(err)
		return
	}
}
