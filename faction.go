package main

import (
	"database/sql"
)

type faction struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getFactions(db *sql.DB) ([]faction, error) {
	rows, err := db.Query("SELECT id, name FROM factions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	factions := []faction{}
	for rows.Next() {
		var f faction
		if err := rows.Scan(&f.ID, &f.Name); err != nil {
			return nil, err
		}
		factions = append(factions, f)
	}
	return factions, nil
}
