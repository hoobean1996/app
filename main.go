package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func testrds() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		"root",                // 用户名
		"test123456",          // 密码
		"mysql-test-app:3306", // 服务名:端口
		"testdb",              // 数据库名
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("connect to mysql failed, err=%s\n", err.Error())
		return
	}
	if _, err := db.Exec(`
CREATE TABLE IF NOT EXISTS user (
   id INT AUTO_INCREMENT PRIMARY KEY,
   name VARCHAR(50) NOT NULL,
   age INT
)`); err != nil {
		fmt.Printf("create table to mysql failed, err=%s\n", err.Error())
		return
	}
	if _, err := db.Exec(`INSERT INTO user (name, age) VALUES (?, ?)`, "Alice", 25); err != nil {
		fmt.Printf("insert data to mysql failed, err=%s\n", err.Error())
		return
	}
	rows, err := db.Query("SELECT id, name, age FROM user")
	if err != nil {
		fmt.Printf("query data from mysql faild, err=%s\n", err.Error())
		return
	}
	for rows.Next() {
		var id int
		var name string
		var age int
		if err := rows.Scan(&id, &name, &age); err != nil {
			fmt.Printf("scan table user failed, err=%s\n", err.Error())
			continue
		}
		fmt.Printf("User: ID=%d, Name=%s, Age=%d\n", id, name, age)
	}
}

func testblob() {
	endpoint := "blob-test-app:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	err = minioClient.MakeBucket(context.Background(), "test-bucket", minio.MakeBucketOptions{
		Region: "asia",
	})
	if err != nil {
		fmt.Printf("create bucket failed, err=%s\n", err.Error())
	}
}

func main() {
	testrds()
	testblob()
	select {}
}
