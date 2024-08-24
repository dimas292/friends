package config

import (
	"database/sql"
	"fmt"
	 _ "github.com/lib/pq"
)


func ConnectDb()(*sql.DB,error){

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		"localhost",
		"5432",
		"postgres",
		"root",
		"friends",
	)

	db , err := sql.Open("postgres", dsn)
	if err != nil{
		return nil, err
	}

	if err := db.Ping(); err != nil{
		return nil, err
	}

	return db, nil
}
