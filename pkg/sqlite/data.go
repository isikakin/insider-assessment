package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)

const dataSource = "./data/sqlite-database.db"

func Insert() {
	var (
		db  *sql.DB
		err error
	)

	if db, err = sql.Open("sqlite", dataSource); err != nil {
		panic(err)
	}

	db.Exec("delete from messages where 1=1")

	defer db.Close()

	for i := range 50 {

		db.Exec("insert into messages (message_id, recipient, content, status, created_at) values($1, $2, $3, $4, $5)",
			uuid.NewString(),
			fmt.Sprintf("%d%d%d%d", i, i, i, i),
			fmt.Sprintf("%d.Recipient", i),
			1,
			time.Now().Format(time.DateTime))
	}
}

func CreateTable() {
	var (
		db  *sql.DB
		err error
	)

	os.Remove("./data/sqlite-database.db")

	if db, err = sql.Open("sqlite", dataSource); err != nil {
		panic(err)
	}

	defer db.Close()

	createMessageTableSQL := `CREATE TABLE IF NOT EXISTS messages (
		"message_id" TEXT NOT NULL PRIMARY KEY,		
		"recipient" TEXT NOT NULL,
		"content" VARCHAR(50) NOT NULL,
		"status" TINYINT,
		"sent_date" TEXT,
		"created_at" TEXT NOT NULL
	  );` // SQL Statement for Create Table

	log.Println("Create messages table...")
	statement, err := db.Prepare(createMessageTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("messages table created")

}
