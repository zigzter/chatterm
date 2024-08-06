package db

import (
	"database/sql"
	"log"
)

func InsertUserMap(db *sql.DB, username, userID string) {
	sqlStatement := `
        INSERT INTO userid_map (username, user_id) VALUES (?, ?)
        ON CONFLICT(username) DO UPDATE SET user_id = excluded.user_id;`
	_, err := db.Exec(sqlStatement, username, userID)
	if err != nil {
		log.Fatal("Cannot insert or update user map:", err)
	}
}

func GetUsername(db *sql.DB, userID string) (string, error) {
	var username string
	sqlStatement := "SELECT username FROM userid_map WHERE user_id = ? LIMIT 1"
	row := db.QueryRow(sqlStatement, userID)
	err := row.Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return username, nil
}

func GetUserId(db *sql.DB, username string) (string, error) {
	var userId string
	sqlStatement := "SELECT user_id FROM userid_map WHERE username = ? LIMIT 1"
	row := db.QueryRow(sqlStatement, username)

	err := row.Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return userId, nil
}
