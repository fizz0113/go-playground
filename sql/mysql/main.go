package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	mysqlHostName = flag.String("mysqlHostName", "mysql", "Specify mysql host name")
	mysqlHostPort = flag.Uint("mysqlHostPort", 3306, "Specify mysql host port")
	mysqlUsername = flag.String("mysqlUsername", "admin", "Specify mysql username")
	mysqlPassword = flag.String("mysqlPassword", "admin", "Specify mysql password")
	mysqlDatabase = flag.String("mysqlDatabase", "defaultdb", "Specify mysql database")
)

func panicWhenErrorNotNil(msg string, err error) {
	if err != nil {
		panic(msg + err.Error())
	}
}

func main() {
	flag.Parse()

	// 連接 mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		*mysqlUsername, *mysqlPassword, *mysqlHostName, *mysqlHostPort, *mysqlDatabase)
	db, err := sql.Open("mysql", dsn)
	panicWhenErrorNotNil("Failed to connect mysql:", err)
	defer db.Close()

	// 移除 table
	stmt, err := db.Prepare(`
		DROP TABLE Users;
	`)
	panicWhenErrorNotNil("Failed to prepare delete table sql:", err)
	_, _ = stmt.Exec()
	stmt.Close()

	// 建立 table
	stmt, err = db.Prepare(`
		CREATE TABLE Users (
			ID int NOT NULL AUTO_INCREMENT,
			Name varchar(255) NOT NULL,
			Password varchar(255),
			PRIMARY KEY (ID,Name)
		);
	`)
	panicWhenErrorNotNil("Failed to prepare create table sql:", err)
	_, err = stmt.Exec()
	panicWhenErrorNotNil("Failed to exec table creation:", err)
	stmt.Close()

	// 插入資料
	stmt, err = db.Prepare(`
		INSERT INTO Users (Name,Password)
		VALUES(?,?);
	`)
	panicWhenErrorNotNil("Failed to prepare insert sql:", err)
	res, err := stmt.Exec("ianian", "init password")
	panicWhenErrorNotNil("Failed to exec insert:", err)
	stmt.Close()

	rowAffected, err := res.RowsAffected()
	panicWhenErrorNotNil("", err)
	log.Printf("rowAffected=%d", rowAffected)

	lastInsertId, err := res.LastInsertId()
	panicWhenErrorNotNil("", err)
	log.Printf("lastInsertId=%d", lastInsertId)

	// 查詢資料
	rows, err := db.Query(`
		SELECT * FROM Users;
	`)
	panicWhenErrorNotNil("Failed to query:", err)
	for rows.Next() {
		var id int
		var name, password string
		err = rows.Scan(&id, &name, &password)
		panicWhenErrorNotNil("Failed to query scan:", err)
		log.Printf("id=%d, name=%q, passoword=%q", id, name, password)
	}

	// 更新資料
	stmt, err = db.Prepare(`
		UPDATE Users SET Password=? WHERE Name=?;
	`)
	panicWhenErrorNotNil("Failed to prepare update sql:", err)
	_, err = stmt.Exec("updated password", "ianian")
	panicWhenErrorNotNil("Failed to exec prepare:", err)
	stmt.Close()

	// 查詢資料
	rows, err = db.Query(`
		SELECT * FROM Users;
	`)
	panicWhenErrorNotNil("Failed to query:", err)
	for rows.Next() {
		var id int
		var name, password string
		err = rows.Scan(&id, &name, &password)
		panicWhenErrorNotNil("Failed to query scan:", err)
		log.Printf("id=%d, name=%q, passoword=%q", id, name, password)
	}

	// 刪除資料
	stmt, err = db.Prepare(`
		DELETE FROM Users WHERE Name=?;
	`)
	panicWhenErrorNotNil("Failed to prepare update sql:", err)
	_, err = stmt.Exec("ianian")
	panicWhenErrorNotNil("Failed to exec prepare:", err)
	stmt.Close()
}
