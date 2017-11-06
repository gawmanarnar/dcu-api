package main

import "database/sql"
import "fmt"

type character struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	ProperName   string   `json:"proper_name"`
	Level        int      `json:"level"`
	ProductImage string   `json:"product_image"`
	CardImage    string   `json:"card_image"`
	Affilitions  []string `json:"affiliation"`
}

func getCharacters(db *sql.DB) ([]character, error) {
	rows, err := db.Query("SELECT * FROM characters")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	characters := []character{}

	for rows.Next() {
		fmt.Println("hello")
		var c character
		if err := rows.Scan(&c.ID, &c.Name, &c.ProperName, &c.Level, &c.ProductImage, &c.CardImage, &c.Affilitions); err != nil {
			return nil, err
		}
		characters = append(characters, c)
	}

	return characters, nil
}
