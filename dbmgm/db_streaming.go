package dbmgm

import (
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"muslog"
	"time"
)

type totalTask struct {
	Total []streamingTask `json:"total"`
}

type streamingTask struct {
	Id                string `json:"id"`
	Task_id           string `json:"task_id"`
	Task_type         string `json:"app_type"`
	Task_name         string `json:"app_name"`
	Task_des          string `json:"task_des"`
	Language          string `json:"language"`
	Data_source       string `json:"dataSource"`
	Submit_time       string `json:"submit_time"`
	Last_exec_time    string `json:"last_exec_time"`
	Exec_times        string `json:"exec_times"`
	Last_modify_time  string `json:"last_modify_time"`
	Modify_times      string `json:"modify_times"`
	Schedule_strategy string `json:"schedule_strategy"`
	Deploy_mode       string `json:"deploy_mode"`
	Driver_id         string `json:"driver_id"`
	App_id            string `json:"app_id"`
	Process_id        string `json:"p_id"`
	Exec_sequence     string `json:"exec_sequence"`
	Source_file       string `json:"source_file"`
}

func AddTask(name string, taskID string, taskType string, taskDes string, language string, schStrategy string, execTimes int, dataSource string, execSeq string, sourceFile string, deployMode string, appId string) error {
	db := createCon(dbName)
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		muslog.Error(err)
		return err
	}

	stmt, err := tx.Prepare("insert into task_streaming" +
		"(id, task_id, task_type, task_name, task_des,language, data_source, schedule_strategy, exec_times, deploy_mode, app_id, exec_sequence, source_file)" +
		"values (NULL,?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		muslog.Error(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(taskID, taskType, name, taskDes, language, dataSource, schStrategy, execTimes, deployMode, appId, execSeq, sourceFile)

	//	var id int
	//	err = tx.QueryRow("select last_insert_rowid() newid;").Scan(&id)
	//	if err != nil {
	//		log.Print(err)
	//		return err
	//	}
	//
	//	stmt2, err := tx.Prepare("insert into task_streaming_detail" +
	//		"(id, app_id, exec_sequence, source_file) values (?, ?, ?, ?);")
	//	_, err = stmt2.Exec(id, appId, execSeq, sourceFile)
	//	if err != nil {
	//		log.Print(err)
	//		return err
	//	}
	tx.Commit()
	return nil
}

func ListStreamingTask() (string, error) {
	db := createCon(dbName)
	defer db.Close()

	var total totalTask

	sql := "select * from task_streaming;"
	rows, err := db.Query(sql)
	if err != nil {
		muslog.Error(err)
		return "", err
	}
	for rows.Next() {
		var task streamingTask

		err = rows.Scan(&task.Id, &task.Task_id, &task.Task_type, &task.Task_name, &task.Task_des, &task.Language, &task.Data_source, &task.Submit_time, &task.Last_exec_time, &task.Exec_times, &task.Last_modify_time, &task.Modify_times, &task.Schedule_strategy, &task.Deploy_mode, &task.Driver_id, &task.App_id, &task.Process_id, &task.Exec_sequence, &task.Source_file)
		if err != nil {
			muslog.Error(err)
			return "", err
		}
		total.Total = append(total.Total, task)
	}
	res, err := json.Marshal(total)
	if err != nil {
		muslog.Error(err)
		return "", err
	}
	muslog.Trace("list streaming task json: " + string(res))
	return string(res), nil
}

func UpdateProcessId(taskId string, pid string) error {
	db := createCon(dbName)
	defer db.Close()
	stmt, err := db.Prepare("update task_streaming set process_id = ? where task_id = ?;")
	defer stmt.Close()
	if err != nil {
		muslog.Error(err)
		return err
	}
	_, err = stmt.Exec(pid, taskId)
	if err != nil {
		muslog.Error(err)
		return err
	}
	return nil
}

func UpdateDriverId(taskId string, driverId string) error {
	db := createCon(dbName)
	defer db.Close()
	stmt, err := db.Prepare("update task_streaming set driver_id = ? where task_id = ?;")
	defer stmt.Close()
	if err != nil {
		muslog.Error(err)
		return err
	}
	_, err = stmt.Exec(driverId, taskId)
	if err != nil {
		muslog.Error(err)
		return err
	}
	return nil
}

func UpdateAppId(taskId string, appId string) error {
	db := createCon(dbName)
	defer db.Close()
	stmt, err := db.Prepare("update task_streaming set app_id = ? where task_id = ?;")
	defer stmt.Close()
	if err != nil {
		muslog.Error(err)
		return err
	}
	_, err = stmt.Exec(appId, taskId)
	if err != nil {
		muslog.Error(err)
		return err
	}
	return nil
}

func UpdateExecTimes(taskId string) error {
	db := createCon(dbName)
	defer db.Close()

	execTime := time.Now().Format("2006-01-02 15:04:05")

	stmt, err := db.Prepare("update task_streaming set exec_times = exec_times+1, Last_exec_time = ? where task_id = ?;")
	defer stmt.Close()
	if err != nil {
		muslog.Error(err)
		return err
	}
	_, err = stmt.Exec(execTime, taskId)
	if err != nil {
		muslog.Error(err)
		return err
	}
	return nil
}

func UpdateStreamingTask(taskId string, name string, taskDes string, schStrategy string, execSeq string, appId string) error {
	db := createCon(dbName)
	defer db.Close()
	stmt, err := db.Prepare("update task_streaming set task_name = ?, task_des = ?, schedule_strategy = ?, exec_sequence = ?, app_id =?, Last_modify_time = ?, modify_times = modify_times+1 where task_id = ?;")
	defer stmt.Close()
	if err != nil {
		muslog.Error(err)
		return err
	}
	modifyTime := time.Now().Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(name, taskDes, schStrategy, execSeq, appId, modifyTime, taskId)
	if err != nil {
		muslog.Error(err)
		return err
	}
	return nil
}
