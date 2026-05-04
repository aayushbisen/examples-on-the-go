package storage

import (
	"database/sql"
	"fmt"
	"my-go-learning/sentinel/internal/types"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage() (*Storage, error) {
	//open a sql file
	db, err := sql.Open("sqlite", "sentinel.db")
	if err != nil {
		fmt.Printf("Failed to open db file %s", err)
		return nil, err
	}
	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Printf("Failed to connect to db %s", pingErr)
		return nil, pingErr
	}
	// make the table and columns
	// raw sel query string
	sqlQueryString := "CREATE TABLE IF NOT EXISTS results (website_id TEXT, timestamp DATETIME, status BOOLEAN);"

	db.Exec(sqlQueryString)

	return &Storage{DB: db}, nil
}

func (s *Storage) SaveResult(res types.CheckResult) error {
	insertQuery := "INSERT INTO results (website_id, timestamp, status) VALUES (?, ?, ?)"

	_, err := s.DB.Exec(insertQuery, res.WebsiteID, res.Time, res.Status)

	if err != nil {
		return err
	}

	return nil
}
