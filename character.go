package main

import "database/sql"
import "github.com/lib/pq"

type character struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	ProperName   string   `json:"proper_name"`
	Level        int      `json:"level"`
	ProductImage string   `json:"product_image"`
	CardImage    string   `json:"card_image"`
	Affilitions  []string `json:"affiliation"`
}

func (c *character) getCharacter(db *sql.DB) error {
	row := db.QueryRow("SELECT id, name, proper_name, level, product_image, card_image, affiliation FROM characters WHERE id=$1", c.ID)
	return row.Scan(&c.ID, &c.Name, &c.ProperName, &c.Level, &c.ProductImage, &c.CardImage, pq.Array(&c.Affilitions))
}

func getCharacters(db *sql.DB) ([]character, error) {
	rows, err := db.Query("SELECT id, name, proper_name, level, product_image, card_image, affiliation FROM characters")

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
