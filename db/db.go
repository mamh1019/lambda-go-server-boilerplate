package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// NewDB는 환경변수 기반으로 MySQL 커넥션을 생성합니다.
// 필요한 환경변수:
//   - DB_HOST
//   - DB_USER
//   - DB_PASS
//   - DB_NAME
func NewDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		user, pass, host, name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
