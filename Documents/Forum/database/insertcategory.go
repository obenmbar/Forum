package forumino

import (
	"database/sql"

	forumino "forumino/models"
)

func INSERTcategory(db *sql.DB) error {
	category := []string{"Sport", "Sience", "Technologie", "Culture", "Busness", "Food", "Other"}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO categories (name) VALUES(?)")
	if err != nil {
		return err
	}

	for _, cataegoryName := range category {
		result, err := stmt.Exec(cataegoryName)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		forumino.CategorID.Categoryid[cataegoryName] = int(id)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
	
}
