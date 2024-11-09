package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	postgresHostName = flag.String("postgresHostName", "postgres", "Specify postgres host name")
	postgresHostPort = flag.Uint("postgresHostPort", 5432, "Specify postgres host port")
	postgresUsername = flag.String("postgresUsername", "admin", "Specify postgres username")
	postgresPassword = flag.String("postgresPassword", "admin", "Specify postgres password")
	postgresDatabase = flag.String("postgresDatabase", "defaultdb", "Specify postgres database")
)

func panicWhenErrorNotNil(msg string, err error) {
	if err != nil {
		panic(msg + err.Error())
	}
}

func main() {
	flag.Parse()

	// 連接 postgres
	// postgres://admin:admin@postgres/defaultdb?sslmode=disable
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		*postgresUsername, *postgresPassword, *postgresHostName, *postgresHostPort, *postgresDatabase)
	db, err := sql.Open("postgres", dsn)
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
			ID         serial       PRIMARY KEY,
			Name       varchar(255) NOT NULL,
			Password   varchar(255)
		);
	`)
	panicWhenErrorNotNil("Failed to prepare create table sql:", err)
	_, err = stmt.Exec()
	panicWhenErrorNotNil("Failed to exec table creation:", err)
	stmt.Close()

	// 插入資料
	stmt, err = db.Prepare(`
		INSERT INTO Users (Name,Password)
		VALUES($1,$2) RETURNING ID;
	`)
	panicWhenErrorNotNil("Failed to prepare insert sql:", err)
	res, err := stmt.Exec("ianian", "init password")
	panicWhenErrorNotNil("Failed to exec insert:", err)
	stmt.Close()

	rowAffected, err := res.RowsAffected()
	panicWhenErrorNotNil("", err)
	log.Printf("rowAffected=%d", rowAffected)

	// // NOT SUPPORT
	// lastInsertId, err := res.LastInsertId()
	// panicWhenErrorNotNil("", err)
	// log.Printf("lastInsertId=%d", lastInsertId)

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
		UPDATE Users SET Password=$1 WHERE Name=$2;
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
		DELETE FROM Users WHERE Name=$1;
	`)
	panicWhenErrorNotNil("Failed to prepare update sql:", err)
	_, err = stmt.Exec("ianian")
	panicWhenErrorNotNil("Failed to exec prepare:", err)
	stmt.Close()
}
