package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const RecordCount = 100

func main() {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=dev dbname=postgres password=pass sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	now := time.Now()
	dropTable(db)
	timeAfterDropTable := time.Now()
	fmt.Printf("dropTable: %s\n", timeAfterDropTable.Sub(now).String())

	migrate(db)
	timeAfterMigrate := time.Now()
	fmt.Printf("migrate: %s\n", timeAfterMigrate.Sub(timeAfterDropTable).String())

	insert(db)
	timeAfterInsert := time.Now()
	fmt.Printf("insert: %s\n", timeAfterInsert.Sub(timeAfterMigrate).String())

	selectAndUpdate(db)
	timeAfterSelectAndUpdate := time.Now()
	fmt.Printf("selectAndUpdate: %s\n", timeAfterSelectAndUpdate.Sub(timeAfterInsert).String())

	selectAndDelete(db)
	timeAfterSelectAndDelete := time.Now()
	fmt.Printf("selectAndDelete: %s\n", timeAfterSelectAndDelete.Sub(timeAfterSelectAndUpdate).String())
}

func dropTable(db *sqlx.DB) {
	// テーブルを削除
	db.MustExec("DROP TABLE IF EXISTS users")
}

func migrate(db *sqlx.DB) {
	// テーブルを作成
	sql := `
	CREATE TABLE users (
		id serial,
		created_at timestamp with time zone,
		updated_at timestamp with time zone,
		deleted_at timestamp with time zone,
		name text,
		PRIMARY KEY (id)
	);
	
	CREATE INDEX idx_users_deleted_at ON users(deleted_at);
	`
	db.MustExec(sql)
}

func insert(db *sqlx.DB) {
	for i := 1; i <= RecordCount; i++ {
		now := time.Now()
		// データを登録
		db.MustExec("INSERT INTO users(created_at, updated_at, name) VALUES ($1, $2, $3)", now, now, fmt.Sprintf("sqlx_test_user_%03d", i))
	}
}

func selectAndUpdate(db *sqlx.DB) {
	for i := 1; i <= RecordCount; i++ {
		var u User
		// データを取得
		if err := db.Get(&u, "SELECT * FROM users WHERE id=$1", i); err != nil {
			fmt.Printf("err(id = %d): %s\n", i, err.Error())
			continue
		}
		// データを更新
		db.MustExec("UPDATE users SET name = $1 WHERE id = $2", fmt.Sprintf("sqlx_test_user_%03d_updated", i), i)
	}
}

func selectAndDelete(db *sqlx.DB) {
	for i := 1; i <= RecordCount; i++ {
		// データを更新
		db.MustExec("DELETE FROM users WHERE id = $1", i)
	}
}
