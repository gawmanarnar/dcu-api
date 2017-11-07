package main

import (
	"database/sql"

	"github.com/lib/pq"
)

type faction struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getFaction(db *sql.DB, id int) ([]character, error) {
	query := "SELECT c.id, c.name, c.proper_name, c.level, c.product_image, c.card_image, c.affiliation FROM characters c" +
		" JOIN factions f ON f.id = ANY(c.affiliation) WHERE f.id=$1"
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	characters := []character{}
	for rows.Next() {
		var c character
		if err := rows.Scan(&c.ID, &c.Name, &c.ProperName, &c.Level, &c.ProductImage, &c.CardImage, pq.Array(&c.Affilitions)); err != nil {
			return nil, err
		}
		characters = append(characters, c)
	}
	return characters, nil
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
